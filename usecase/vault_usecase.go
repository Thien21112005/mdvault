package usecase

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"mdvault/domain/model"
	"mdvault/domain/repository"
)

type VaultUseCase interface {
	OpenWorkspace(dirPath string) (*model.WorkspaceData, error)
	GetCurrentWorkspace() (*model.WorkspaceData, error)
	RefreshCache() (*model.WorkspaceData, error)
	GetNote(filePath string) (*model.Note, error)
	SaveNote(filePath string, content string, title string, tags []string) error
	CreateFolder(parentPath string, folderName string) error
	CreateNote(folderPath string, noteTitle string) (*model.Note, error)
	SearchNotes(query string) ([]model.Note, error)
	SetOnChangeCallback(callback func(path string))
}

type vaultUseCaseImpl struct {
	repo          repository.VaultRepository
	currentWsPath string
	cachedData    *model.WorkspaceData
	mu            sync.RWMutex
	onChange      func(path string)
}

func NewVaultUseCase(repo repository.VaultRepository) VaultUseCase {
	return &vaultUseCaseImpl{
		repo: repo,
	}
}

func (u *vaultUseCaseImpl) SetOnChangeCallback(callback func(path string)) {
	wrapped := func(changedPath string) {
		u.mu.RLock()
		ws := u.currentWsPath
		u.mu.RUnlock()

		if ws != "" {
			// Làm tươi cache RAM ngầm ngay khi OS báo file thay đổi
			if freshData, err := u.repo.ScanWorkspace(ws); err == nil && freshData != nil {
				u.mu.Lock()
				u.cachedData = freshData
				u.mu.Unlock()
			}
		}

		if callback != nil {
			callback(changedPath)
		}
	}

	u.mu.Lock()
	u.onChange = wrapped
	ws := u.currentWsPath
	u.mu.Unlock()

	if ws != "" {
		_ = u.repo.StartWatcher(ws, wrapped)
	}
}

func (u *vaultUseCaseImpl) OpenWorkspace(dirPath string) (*model.WorkspaceData, error) {
	trimmed := strings.TrimSpace(dirPath)
	if trimmed == "" {
		return nil, fmt.Errorf("đường dẫn workspace không được để trống")
	}

	absPath, err := filepath.Abs(trimmed)
	if err != nil {
		absPath = trimmed
	}
	cleanPath := filepath.ToSlash(filepath.Clean(absPath))

	data, err := u.repo.ScanWorkspace(cleanPath)
	if err != nil {
		return nil, err
	}

	u.mu.Lock()
	u.currentWsPath = cleanPath
	u.cachedData = data
	cb := u.onChange
	u.mu.Unlock()

	if cb != nil {
		_ = u.repo.StartWatcher(cleanPath, cb)
	}
	return data, nil
}

func (u *vaultUseCaseImpl) GetCurrentWorkspace() (*model.WorkspaceData, error) {
	u.mu.RLock()
	ws := u.currentWsPath
	cached := u.cachedData
	u.mu.RUnlock()

	if ws == "" {
		return nil, fmt.Errorf("chưa có workspace nào được thiết lập")
	}

	if cached != nil {
		return cached, nil
	}

	return u.RefreshCache()
}

func (u *vaultUseCaseImpl) RefreshCache() (*model.WorkspaceData, error) {
	u.mu.RLock()
	ws := u.currentWsPath
	u.mu.RUnlock()

	if ws == "" {
		return nil, fmt.Errorf("chưa có workspace nào được thiết lập")
	}

	data, err := u.repo.ScanWorkspace(ws)
	if err != nil {
		return nil, err
	}

	u.mu.Lock()
	u.cachedData = data
	u.mu.Unlock()

	return data, nil
}

func (u *vaultUseCaseImpl) GetNote(filePath string) (*model.Note, error) {
	u.mu.RLock()
	ws := u.currentWsPath
	u.mu.RUnlock()

	cleanFilePath := filepath.ToSlash(filepath.Clean(filePath))
	parentDir := filepath.ToSlash(filepath.Dir(cleanFilePath))

	folderID := parentDir
	if ws != "" && (parentDir == ws || strings.EqualFold(parentDir, ws)) {
		folderID = "root"
	}

	note, err := u.repo.ParseNote(cleanFilePath, folderID)
	if err != nil {
		return nil, err
	}
	note.ID = cleanFilePath
	note.Folder = folderID
	return note, nil
}

func (u *vaultUseCaseImpl) SaveNote(filePath string, content string, title string, tags []string) error {
	cleanPath := filepath.ToSlash(filepath.Clean(filePath))
	err := u.repo.SaveNote(cleanPath, content, title, tags)
	if err == nil {
		_, _ = u.RefreshCache()
	}
	return err
}

func (u *vaultUseCaseImpl) CreateFolder(parentPath string, folderName string) error {
	u.mu.RLock()
	ws := u.currentWsPath
	u.mu.RUnlock()

	cleanParent := filepath.ToSlash(filepath.Clean(parentPath))
	targetParent := cleanParent
	if targetParent == "" || targetParent == "root" || targetParent == "." || (ws != "" && strings.EqualFold(targetParent, ws)) {
		targetParent = ws
	}
	if targetParent == "" {
		return fmt.Errorf("chưa có workspace để tạo thư mục")
	}

	err := u.repo.CreateFolder(targetParent, folderName)
	if err == nil {
		_, _ = u.RefreshCache()
	}
	return err
}

func (u *vaultUseCaseImpl) CreateNote(folderPath string, noteTitle string) (*model.Note, error) {
	u.mu.RLock()
	ws := u.currentWsPath
	u.mu.RUnlock()

	if ws == "" {
		return nil, fmt.Errorf("chưa có workspace để tạo ghi chú")
	}

	cleanFolder := filepath.ToSlash(filepath.Clean(folderPath))
	targetFolder := cleanFolder
	isRoot := false

	if targetFolder == "" || targetFolder == "root" || targetFolder == "." || strings.EqualFold(targetFolder, ws) {
		targetFolder = ws
		isRoot = true
	}

	note, err := u.repo.CreateNote(targetFolder, noteTitle)
	if err != nil {
		return nil, err
	}

	if isRoot {
		note.Folder = "root"
	} else {
		note.Folder = filepath.ToSlash(targetFolder)
	}
	note.ID = filepath.ToSlash(note.ID)

	_, _ = u.RefreshCache()
	return note, nil
}

func (u *vaultUseCaseImpl) SearchNotes(query string) ([]model.Note, error) {
	data, err := u.GetCurrentWorkspace()
	if err != nil || data == nil {
		return nil, err
	}

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return data.Notes, nil
	}

	var results []model.Note
	for _, n := range data.Notes {
		tagMatches := false
		for _, t := range n.Tags {
			if strings.Contains(strings.ToLower(t), q) {
				tagMatches = true
				break
			}
		}
		if strings.Contains(strings.ToLower(n.Title), q) || strings.Contains(strings.ToLower(n.Content), q) || tagMatches {
			results = append(results, n)
		}
	}
	return results, nil
}
