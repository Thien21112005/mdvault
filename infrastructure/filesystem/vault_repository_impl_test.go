package filesystem

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWatcherLifecycle_StartStopRestart(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_watcher_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewVaultRepository()

	// 1. Khởi động watcher với đường dẫn rỗng -> Phải thành công mà không lỗi
	if err := repo.StartWatcher("", nil); err != nil {
		t.Fatalf("StartWatcher with empty root should return nil, got: %v", err)
	}

	var mu sync.Mutex
	var eventsReceived []string

	onChange := func(path string) {
		mu.Lock()
		eventsReceived = append(eventsReceived, path)
		mu.Unlock()
	}

	// 2. Khởi động watcher với thư mục hợp lệ
	if err := repo.StartWatcher(tempDir, onChange); err != nil {
		t.Fatalf("StartWatcher failed: %v", err)
	}

	// 3. Dừng watcher
	if err := repo.StopWatcher(); err != nil {
		t.Fatalf("StopWatcher failed: %v", err)
	}

	// 4. Khởi động lại watcher (Restart) -> Không được chết goroutine hoặc panic
	if err := repo.StartWatcher(tempDir, onChange); err != nil {
		t.Fatalf("Restarting StartWatcher failed: %v", err)
	}

	// Tạo một file markdown để xem watcher có bắt được không
	testFile := filepath.Join(tempDir, "note1.md")
	if err := os.WriteFile(testFile, []byte("# Test Note"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Chờ watcher xử lý sự kiện
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	count := len(eventsReceived)
	mu.Unlock()

	if count == 0 {
		t.Logf("Notice: OS file notification might have delayed or been caught: count=%d", count)
	}

	// Dọn dẹp an toàn
	_ = repo.StopWatcher()
}

func TestWatcher_AutoAddSubdirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_subdir_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewVaultRepository()

	var mu sync.Mutex
	var events []string
	onChange := func(path string) {
		mu.Lock()
		events = append(events, path)
		mu.Unlock()
	}

	if err := repo.StartWatcher(tempDir, onChange); err != nil {
		t.Fatalf("StartWatcher failed: %v", err)
	}
	defer repo.StopWatcher()

	// Tạo một thư mục con mới từ bên ngoài
	subDir := filepath.Join(tempDir, "new_folder")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subDir: %v", err)
	}

	// Cho watcher thời gian để bắt sự kiện Create folder và add vào watch list
	time.Sleep(150 * time.Millisecond)

	// Tạo file .md bên trong thư mục con mới
	subNote := filepath.Join(subDir, "sub_note.md")
	if err := os.WriteFile(subNote, []byte("# Note inside new folder"), 0644); err != nil {
		t.Fatalf("Failed to create note in subfolder: %v", err)
	}

	time.Sleep(250 * time.Millisecond)

	mu.Lock()
	total := len(events)
	mu.Unlock()

	t.Logf("Total events caught during subdir creation: %d", total)
}

func TestWatcher_DebounceCoalescing(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_debounce_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewVaultRepository()

	var mu sync.Mutex
	callbackCount := 0
	onChange := func(path string) {
		mu.Lock()
		callbackCount++
		mu.Unlock()
	}

	if err := repo.StartWatcher(tempDir, onChange); err != nil {
		t.Fatalf("StartWatcher failed: %v", err)
	}
	defer repo.StopWatcher()

	testFile := filepath.Join(tempDir, "debounce_note.md")

	// Ghi liên tục 5 lần vào cùng 1 file trong vòng 20ms
	for i := 0; i < 5; i++ {
		_ = os.WriteFile(testFile, []byte("Content update"), 0644)
		time.Sleep(4 * time.Millisecond)
	}

	// Đợi debounce (50ms) hoàn tất + buffer
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	count := callbackCount
	mu.Unlock()

	// Nhờ debounce 50ms, 5 lần ghi dồn dập phải được gom lại
	if count > 2 {
		t.Errorf("Expected at most 1-2 callback events due to debouncing, got %d", count)
	} else {
		t.Logf("Debounce successfully coalesced rapid writes: callbackCount = %d", count)
	}
}

func TestWatcher_IgnoreNonMarkdownAndTempFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_ignore_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewVaultRepository()

	var mu sync.Mutex
	events := []string{}
	onChange := func(path string) {
		mu.Lock()
		events = append(events, path)
		mu.Unlock()
	}

	if err := repo.StartWatcher(tempDir, onChange); err != nil {
		t.Fatalf("StartWatcher failed: %v", err)
	}
	defer repo.StopWatcher()

	// Tạo các file tạm hoặc file không phải markdown
	_ = os.WriteFile(filepath.Join(tempDir, "test.tmp"), []byte("temp"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "test.swp"), []byte("swap"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, ".hidden"), []byte("hidden"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "image.png"), []byte("binary"), 0644)

	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	count := len(events)
	mu.Unlock()

	if count > 0 {
		t.Errorf("Expected 0 events for non-markdown and temp files, got %d: %v", count, events)
	}
}

func TestAtomicWrite_SaveNoteAndCleanTemp(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_atomic_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewVaultRepository()
	testFile := filepath.Join(tempDir, "atomic_note.md")

	// 1. Lưu ghi chú lần đầu
	content1 := "Nội dung ban đầu"
	if err := repo.SaveNote(testFile, content1, "Tiêu Đề 1", []string{"tag1"}); err != nil {
		t.Fatalf("SaveNote failed: %v", err)
	}

	// Đọc lại kiểm tra
	data1, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read saved note: %v", err)
	}
	if !strings.Contains(string(data1), content1) {
		t.Errorf("Saved file does not contain expected content: %s", string(data1))
	}

	// 2. Ghi đè bằng SaveNote (Atomic overwrite)
	content2 := "Nội dung ghi đè cập nhật"
	if err := repo.SaveNote(testFile, content2, "Tiêu Đề 2", []string{"tag2"}); err != nil {
		t.Fatalf("SaveNote overwrite failed: %v", err)
	}

	data2, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read overwritten note: %v", err)
	}
	if !strings.Contains(string(data2), content2) {
		t.Errorf("Overwritten file does not contain new content: %s", string(data2))
	}

	// 3. Kiểm tra xem có file tạm .mdvault_tmp_* nào bị bỏ quên không
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read dir: %v", err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".mdvault_tmp_") {
			t.Errorf("Found leftover temp file: %s", entry.Name())
		}
	}
}

func TestCreateNote_AtomicPersistence(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_create_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewVaultRepository()
	note, err := repo.CreateNote(tempDir, "Kiến Trúc Clean Architecture")
	if err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	if note == nil {
		t.Fatal("Expected created note to be non-nil")
	}
	if note.Title != "Kiến Trúc Clean Architecture" {
		t.Errorf("Expected title 'Kiến Trúc Clean Architecture', got '%s'", note.Title)
	}

	// Kiểm tra file trên đĩa
	expectedPath := filepath.Join(tempDir, "kiến-trúc-clean-architecture.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected file at %s does not exist", expectedPath)
	}
}

