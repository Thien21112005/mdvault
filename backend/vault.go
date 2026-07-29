package backend

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/goccy/go-yaml"
)

// ParseMarkdownFile đọc và parse frontmatter cùng nội dung tệp .md
func ParseMarkdownFile(filePath string, folderID string) (*Note, error) {
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

	var frontmatter Frontmatter
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

	// Xác định tiêu đề
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

	// Tạo excerpt từ body
	bodyText := strings.Join(bodyLines, "\n")
	excerpt := extractExcerpt(bodyLines)

	normalizedPath := filepath.ToSlash(filePath)

	return &Note{
		ID:           normalizedPath,
		Title:        title,
		Content:      bodyText,
		Tags:         frontmatter.Tags,
		Folder:       folderID,
		LastModified: info.ModTime(),
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

// SaveMarkdownFile ghi lại content và frontmatter vào tệp .md thật
func SaveMarkdownFile(filePath string, content string, title string, tags []string) error {
	var sb strings.Builder
	sb.WriteString("---\n")
	fm := Frontmatter{
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

	return os.WriteFile(filePath, []byte(sb.String()), 0644)
}

// ScanWorkspace quét cây thư mục từ root và trả về danh sách file markdown và folders
func ScanWorkspace(rootPath string) (*WorkspaceData, error) {
	if rootPath == "" {
		return nil, fmt.Errorf("root path is empty")
	}

	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}
	absRoot = filepath.ToSlash(absRoot)

	rootFolderMap := make(map[string]*Folder)
	var allNotes []Note
	tagSet := make(map[string]struct{})
	var rootFolders []*Folder

	// Dùng filepath.WalkDir để duyệt filesystem
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
			folder := &Folder{
				ID:       slashPath,
				Name:     name,
				Icon:     "📁",
				Level:    depth,
				Children: []*Folder{},
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
			// Nếu file mằn ngay ở root workspace
			folderID := parentDir
			if parentDir == absRoot {
				folderID = "root"
			}
			note, err := ParseMarkdownFile(path, folderID)
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

	return &WorkspaceData{
		RootPath: absRoot,
		Folders:  rootFolders,
		Notes:    allNotes,
		AllTags:  tags,
	}, nil
}

// CopyFile hỗ trợ sao lưu nếu cần
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// CreateNewFolder tạo tệp thư mục mới
func CreateNewFolder(path string) error {
	return os.MkdirAll(path, 0755)
}

// CreateNewNote tạo 1 file .md mới với template frontmatter
func CreateNewNote(filePath string, title string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(fmt.Sprintf("---\ntitle: \"%s\"\ntags: []\n---\n\n# %s\n\nBắt đầu viết nội dung của bạn tại đây...\n", title, title))
	return writer.Flush()
}
