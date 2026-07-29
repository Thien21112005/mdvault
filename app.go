package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"mdvault/backend"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	workspacePath string
	watcher       *fsnotify.Watcher
	mu            sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Warning: failed to initialize fsnotify watcher: %v", err)
	}
	return &App{
		watcher: watcher,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.watcher != nil {
		go a.listenToWatchEvents()
	}

	// Mở mặc định thư mục markdown/ của dự án (nếu tồn tại) làm mẫu ban đầu
	cwd, err := os.Getwd()
	if err == nil {
		defaultWs := filepath.Join(cwd, "markdown")
		if info, err := os.Stat(defaultWs); err == nil && info.IsDir() {
			_, _ = a.OpenWorkspacePath(defaultWs)
		}
	}
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.watcher != nil {
		_ = a.watcher.Close()
	}
}

func (a *App) listenToWatchEvents() {
	for {
		select {
		case event, ok := <-a.watcher.Events:
			if !ok {
				return
			}
			// Phát ra sự kiện thay đổi đến phía Vite React Frontend
			if strings.HasSuffix(strings.ToLower(event.Name), ".md") || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
				runtime.EventsEmit(a.ctx, "on-file-change", event.Name)
			}
		case err, ok := <-a.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (a *App) addWatchersRecursively(dir string) {
	if a.watcher == nil {
		return
	}
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "dist" || name == "build" || name == "wailsjs" {
				return filepath.SkipDir
			}
			_ = a.watcher.Add(path)
		}
		return nil
	})
}

// SelectWorkspace mở hộp thoại chọn thư mục và kích hoạt theo dõi
func (a *App) SelectWorkspace() (*backend.WorkspaceData, error) {
	selectedDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Chọn thư mục Workspace ghi chú Markdown",
	})
	if err != nil {
		return nil, err
	}
	if selectedDir == "" {
		return nil, fmt.Errorf("không có thư mục nào được chọn")
	}
	return a.OpenWorkspacePath(selectedDir)
}

// OpenWorkspacePath trực tiếp mở và quét workspace theo đường dẫn
func (a *App) OpenWorkspacePath(dirPath string) (*backend.WorkspaceData, error) {
	a.mu.Lock()
	a.workspacePath = dirPath
	a.mu.Unlock()

	// Cập nhật fsnotify watch list
	if a.watcher != nil {
		for _, w := range a.watcher.WatchList() {
			_ = a.watcher.Remove(w)
		}
		a.addWatchersRecursively(dirPath)
	}

	return backend.ScanWorkspace(dirPath)
}

// ScanFolder làm mới dữ liệu cho workspace hiện tại
func (a *App) ScanFolder() (*backend.WorkspaceData, error) {
	a.mu.Lock()
	path := a.workspacePath
	a.mu.Unlock()

	if path == "" {
		return nil, fmt.Errorf("workspace chưa được thiết lập")
	}
	return backend.ScanWorkspace(path)
}

// GetNoteContent lấy nôi dung chi tiết một file note
func (a *App) GetNoteContent(filePath string) (*backend.Note, error) {
	return backend.ParseMarkdownFile(filePath, filepath.Dir(filePath))
}

// SaveNoteContent lưu lại nội dung và thông tin note xuống ổ cứng
func (a *App) SaveNoteContent(filePath string, content string, title string, tags []string) error {
	err := backend.SaveMarkdownFile(filePath, content, title, tags)
	if err == nil {
		runtime.EventsEmit(a.ctx, "on-file-change", filePath)
	}
	return err
}

// CreateFolder tạo một thư mục con trong workspace
func (a *App) CreateFolder(parentPath string, folderName string) error {
	a.mu.Lock()
	ws := a.workspacePath
	a.mu.Unlock()

	targetParent := parentPath
	if targetParent == "" || targetParent == "root" {
		targetParent = ws
	}
	if targetParent == "" {
		return fmt.Errorf("chưa có workspace để tạo thư mục")
	}

	newPath := filepath.Join(targetParent, folderName)
	err := backend.CreateNewFolder(newPath)
	if err == nil && a.watcher != nil {
		_ = a.watcher.Add(newPath)
		runtime.EventsEmit(a.ctx, "on-file-change", newPath)
	}
	return err
}

// CreateNote tạo file markdown .md mới và trả về object Note
func (a *App) CreateNote(folderPath string, noteTitle string) (*backend.Note, error) {
	a.mu.Lock()
	ws := a.workspacePath
	a.mu.Unlock()

	targetFolder := folderPath
	if targetFolder == "" || targetFolder == "root" {
		targetFolder = ws
	}
	if targetFolder == "" {
		return nil, fmt.Errorf("chưa mở workspace nào")
	}

	fileName := strings.ToLower(strings.ReplaceAll(noteTitle, " ", "-"))
	if !strings.HasSuffix(fileName, ".md") {
		fileName += ".md"
	}
	filePath := filepath.Join(targetFolder, fileName)

	err := backend.CreateNewNote(filePath, noteTitle)
	if err != nil {
		return nil, err
	}

	runtime.EventsEmit(a.ctx, "on-file-change", filePath)
	return backend.ParseMarkdownFile(filePath, targetFolder)
}

// SearchNotes tìm kiếm từ khoá trong title hoặc nội dung note thuộc workspace
func (a *App) SearchNotes(query string) ([]backend.Note, error) {
	data, err := a.ScanFolder()
	if err != nil || data == nil {
		return nil, err
	}

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return data.Notes, nil
	}

	var results []backend.Note
	for _, n := range data.Notes {
		if strings.Contains(strings.ToLower(n.Title), q) || strings.Contains(strings.ToLower(n.Content), q) {
			results = append(results, n)
		}
	}
	return results, nil
}
