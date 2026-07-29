package filesystem

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"mdvault/domain/model"
	"mdvault/domain/repository"

	"github.com/fsnotify/fsnotify"
	"github.com/goccy/go-yaml"
)

type vaultRepositoryImpl struct {
	watcher  *fsnotify.Watcher
	onChange func(path string)
	mu       sync.Mutex
}

func NewVaultRepository() repository.VaultRepository {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Warning: failed to init fsnotify watcher: %v", err)
	}
	repo := &vaultRepositoryImpl{
		watcher: w,
	}
	if w != nil {
		go repo.listenEvents()
	}
	return repo
}

func (r *vaultRepositoryImpl) listenEvents() {
	for {
		select {
		case event, ok := <-r.watcher.Events:
			if !ok {
				return
			}
			if strings.HasSuffix(strings.ToLower(event.Name), ".md") || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
				r.mu.Lock()
				cb := r.onChange
				r.mu.Unlock()
				if cb != nil {
					cb(event.Name)
				}
			}
		case err, ok := <-r.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (r *vaultRepositoryImpl) StartWatcher(rootPath string, onChange func(path string)) error {
	r.mu.Lock()
	r.onChange = onChange
	r.mu.Unlock()
	return r.UpdateWatchPath(rootPath)
}

func (r *vaultRepositoryImpl) StopWatcher() error {
	if r.watcher != nil {
		return r.watcher.Close()
	}
	return nil
}

func (r *vaultRepositoryImpl) UpdateWatchPath(dirPath string) error {
	if r.watcher == nil {
		return nil
	}
	for _, w := range r.watcher.WatchList() {
		_ = r.watcher.Remove(w)
	}
	return filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "dist" || name == "build" || name == "wailsjs" {
				return filepath.SkipDir
			}
			_ = r.watcher.Add(path)
		}
		return nil
	})
}

func (r *vaultRepositoryImpl) ParseNote(filePath string, folderID string) (*model.Note, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	contentStr := string(contentBytes)
	lines := strings.Split(contentStr, "\n")

	var frontmatter model.Frontmatter
	var bodyLines []string
	inFrontmatter := false
	frontmatterFound := false
	var yamlLines []string

	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		inFrontmatter = true
		for i := 1; i < len(lines); i++ {
			line := lines[i]
			if strings.TrimSpace(line) == "---" && inFrontmatter {
				inFrontmatter = false
				frontmatterFound = true
				bodyLines = lines[i+1:]
				break
			}
			if inFrontmatter {
				yamlLines = append(yamlLines, line)
			}
		}
	}

	if !frontmatterFound {
		bodyLines = lines
	} else {
		yamlStr := strings.Join(yamlLines, "\n")
		_ = yaml.Unmarshal([]byte(yamlStr), &frontmatter)
	}

	title := frontmatter.Title
	if strings.TrimSpace(title) == "" {
		for _, line := range bodyLines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "# ") {
				title = strings.TrimPrefix(trimmed, "# ")
				break
			}
		}
	}
	if strings.TrimSpace(title) == "" {
		base := filepath.Base(filePath)
		title = strings.TrimSuffix(base, filepath.Ext(base))
	}

	if frontmatter.Tags == nil {
		frontmatter.Tags = []string{}
	}

	bodyText := strings.Join(bodyLines, "\n")
	excerpt := extractExcerpt(bodyLines)
	normalizedPath := filepath.ToSlash(filePath)

	return &model.Note{
		ID:           normalizedPath,
		Title:        title,
		Content:      bodyText,
		Tags:         frontmatter.Tags,
		Folder:       folderID,
		LastModified: info.ModTime().Format(time.RFC3339),
		Excerpt:      excerpt,
	}, nil
}

func extractExcerpt(lines []string) string {
	var cleanText []string
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if t != "" && !strings.HasPrefix(t, "---") && !strings.HasPrefix(t, "```") && !strings.HasPrefix(t, "#") {
			cleanText = append(cleanText, t)
		} else if strings.HasPrefix(t, "#") {
			cleanText = append(cleanText, strings.TrimSpace(strings.TrimLeft(t, "#")))
		}
		if len(cleanText) >= 3 {
			break
		}
	}
	joined := strings.Join(cleanText, " ")
	if utf8.RuneCountInString(joined) > 100 {
		runes := []rune(joined)
		return string(runes[:100]) + "..."
	}
	if joined == "" {
		return "Ghi chú trống..."
	}
	return joined
}

