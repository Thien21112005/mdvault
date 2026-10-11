package filesystem

import (
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
	watcher        *fsnotify.Watcher
	onChange       func(path string)
	done           chan struct{}
	debounceLock   sync.Mutex
	debounceTimers map[string]*time.Timer
	mu             sync.Mutex
}

func NewVaultRepository() repository.VaultRepository {
	return &vaultRepositoryImpl{
		debounceTimers: make(map[string]*time.Timer),
	}
}

func (r *vaultRepositoryImpl) stopWatcherLocked() error {
	r.debounceLock.Lock()
	for p, t := range r.debounceTimers {
		t.Stop()
		delete(r.debounceTimers, p)
	}
	r.debounceLock.Unlock()

	if r.done != nil {
		select {
		case <-r.done:
			// already closed
		default:
			close(r.done)
		}
		r.done = nil
	}
	var err error
	if r.watcher != nil {
		err = r.watcher.Close()
		r.watcher = nil
	}
	return err
}

func (r *vaultRepositoryImpl) StopWatcher() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.stopWatcherLocked()
}

func (r *vaultRepositoryImpl) StartWatcher(rootPath string, onChange func(path string)) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.onChange = onChange

	cleanRoot := strings.TrimSpace(rootPath)
	if cleanRoot == "" {
		return nil
	}

	// Dọn dẹp an toàn watcher trước đó nếu đang chạy
	_ = r.stopWatcherLocked()

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("không thể khởi tạo fsnotify watcher: %w", err)
	}

	done := make(chan struct{})
	r.watcher = w
	r.done = done

	// Đăng ký toàn bộ thư mục đệ quy trong workspace
	if err := r.addDirRecursiveLocked(cleanRoot); err != nil {
		_ = r.stopWatcherLocked()
		return err
	}

	go r.listenEvents(w, done)
	return nil
}

func (r *vaultRepositoryImpl) addDirRecursiveLocked(dirPath string) error {
	if r.watcher == nil || strings.TrimSpace(dirPath) == "" {
		return nil
	}
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		absPath = dirPath
	}
	return filepath.WalkDir(absPath, func(path string, d os.DirEntry, err error) error {
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

func (r *vaultRepositoryImpl) UpdateWatchPath(dirPath string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cleanDir := strings.TrimSpace(dirPath)
	if cleanDir == "" {
		return nil
	}

	if r.watcher == nil {
		r.mu.Unlock()
		err := r.StartWatcher(cleanDir, r.onChange)
		r.mu.Lock()
		return err
	}

	for _, w := range r.watcher.WatchList() {
		_ = r.watcher.Remove(w)
	}
	return r.addDirRecursiveLocked(cleanDir)
}

func (r *vaultRepositoryImpl) listenEvents(w *fsnotify.Watcher, done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			// Tự động nhận diện và theo dõi thư mục mới tạo từ bên ngoài
			if event.Op&fsnotify.Create == fsnotify.Create {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					r.mu.Lock()
					_ = r.addDirRecursiveLocked(event.Name)
					r.mu.Unlock()
				}
			}

			// Chỉ thông báo sự kiện liên quan đến file markdown hoặc cấu trúc thư mục (qua Debounce 50ms)
			if isRelevantWatchEvent(event) {
				r.triggerDebounced(event.Name)
			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (r *vaultRepositoryImpl) triggerDebounced(path string) {
	r.debounceLock.Lock()
	defer r.debounceLock.Unlock()

	if r.debounceTimers == nil {
		r.debounceTimers = make(map[string]*time.Timer)
	}

	if timer, exists := r.debounceTimers[path]; exists {
		timer.Stop()
	}

	r.debounceTimers[path] = time.AfterFunc(50*time.Millisecond, func() {
		r.debounceLock.Lock()
		delete(r.debounceTimers, path)
		r.debounceLock.Unlock()

		r.mu.Lock()
		cb := r.onChange
		r.mu.Unlock()

		if cb != nil {
			cb(path)
		}
	})
}

func isRelevantWatchEvent(event fsnotify.Event) bool {
	base := filepath.Base(event.Name)
	if strings.HasPrefix(base, ".") || strings.HasSuffix(base, ".tmp") || strings.HasSuffix(base, ".swp") {
		return false
	}
	ext := strings.ToLower(filepath.Ext(base))
	if ext == ".md" || ext == ".markdown" {
		return true
	}
	if info, err := os.Stat(event.Name); err == nil {
		return info.IsDir()
	}
	if ext == "" && event.Op&(fsnotify.Remove|fsnotify.Rename) != 0 {
		return true
	}
	return false
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
	var rawMeta map[string]interface{}
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
		_ = yaml.Unmarshal([]byte(yamlStr), &rawMeta)
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
		Metadata:     rawMeta,
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

func atomicWrite(filePath string, data []byte) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dir, ".mdvault_tmp_*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return err
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	var renameErr error
	for attempts := 0; attempts < 3; attempts++ {
		renameErr = os.Rename(tmpPath, filePath)
		if renameErr == nil {
			cleanup = false
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}

	return renameErr
}

func extractExistingFrontmatterMap(filePath string) map[string]interface{} {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return make(map[string]interface{})
	}
	contentStr := string(contentBytes)
	lines := strings.Split(contentStr, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return make(map[string]interface{})
	}

	var yamlLines []string
	found := false
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "---" {
			found = true
			break
		}
		yamlLines = append(yamlLines, line)
	}
	if !found {
		return make(map[string]interface{})
	}

	metaMap := make(map[string]interface{})
	_ = yaml.Unmarshal([]byte(strings.Join(yamlLines, "\n")), &metaMap)
	if metaMap == nil {
		return make(map[string]interface{})
	}
	return metaMap
}

func (r *vaultRepositoryImpl) SaveNote(filePath string, content string, title string, tags []string) error {
	metaMap := extractExistingFrontmatterMap(filePath)
	if metaMap == nil {
		metaMap = make(map[string]interface{})
	}
	metaMap["title"] = title
	if tags == nil {
		metaMap["tags"] = []string{}
	} else {
		metaMap["tags"] = tags
	}

	out, err := yaml.Marshal(metaMap)
	var sb strings.Builder
	sb.WriteString("---\n")
	if err == nil {
		sb.Write(out)
	} else {
		sb.WriteString(fmt.Sprintf("title: %s\n", title))
	}
	sb.WriteString("---\n\n")
	sb.WriteString(strings.TrimSpace(content))
	sb.WriteString("\n")

	return atomicWrite(filePath, []byte(sb.String()))
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
	if err == nil {
		r.mu.Lock()
		if r.watcher != nil {
			_ = r.watcher.Add(newPath)
		}
		r.mu.Unlock()
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

	initialContent := fmt.Sprintf("---\ntitle: \"%s\"\ntags: []\n---\n\n# %s\n\nBắt đầu viết nội dung Ghi chú của bạn tại đây...\n", noteTitle, noteTitle)
	if err := atomicWrite(filePath, []byte(initialContent)); err != nil {
		return nil, err
	}

	return r.ParseNote(filePath, folderPath)
}

// Verify compilation at compile-time
var _ repository.VaultRepository = (*vaultRepositoryImpl)(nil)
