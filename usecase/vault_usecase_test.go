package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func TestVaultUseCase_PathNormalization_And_RootFolder(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_uc_path_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := filesystem.NewVaultRepository()
	uc := NewVaultUseCase(repo)

	// 1. OpenWorkspace
	wsData, err := uc.OpenWorkspace(tempDir)
	if err != nil {
		t.Fatalf("OpenWorkspace failed: %v", err)
	}
	if strings.Contains(wsData.RootPath, "\\") {
		t.Errorf("Workspace RootPath contains backslashes: %s", wsData.RootPath)
	}

	// 2. Create note in "root"
	rootNote1, err := uc.CreateNote("root", "Root Note Alpha")
	if err != nil {
		t.Fatalf("CreateNote in root failed: %v", err)
	}
	if rootNote1.Folder != "root" {
		t.Errorf("Expected rootNote1.Folder to be 'root', got '%s'", rootNote1.Folder)
	}
	if strings.Contains(rootNote1.ID, "\\") {
		t.Errorf("rootNote1.ID contains backslashes: %s", rootNote1.ID)
	}

	// 3. Create note with empty folder path ("") -> Should also resolve to "root"
	rootNote2, err := uc.CreateNote("", "Root Note Beta")
	if err != nil {
		t.Fatalf("CreateNote with empty folder failed: %v", err)
	}
	if rootNote2.Folder != "root" {
		t.Errorf("Expected rootNote2.Folder to be 'root', got '%s'", rootNote2.Folder)
	}
	if strings.Contains(rootNote2.ID, "\\") {
		t.Errorf("rootNote2.ID contains backslashes: %s", rootNote2.ID)
	}

	// 4. Create subfolder
	if err := uc.CreateFolder("root", "projects"); err != nil {
		t.Fatalf("CreateFolder projects failed: %v", err)
	}

	// 5. Create note in subfolder
	subDirPath := filepath.Join(tempDir, "projects")
	subNote, err := uc.CreateNote(subDirPath, "Project Spec")
	if err != nil {
		t.Fatalf("CreateNote in subfolder failed: %v", err)
	}
	expectedSubFolder := filepath.ToSlash(filepath.Clean(subDirPath))
	if subNote.Folder != expectedSubFolder {
		t.Errorf("Expected subNote.Folder to be '%s', got '%s'", expectedSubFolder, subNote.Folder)
	}
	if strings.Contains(subNote.Folder, "\\") {
		t.Errorf("subNote.Folder contains backslashes: %s", subNote.Folder)
	}
	if strings.Contains(subNote.ID, "\\") {
		t.Errorf("subNote.ID contains backslashes: %s", subNote.ID)
	}

	// 6. GetNote on root note
	retrievedRoot, err := uc.GetNote(rootNote1.ID)
	if err != nil {
		t.Fatalf("GetNote rootNote1 failed: %v", err)
	}
	if retrievedRoot.Folder != "root" {
		t.Errorf("GetNote root note expected Folder 'root', got '%s'", retrievedRoot.Folder)
	}
	if strings.Contains(retrievedRoot.ID, "\\") || strings.Contains(retrievedRoot.Folder, "\\") {
		t.Errorf("retrievedRoot contains backslashes: ID=%s, Folder=%s", retrievedRoot.ID, retrievedRoot.Folder)
	}

	// 7. GetNote on subfolder note
	retrievedSub, err := uc.GetNote(subNote.ID)
	if err != nil {
		t.Fatalf("GetNote subNote failed: %v", err)
	}
	if retrievedSub.Folder != expectedSubFolder {
		t.Errorf("GetNote sub note expected Folder '%s', got '%s'", expectedSubFolder, retrievedSub.Folder)
	}
	if strings.Contains(retrievedSub.ID, "\\") || strings.Contains(retrievedSub.Folder, "\\") {
		t.Errorf("retrievedSub contains backslashes: ID=%s, Folder=%s", retrievedSub.ID, retrievedSub.Folder)
	}

	// 8. GetCurrentWorkspace cache consistency
	cached, err := uc.GetCurrentWorkspace()
	if err != nil {
		t.Fatalf("GetCurrentWorkspace failed: %v", err)
	}
	for _, n := range cached.Notes {
		if strings.Contains(n.ID, "\\") {
			t.Errorf("Cached note ID contains backslashes: %s", n.ID)
		}
		if strings.Contains(n.Folder, "\\") {
			t.Errorf("Cached note Folder contains backslashes: %s", n.Folder)
		}
		if n.Title == "Root Note Alpha" || n.Title == "Root Note Beta" {
			if n.Folder != "root" {
				t.Errorf("Cached note '%s' expected Folder 'root', got '%s'", n.Title, n.Folder)
			}
		}
		if n.Title == "Project Spec" {
			if n.Folder != expectedSubFolder {
				t.Errorf("Cached note '%s' expected Folder '%s', got '%s'", n.Title, expectedSubFolder, n.Folder)
			}
		}
	}
	for _, f := range cached.Folders {
		if strings.Contains(f.ID, "\\") {
			t.Errorf("Cached folder ID contains backslashes: %s", f.ID)
		}
	}
}