func (r *vaultRepositoryImpl) SaveNote(filePath string, content string, title string, tags []string) error {
	var sb strings.Builder
	sb.WriteString("---\n")
	fm := model.Frontmatter{
		Title: title,
		Tags:  tags,
	}
	out, err := yaml.Marshal(fm)
	if err == nil {
		sb.Write(out)
	} else {
		sb.WriteString(fmt.Sprintf("title: %s\n", title))
	}
	sb.WriteString("---\n\n")
	sb.WriteString(strings.TrimSpace(content))
	sb.WriteString("\n")

	err = os.WriteFile(filePath, []byte(sb.String()), 0644)
	if err == nil {
		r.mu.Lock()
		cb := r.onChange
		r.mu.Unlock()
		if cb != nil {
			cb(filePath)
		}
	}
	return err
}

func (r *vaultRepositoryImpl) ScanWorkspace(rootPath string) (*model.WorkspaceData, error) {
	if rootPath == "" {
		return nil, fmt.Errorf("root path is empty")
	}

	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}
	absRoot = filepath.ToSlash(absRoot)

	rootFolderMap := make(map[string]*model.Folder)
	var allNotes []model.Note
	tagSet := make(map[string]struct{})
	var rootFolders []*model.Folder

	err = filepath.WalkDir(absRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		slashPath := filepath.ToSlash(path)
		if slashPath == absRoot {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == "dist" || name == "build" || name == "wailsjs" {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(absRoot, path)
		if err != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		parts := strings.Split(rel, "/")
		depth := len(parts) - 1

		if d.IsDir() {
			folder := &model.Folder{
				ID:       slashPath,
				Name:     name,
				Icon:     "📁",
				Level:    depth,
				Children: []*model.Folder{},
				Type:     "markdown",
			}
			rootFolderMap[slashPath] = folder

			parentPath := filepath.ToSlash(filepath.Dir(path))
			if parentPath == absRoot || depth == 0 {
				rootFolders = append(rootFolders, folder)
			} else {
				if parent, ok := rootFolderMap[parentPath]; ok {
					parent.Children = append(parent.Children, folder)
				}
			}
		} else if strings.ToLower(filepath.Ext(name)) == ".md" {
			parentDir := filepath.ToSlash(filepath.Dir(path))
			folderID := parentDir
			if parentDir == absRoot {
				folderID = "root"
			}
			note, err := r.ParseNote(path, folderID)
			if err == nil && note != nil {
				allNotes = append(allNotes, *note)
				for _, t := range note.Tags {
					tagSet[t] = struct{}{}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	var tags []string
	for t := range tagSet {
		tags = append(tags, t)
	}
	sort.Strings(tags)

	return &model.WorkspaceData{
		RootPath: absRoot,
		Folders:  rootFolders,
		Notes:    allNotes,
		AllTags:  tags,
	}, nil
}

func (r *vaultRepositoryImpl) CreateFolder(parentPath string, folderName string) error {
	newPath := filepath.Join(parentPath, folderName)
	err := os.MkdirAll(newPath, 0755)
	if err == nil && r.watcher != nil {
		_ = r.watcher.Add(newPath)
		r.mu.Lock()
		cb := r.onChange
		r.mu.Unlock()
		if cb != nil {
			cb(newPath)
		}
	}
	return err
}

func (r *vaultRepositoryImpl) CreateNote(folderPath string, noteTitle string) (*model.Note, error) {
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return nil, err
	}
	fileName := strings.ToLower(strings.ReplaceAll(noteTitle, " ", "-"))
	if !strings.HasSuffix(fileName, ".md") {
		fileName += ".md"
	}
	filePath := filepath.Join(folderPath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(file)
	writer.WriteString(fmt.Sprintf("---\ntitle: \"%s\"\ntags: []\n---\n\n# %s\n\nBắt đầu viết nội dung Ghi chú của bạn tại đây...\n", noteTitle, noteTitle))
	_ = writer.Flush()
	file.Close()

	r.mu.Lock()
	cb := r.onChange
	r.mu.Unlock()
	if cb != nil {
		cb(filePath)
	}
	return r.ParseNote(filePath, folderPath)
}

// Verify compilation at compile-time
var _ repository.VaultRepository = (*vaultRepositoryImpl)(nil)