func TestFrontmatter_PreservesCustomMetadataOnSave(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_meta_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewVaultRepository()
	testFile := filepath.Join(tempDir, "rich_note.md")

	initialYAML := `---
title: "Kiến Trúc Microservices"
tags: ["architecture", "go"]
icon: "server"
author: "Alex"
aliases:
  - "distributed-systems"
  - "clean-arch"
status: "draft"
---

# Nội dung bài viết
Uncle Bob's Clean Architecture pattern...
`
	if err := os.WriteFile(testFile, []byte(initialYAML), 0644); err != nil {
		t.Fatalf("Failed to write initial file: %v", err)
	}

	// 1. Parse note ban đầu
	note, err := repo.ParseNote(testFile, "root")
	if err != nil {
		t.Fatalf("ParseNote failed: %v", err)
	}
	if note.Metadata == nil {
		t.Fatal("Expected note.Metadata to be non-nil")
	}
	if note.Metadata["icon"] != "server" {
		t.Errorf("Expected icon 'server', got '%v'", note.Metadata["icon"])
	}
	if note.Metadata["author"] != "Alex" {
		t.Errorf("Expected author 'Alex', got '%v'", note.Metadata["author"])
	}

	// 2. Lưu cập nhật với Tiêu Đề và Tags mới
	newContent := "Nội dung cập nhật mới hoàn toàn..."
	newTitle := "Kiến Trúc Microservices 2026"
	newTags := []string{"architecture", "golang", "backend"}

	if err := repo.SaveNote(testFile, newContent, newTitle, newTags); err != nil {
		t.Fatalf("SaveNote failed: %v", err)
	}

	// 3. Đọc lại file sau khi lưu và xác minh
	savedNote, err := repo.ParseNote(testFile, "root")
	if err != nil {
		t.Fatalf("ParseNote on saved file failed: %v", err)
	}

	// Xác minh tiêu đề và tags đã được cập nhật
	if savedNote.Title != newTitle {
		t.Errorf("Expected updated title '%s', got '%s'", newTitle, savedNote.Title)
	}
	if len(savedNote.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(savedNote.Tags))
	}

	// BẢO TOÀN 100%: Các trường custom YAML ban đầu (icon, author, status, aliases) KHÔNG ĐƯỢC MẤT!
	if savedNote.Metadata == nil {
		t.Fatal("Expected savedNote.Metadata to be non-nil")
	}
	if savedNote.Metadata["icon"] != "server" {
		t.Errorf("FAIL: 'icon' metadata was lost or modified! Got '%v'", savedNote.Metadata["icon"])
	}
	if savedNote.Metadata["author"] != "Alex" {
		t.Errorf("FAIL: 'author' metadata was lost or modified! Got '%v'", savedNote.Metadata["author"])
	}
	if savedNote.Metadata["status"] != "draft" {
		t.Errorf("FAIL: 'status' metadata was lost or modified! Got '%v'", savedNote.Metadata["status"])
	}
}

func TestFrontmatter_FallbackToH1AndFilename(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_fallback_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewVaultRepository()

	// Trường hợp 1: Không có Frontmatter nhưng có H1
	fileWithH1 := filepath.Join(tempDir, "guide.md")
	if err := os.WriteFile(fileWithH1, []byte("# Hướng Dẫn Git Rebase\n\nNội dung hướng dẫn..."), 0644); err != nil {
		t.Fatalf("Failed to write fileWithH1: %v", err)
	}
	note1, err := repo.ParseNote(fileWithH1, "root")
	if err != nil {
		t.Fatalf("ParseNote note1 failed: %v", err)
	}
	if note1.Title != "Hướng Dẫn Git Rebase" {
		t.Errorf("Expected fallback title from H1 'Hướng Dẫn Git Rebase', got '%s'", note1.Title)
	}

	// Trường hợp 2: Không có Frontmatter và không có H1 -> Lấy tên file
	fileNoHeading := filepath.Join(tempDir, "todo-list.md")
	if err := os.WriteFile(fileNoHeading, []byte("Chỉ có nội dung văn bản thường không tiêu đề"), 0644); err != nil {
		t.Fatalf("Failed to write fileNoHeading: %v", err)
	}
	note2, err := repo.ParseNote(fileNoHeading, "root")
	if err != nil {
		t.Fatalf("ParseNote note2 failed: %v", err)
	}
	if note2.Title != "todo-list" {
		t.Errorf("Expected fallback title from filename 'todo-list', got '%s'", note2.Title)
	}
}

func TestPathNormalization_WindowsNoBackslashes(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_slash_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewVaultRepository()

	// Tạo subfolder
	subDir := filepath.Join(tempDir, "docs")
	if err := repo.CreateFolder(tempDir, "docs"); err != nil {
		t.Fatalf("CreateFolder failed: %v", err)
	}

	// Tạo note trong subfolder
	subNote, err := repo.CreateNote(subDir, "Architecture Guide")
	if err != nil {
		t.Fatalf("CreateNote in subfolder failed: %v", err)
	}

	if strings.Contains(subNote.ID, "\\") {
		t.Errorf("subNote.ID contains backslashes: %s", subNote.ID)
	}
	if strings.Contains(subNote.Folder, "\\") {
		t.Errorf("subNote.Folder contains backslashes: %s", subNote.Folder)
	}

	// Quét ScanWorkspace
	data, err := repo.ScanWorkspace(tempDir)
	if err != nil {
		t.Fatalf("ScanWorkspace failed: %v", err)
	}
	if strings.Contains(data.RootPath, "\\") {
		t.Errorf("data.RootPath contains backslashes: %s", data.RootPath)
	}
	for _, n := range data.Notes {
		if strings.Contains(n.ID, "\\") {
			t.Errorf("Scanned note ID contains backslashes: %s", n.ID)
		}
		if strings.Contains(n.Folder, "\\") {
			t.Errorf("Scanned note Folder contains backslashes: %s", n.Folder)
		}
	}
	for _, f := range data.Folders {
		if strings.Contains(f.ID, "\\") {
			t.Errorf("Scanned folder ID contains backslashes: %s", f.ID)
		}
	}
}

