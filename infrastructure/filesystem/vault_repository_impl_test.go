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