func TestVaultUseCase_PathTraversal_Rejected(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdvault_uc_sec_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo := filesystem.NewVaultRepository()
	uc := NewVaultUseCase(repo)

	if _, err := uc.OpenWorkspace(tempDir); err != nil {
		t.Fatalf("OpenWorkspace failed: %v", err)
	}

	// 1. GetNote with relative path traversal
	if _, err := uc.GetNote("../../system_file.md"); err == nil {
		t.Errorf("Expected GetNote with ../../ to fail with access denied error, but got nil")
	} else if !strings.Contains(err.Error(), "truy cập bị từ chối") {
		t.Errorf("Expected 'truy cập bị từ chối' error message, got: %v", err)
	}

	// 2. GetNote with absolute path traversal
	absOutside := filepath.Join(tempDir, "..", "..", "system_file.md")
	if _, err := uc.GetNote(absOutside); err == nil {
		t.Errorf("Expected GetNote with outside absolute path to fail, got nil")
	} else if !strings.Contains(err.Error(), "truy cập bị từ chối") {
		t.Errorf("Expected 'truy cập bị từ chối' error message, got: %v", err)
	}

	// 3. SaveNote with path traversal
	if err := uc.SaveNote("../../shadow.txt", "exploit", "Exploit", nil); err == nil {
		t.Errorf("Expected SaveNote with ../../ to fail, got nil")
	} else if !strings.Contains(err.Error(), "truy cập bị từ chối") {
		t.Errorf("Expected 'truy cập bị từ chối', got: %v", err)
	}

	// 4. CreateFolder outside workspace
	if err := uc.CreateFolder("../../outside_folder", "sub"); err == nil {
		t.Errorf("Expected CreateFolder with ../../ to fail, got nil")
	} else if !strings.Contains(err.Error(), "truy cập bị từ chối") {
		t.Errorf("Expected 'truy cập bị từ chối', got: %v", err)
	}

	// 5. CreateFolder with folderName containing traversal
	if err := uc.CreateFolder("root", "../../evil_folder"); err == nil {
		t.Errorf("Expected CreateFolder with malicious folderName to fail, got nil")
	} else if !strings.Contains(err.Error(), "truy cập bị từ chối") {
		t.Errorf("Expected 'truy cập bị từ chối', got: %v", err)
	}

	// 6. CreateNote outside workspace
	if _, err := uc.CreateNote("../../outside_folder", "Evil Note"); err == nil {
		t.Errorf("Expected CreateNote with ../../ to fail, got nil")
	} else if !strings.Contains(err.Error(), "truy cập bị từ chối") {
		t.Errorf("Expected 'truy cập bị từ chối', got: %v", err)
	}

	// 7. CreateNote with noteTitle containing traversal
	if _, err := uc.CreateNote("root", "../../Evil Note"); err == nil {
		t.Errorf("Expected CreateNote with malicious noteTitle to fail, got nil")
	} else if !strings.Contains(err.Error(), "truy cập bị từ chối") {
		t.Errorf("Expected 'truy cập bị từ chối', got: %v", err)
	}

	// 8. Thao tác hợp lệ bên trong workspace vẫn hoạt động hoàn hảo
	safeNote, err := uc.CreateNote("root", "Safe Note")
	if err != nil {
		t.Fatalf("CreateNote inside workspace failed: %v", err)
	}
	if safeNote.Folder != "root" {
		t.Errorf("Expected safeNote.Folder == 'root', got '%s'", safeNote.Folder)
	}

	retrieved, err := uc.GetNote(safeNote.ID)
	if err != nil {
		t.Fatalf("GetNote for valid note failed: %v", err)
	}
	if retrieved.Title != "Safe Note" {
		t.Errorf("Expected title 'Safe Note', got '%s'", retrieved.Title)
	}
}