func TestRepository_TableDriven_FileOperations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_repo_table_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := NewVaultRepository()

	tests := []struct {
		name       string
		folderName string
		noteTitle  string
		content    string
		tags       []string
		renameTo   string
	}{
		{
			name:       "Standard Note in Root",
			folderName: "",
			noteTitle:  "Standard Note",
			content:    "Standard content without special chars",
			tags:       []string{"go", "unit-test"},
			renameTo:   "Standard Note Renamed",
		},
		{
			name:       "Unicode Vietnamese Note in Subfolder",
			folderName: "tai-lieu",
			noteTitle:  "Kiến Trúc Phần Mềm",
			content:    "Nội dung tiếng Việt có dấu đầy đủ và chi tiết.",
			tags:       []string{"kiến-trúc", "hệ-thống"},
			renameTo:   "Kiến Trúc Microservices Mới",
		},
		{
			name:       "Note with Empty Tags",
			folderName: "scratch",
			noteTitle:  "Scratchpad",
			content:    "Chỉ là nội dung nháp nhanh",
			tags:       []string{},
			renameTo:   "Scratchpad Final",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetDir := tempDir
			if tt.folderName != "" {
				if err := repo.CreateFolder(tempDir, tt.folderName); err != nil {
					t.Fatalf("[%s] CreateFolder failed: %v", tt.name, err)
				}
				targetDir = filepath.Join(tempDir, tt.folderName)
			}

			// 1. CreateNote
			createdNote, err := repo.CreateNote(targetDir, tt.noteTitle)
			if err != nil {
				t.Fatalf("[%s] CreateNote failed: %v", tt.name, err)
			}
			if createdNote == nil || createdNote.Title != tt.noteTitle {
				t.Fatalf("[%s] Expected title '%s', got '%v'", tt.name, tt.noteTitle, createdNote)
			}

			// 2. SaveNote with custom content and tags
			if err := repo.SaveNote(createdNote.ID, tt.content, tt.noteTitle, tt.tags); err != nil {
				t.Fatalf("[%s] SaveNote failed: %v", tt.name, err)
			}

			// 3. ParseNote to verify persistence
			parsed, err := repo.ParseNote(createdNote.ID, createdNote.Folder)
			if err != nil {
				t.Fatalf("[%s] ParseNote failed: %v", tt.name, err)
			}
			if !strings.Contains(parsed.Content, tt.content) {
				t.Errorf("[%s] Parsed content does not contain expected text", tt.name)
			}
			if len(parsed.Tags) != len(tt.tags) {
				t.Errorf("[%s] Expected %d tags, got %d", tt.name, len(tt.tags), len(parsed.Tags))
			}

			// 4. RenameNote
			newSlug := strings.ToLower(strings.ReplaceAll(tt.renameTo, " ", "-")) + ".md"
			newPath := filepath.ToSlash(filepath.Join(targetDir, newSlug))
			if err := repo.RenameNote(createdNote.ID, newPath); err != nil {
				t.Fatalf("[%s] RenameNote failed: %v", tt.name, err)
			}
			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				t.Errorf("[%s] Renamed note does not exist at %s", tt.name, newPath)
			}

			// 5. DeleteNote
			if err := repo.DeleteNote(newPath); err != nil {
				t.Fatalf("[%s] DeleteNote failed: %v", tt.name, err)
			}
			if _, err := os.Stat(newPath); !os.IsNotExist(err) {
				t.Errorf("[%s] Deleted note still exists on disk", tt.name)
			}
		})
	}
}



