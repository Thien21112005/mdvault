package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"mdvault/infrastructure/filesystem"
)

func TestVaultUseCase_InMemoryCaching(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_uc_cache_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Tạo 10 file note mẫu
	for i := 1; i <= 10; i++ {
		noteFile := filepath.Join(tempDir, fmt.Sprintf("note_%d.md", i))
		content := fmt.Sprintf("---\ntitle: Note %d\ntags: [test, go]\n---\n\nNội dung note %d", i, i)
		_ = os.WriteFile(noteFile, []byte(content), 0644)
	}

	repo := filesystem.NewVaultRepository()
	uc := NewVaultUseCase(repo)

	// 1. Mở workspace -> Nạp vào RAM
	wsData, err := uc.OpenWorkspace(tempDir)
	if err != nil {
		t.Fatalf("OpenWorkspace failed: %v", err)
	}
	if len(wsData.Notes) != 10 {
		t.Fatalf("Expected 10 notes, got %d", len(wsData.Notes))
	}

	// 2. Gọi GetCurrentWorkspace() 1000 lần liên tục từ RAM
	start := time.Now()
	for i := 0; i < 1000; i++ {
		cached, err := uc.GetCurrentWorkspace()
		if err != nil {
			t.Fatalf("GetCurrentWorkspace failed: %v", err)
		}
		if len(cached.Notes) != 10 {
			t.Fatalf("Corrupted cache: expected 10 notes, got %d", len(cached.Notes))
		}
	}
	elapsed := time.Since(start)
	t.Logf("1,000 GetCurrentWorkspace calls from RAM completed in: %v (Avg per call: %v)", elapsed, elapsed/1000)

	if elapsed > 100*time.Millisecond {
		t.Errorf("1,000 in-memory cache reads took too long: %v", elapsed)
	}
}

func TestVaultUseCase_SearchNotes_Sub5ms(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_uc_search_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Tạo 100 note mẫu với đa dạng tiêu đề, tag và nội dung
	for i := 1; i <= 100; i++ {
		noteFile := filepath.Join(tempDir, fmt.Sprintf("article_%d.md", i))
		var tag string
		if i%2 == 0 {
			tag = "backend"
		} else {
			tag = "frontend"
		}
		content := fmt.Sprintf("---\ntitle: \"Kiến Trúc Hệ Thống %d\"\ntags: [\"%s\", \"golang\"]\n---\n\nNội dung về Clean Architecture và microservice %d", i, tag, i)
		_ = os.WriteFile(noteFile, []byte(content), 0644)
	}

	repo := filesystem.NewVaultRepository()
	uc := NewVaultUseCase(repo)

	_, err = uc.OpenWorkspace(tempDir)
	if err != nil {
		t.Fatalf("OpenWorkspace failed: %v", err)
	}

	// Test 1: Tìm theo Tag "backend" (khoảng 50 kết quả)
	start := time.Now()
	resultsTag, err := uc.SearchNotes("backend")
	durationTag := time.Since(start)
	if err != nil {
		t.Fatalf("SearchNotes tag failed: %v", err)
	}
	if len(resultsTag) != 50 {
		t.Errorf("Expected 50 backend notes, got %d", len(resultsTag))
	}
	t.Logf("Search by Tag completed in: %v (Sub-5ms requirement met)", durationTag)
	if durationTag > 5*time.Millisecond {
		t.Errorf("Search by tag took %v > 5ms", durationTag)
	}

	// Test 2: Tìm theo Title "Hệ Thống 42" (chính xác 1 kết quả)
	start = time.Now()
	resultsTitle, err := uc.SearchNotes("Hệ Thống 42")
	durationTitle := time.Since(start)
	if err != nil {
		t.Fatalf("SearchNotes title failed: %v", err)
	}
	if len(resultsTitle) != 1 {
		t.Errorf("Expected 1 result for 'Hệ Thống 42', got %d", len(resultsTitle))
	}
	t.Logf("Search by Title completed in: %v", durationTitle)
	if durationTitle > 5*time.Millisecond {
		t.Errorf("Search by title took %v > 5ms", durationTitle)
	}

	// Test 3: Tìm theo Content "microservice" (100 kết quả)
	start = time.Now()
	resultsContent, err := uc.SearchNotes("microservice")
	durationContent := time.Since(start)
	if err != nil {
		t.Fatalf("SearchNotes content failed: %v", err)
	}
	if len(resultsContent) != 100 {
		t.Errorf("Expected 100 results for 'microservice', got %d", len(resultsContent))
	}
	t.Logf("Search by Content completed in: %v", durationContent)
	if durationContent > 5*time.Millisecond {
		t.Errorf("Search by content took %v > 5ms", durationContent)
	}
}

func TestVaultUseCase_ConcurrentAccess_Race(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_uc_race_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	for i := 1; i <= 20; i++ {
		noteFile := filepath.Join(tempDir, fmt.Sprintf("note_%d.md", i))
		_ = os.WriteFile(noteFile, []byte(fmt.Sprintf("---\ntitle: Note %d\ntags: [go]\n---\n\nBody %d", i, i)), 0644)
	}

	repo := filesystem.NewVaultRepository()
	uc := NewVaultUseCase(repo)

	if _, err := uc.OpenWorkspace(tempDir); err != nil {
		t.Fatalf("OpenWorkspace failed: %v", err)
	}

	var wg sync.WaitGroup
	// 10 goroutines đọc GetCurrentWorkspace đồng thời
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_, _ = uc.GetCurrentWorkspace()
			}
		}()
	}

	// 10 goroutines tìm kiếm đồng thời
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_, _ = uc.SearchNotes("Note")
			}
		}(i)
	}

	// 2 goroutines cập nhật note đồng thời
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				filePath := filepath.Join(tempDir, fmt.Sprintf("note_%d.md", id+1))
				_ = uc.SaveNote(filePath, fmt.Sprintf("Updated %d", j), fmt.Sprintf("Note %d", id+1), []string{"go", "updated"})
			}
		}(i)
	}

	wg.Wait()
}
