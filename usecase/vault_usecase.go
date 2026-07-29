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
	mu            sync.Mutex
	onChange      func(path string)
}

func NewVaultUseCase(repo repository.VaultRepository) VaultUseCase {
	return &vaultUseCaseImpl{
		repo: repo,
	}
}

func (u *vaultUseCaseImpl) SetOnChangeCallback(callback func(path string)) {
	u.mu.Lock()
	u.onChange = callback
	u.mu.Unlock()
	_ = u.repo.StartWatcher(u.currentWsPath, callback)
}

func (u *vaultUseCaseImpl) OpenWorkspace(dirPath string) (*model.WorkspaceData, error) {
	u.mu.Lock()
	u.currentWsPath = dirPath
	cb := u.onChange
	u.mu.Unlock()

	if cb != nil && dirPath != "" {
		_ = u.repo.StartWatcher(dirPath, cb)
	}
	return u.repo.ScanWorkspace(dirPath)
}

func (u *vaultUseCaseImpl) GetCurrentWorkspace() (*model.WorkspaceData, error) {
	u.mu.Lock()
	path := u.currentWsPath
	u.mu.Unlock()
	if path == "" {
		return nil, fmt.Errorf("chưa có workspace nào được thiết lập")
	}
	return u.repo.ScanWorkspace(path)
}

func (u *vaultUseCaseImpl) GetNote(filePath string) (*model.Note, error) {
	return u.repo.ParseNote(filePath, filepath.Dir(filePath))
}

func (u *vaultUseCaseImpl) SaveNote(filePath string, content string, title string, tags []string) error {
	return u.repo.SaveNote(filePath, content, title, tags)
}

func (u *vaultUseCaseImpl) CreateFolder(parentPath string, folderName string) error {
	u.mu.Lock()
	ws := u.currentWsPath
	u.mu.Unlock()

	targetParent := parentPath
	if targetParent == "" || targetParent == "root" {
		targetParent = ws
	}
	if targetParent == "" {
		return fmt.Errorf("chưa có workspace để tạo thư mục")
	}
	return u.repo.CreateFolder(targetParent, folderName)
}

func (u *vaultUseCaseImpl) CreateNote(folderPath string, noteTitle string) (*model.Note, error) {
	u.mu.Lock()
	ws := u.currentWsPath
	u.mu.Unlock()

	targetFolder := folderPath
	if targetFolder == "" || targetFolder == "root" {
		targetFolder = ws
	}
	if targetFolder == "" {
		return nil, fmt.Errorf("chưa có workspace để tạo ghi chú")
	}
	return u.repo.CreateNote(targetFolder, noteTitle)
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
		if strings.Contains(strings.ToLower(n.Title), q) || strings.Contains(strings.ToLower(n.Content), q) {
			results = append(results, n)
		}
	}
	return results, nil
}
