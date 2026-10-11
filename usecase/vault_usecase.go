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

func normalizeDriveLetter(p string) string {
	if len(p) >= 2 && p[1] == ':' {
		return strings.ToUpper(string(p[0])) + p[1:]
	}
	return p
}

func (u *vaultUseCaseImpl) validatePath(userPath string) (string, error) {
	u.mu.RLock()
	ws := u.currentWsPath
	u.mu.RUnlock()

	if strings.TrimSpace(ws) == "" {
		return "", fmt.Errorf("chưa có workspace nào được thiết lập")
	}

	cleanWs, err := filepath.Abs(ws)
	if err != nil {
		cleanWs = ws
	}
	cleanWs = normalizeDriveLetter(filepath.ToSlash(filepath.Clean(cleanWs)))

	trimmed := strings.TrimSpace(userPath)
	if trimmed == "" || trimmed == "root" || trimmed == "." {
		return cleanWs, nil
	}

	target := trimmed
	if !filepath.IsAbs(target) {
		target = filepath.Join(cleanWs, target)
	}

	absTarget, err := filepath.Abs(target)
	if err != nil {
		return "", fmt.Errorf("đường dẫn không hợp lệ: %w", err)
	}
	cleanTarget := normalizeDriveLetter(filepath.ToSlash(filepath.Clean(absTarget)))

	rel, err := filepath.Rel(filepath.FromSlash(cleanWs), filepath.FromSlash(cleanTarget))
	if err != nil {
		return "", fmt.Errorf("truy cập bị từ chối: đường dẫn vượt ra ngoài workspace (%s)", userPath)
	}

	relSlash := filepath.ToSlash(rel)
	if relSlash == ".." || strings.HasPrefix(relSlash, "../") {
		return "", fmt.Errorf("truy cập bị từ chối: đường dẫn vượt ra ngoài ranh giới workspace (%s)", userPath)
	}

	return cleanTarget, nil
}

func (u *vaultUseCaseImpl) GetNote(filePath string) (*model.Note, error) {
	validatedPath, err := u.validatePath(filePath)
	if err != nil {
		return nil, err
	}

	u.mu.RLock()
	ws := u.currentWsPath
	u.mu.RUnlock()

	cleanWs := normalizeDriveLetter(filepath.ToSlash(filepath.Clean(ws)))
	parentDir := normalizeDriveLetter(filepath.ToSlash(filepath.Dir(validatedPath)))

	folderID := parentDir
	if strings.EqualFold(parentDir, cleanWs) {
		folderID = "root"
	}

	note, err := u.repo.ParseNote(validatedPath, folderID)
	if err != nil {
		return nil, err
	}
	note.ID = validatedPath
	note.Folder = folderID
	return note, nil
}

func (u *vaultUseCaseImpl) SaveNote(filePath string, content string, title string, tags []string) error {
	validatedPath, err := u.validatePath(filePath)
	if err != nil {
		return err
	}

	err = u.repo.SaveNote(validatedPath, content, title, tags)
	if err == nil {
		_, _ = u.RefreshCache()
	}
	return err
}

func (u *vaultUseCaseImpl) CreateFolder(parentPath string, folderName string) error {
	trimmedName := strings.TrimSpace(folderName)
	if trimmedName == "" {
		return fmt.Errorf("tên thư mục không được để trống")
	}
	if strings.Contains(trimmedName, "..") || strings.Contains(trimmedName, "/") || strings.Contains(trimmedName, "\\") {
		return fmt.Errorf("truy cập bị từ chối: tên thư mục không được chứa ký tự phân cách hoặc '..'")
	}

	validatedParent, err := u.validatePath(parentPath)
	if err != nil {
		return err
	}

	targetDir := filepath.ToSlash(filepath.Join(validatedParent, trimmedName))
	if _, err := u.validatePath(targetDir); err != nil {
		return err
	}

	err = u.repo.CreateFolder(validatedParent, trimmedName)
	if err == nil {
		_, _ = u.RefreshCache()
	}
	return err
}

func (u *vaultUseCaseImpl) CreateNote(folderPath string, noteTitle string) (*model.Note, error) {
	trimmedTitle := strings.TrimSpace(noteTitle)
	if trimmedTitle == "" {
		return nil, fmt.Errorf("tiêu đề ghi chú không được để trống")
	}
	if strings.Contains(trimmedTitle, "..") || strings.Contains(trimmedTitle, "/") || strings.Contains(trimmedTitle, "\\") {
		return nil, fmt.Errorf("truy cập bị từ chối: tiêu đề ghi chú không được chứa ký tự phân cách hoặc '..'")
	}

	validatedFolder, err := u.validatePath(folderPath)
	if err != nil {
		return nil, err
	}

	u.mu.RLock()
	ws := u.currentWsPath
	u.mu.RUnlock()

	cleanWs := normalizeDriveLetter(filepath.ToSlash(filepath.Clean(ws)))
	isRoot := strings.EqualFold(validatedFolder, cleanWs)

	note, err := u.repo.CreateNote(validatedFolder, trimmedTitle)
	if err != nil {
		return nil, err
	}

	if isRoot {
		note.Folder = "root"
	} else {
		note.Folder = filepath.ToSlash(validatedFolder)
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
