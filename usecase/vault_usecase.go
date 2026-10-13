package usecase

import (
	"fmt"
	"os"
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
	DeleteNote(filePath string) error
	DeleteFolder(folderPath string) error
	RenameNote(oldPath string, newTitle string) (*model.Note, error)
	RenameFolder(oldPath string, newName string) error
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

func matchKeyword(s, qOriginal, qLower string) bool {
	if qLower == "" {
		return true
	}
	if strings.Contains(s, qOriginal) || strings.Contains(s, qLower) {
		return true
	}
	return strings.Contains(strings.ToLower(s), qLower)
}

func (u *vaultUseCaseImpl) SearchNotes(query string) ([]model.Note, error) {
	data, err := u.GetCurrentWorkspace()
	if err != nil || data == nil {
		return nil, err
	}

	qOriginal := strings.TrimSpace(query)
	if qOriginal == "" {
		return data.Notes, nil
	}
	qLower := strings.ToLower(qOriginal)

	var results []model.Note
	for _, n := range data.Notes {
		// 1. Kiểm tra Tiêu đề trước (chuỗi ngắn, tốc độ nhanh nhất)
		if matchKeyword(n.Title, qOriginal, qLower) {
			results = append(results, n)
			continue
		}

		// 2. Kiểm tra Thẻ Tags (chuỗi ngắn)
		tagMatches := false
		for _, t := range n.Tags {
			if matchKeyword(t, qOriginal, qLower) {
				tagMatches = true
				break
			}
		}
		if tagMatches {
			results = append(results, n)
			continue
		}

		// 3. Chỉ kiểm tra Nội dung toàn văn khi Tiêu đề và Tags không khớp
		if matchKeyword(n.Content, qOriginal, qLower) {
			results = append(results, n)
		}
	}
	return results, nil
}

func (u *vaultUseCaseImpl) DeleteNote(filePath string) error {
	validatedPath, err := u.validatePath(filePath)
	if err != nil {
		return err
	}

	u.mu.RLock()
	ws := u.currentWsPath
	u.mu.RUnlock()

	cleanWs := normalizeDriveLetter(filepath.ToSlash(filepath.Clean(ws)))
	if strings.EqualFold(validatedPath, cleanWs) {
		return fmt.Errorf("truy cập bị từ chối: không thể xóa thư mục gốc workspace")
	}

	info, err := os.Stat(validatedPath)
	if err != nil {
		return fmt.Errorf("ghi chú không tồn tại: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("đường dẫn là thư mục, không phải ghi chú")
	}

	if err := u.repo.DeleteNote(validatedPath); err != nil {
		return fmt.Errorf("lỗi khi xóa ghi chú: %w", err)
	}

	_, _ = u.RefreshCache()
	return nil
}

func (u *vaultUseCaseImpl) DeleteFolder(folderPath string) error {
	validatedPath, err := u.validatePath(folderPath)
	if err != nil {
		return err
	}

	u.mu.RLock()
	ws := u.currentWsPath
	u.mu.RUnlock()

	cleanWs := normalizeDriveLetter(filepath.ToSlash(filepath.Clean(ws)))
	if strings.EqualFold(validatedPath, cleanWs) {
		return fmt.Errorf("truy cập bị từ chối: không thể xóa thư mục gốc workspace")
	}

	info, err := os.Stat(validatedPath)
	if err != nil {
		return fmt.Errorf("thư mục không tồn tại: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("đường dẫn là tệp tin, không phải thư mục")
	}

	if err := u.repo.DeleteFolder(validatedPath); err != nil {
		return fmt.Errorf("lỗi khi xóa thư mục: %w", err)
	}

	_, _ = u.RefreshCache()
	return nil
}

func (u *vaultUseCaseImpl) RenameNote(oldPath string, newTitle string) (*model.Note, error) {
	cleanTitle := strings.TrimSpace(newTitle)
	if cleanTitle == "" {
		return nil, fmt.Errorf("tiêu đề mới không được để trống")
	}
	if strings.Contains(cleanTitle, "..") || strings.Contains(cleanTitle, "/") || strings.Contains(cleanTitle, "\\") {
		return nil, fmt.Errorf("truy cập bị từ chối: tiêu đề không được chứa ký tự phân cách hoặc '..'")
	}

	validatedOld, err := u.validatePath(oldPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(validatedOld)
	if err != nil {
		return nil, fmt.Errorf("ghi chú không tồn tại: %w", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("đường dẫn là thư mục, không phải ghi chú")
	}

	oldNote, err := u.GetNote(validatedOld)
	if err != nil {
		return nil, fmt.Errorf("không thể đọc thông tin ghi chú: %w", err)
	}

	dir := filepath.Dir(validatedOld)
	newFileName := strings.ToLower(strings.ReplaceAll(cleanTitle, " ", "-"))
	newFileName = filepath.Base(newFileName)
	if !strings.HasSuffix(newFileName, ".md") {
		newFileName += ".md"
	}
	newPath := filepath.ToSlash(filepath.Join(dir, newFileName))

	if !strings.EqualFold(validatedOld, newPath) {
		if _, err := os.Stat(newPath); err == nil {
			return nil, fmt.Errorf("ghi chú với tên '%s' đã tồn tại", newFileName)
		}
		if err := u.repo.RenameNote(validatedOld, newPath); err != nil {
			return nil, fmt.Errorf("lỗi khi đổi tên tệp ghi chú: %w", err)
		}
	}

	if err := u.repo.SaveNote(newPath, oldNote.Content, cleanTitle, oldNote.Tags); err != nil {
		return nil, fmt.Errorf("lỗi khi cập nhật tiêu đề mới: %w", err)
	}

	_, _ = u.RefreshCache()
	return u.GetNote(newPath)
}

func (u *vaultUseCaseImpl) RenameFolder(oldPath string, newName string) error {
	cleanName := strings.TrimSpace(newName)
	if cleanName == "" {
		return fmt.Errorf("tên thư mục mới không được để trống")
	}
	if strings.Contains(cleanName, "..") || strings.Contains(cleanName, "/") || strings.Contains(cleanName, "\\") {
		return fmt.Errorf("truy cập bị từ chối: tên thư mục không được chứa ký tự phân cách hoặc '..'")
	}

	validatedOld, err := u.validatePath(oldPath)
	if err != nil {
		return err
	}

	u.mu.RLock()
	ws := u.currentWsPath
	u.mu.RUnlock()

	cleanWs := normalizeDriveLetter(filepath.ToSlash(filepath.Clean(ws)))
	if strings.EqualFold(validatedOld, cleanWs) {
		return fmt.Errorf("truy cập bị từ chối: không thể đổi tên thư mục gốc workspace")
	}

	info, err := os.Stat(validatedOld)
	if err != nil {
		return fmt.Errorf("thư mục không tồn tại: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("đường dẫn là tệp tin, không phải thư mục")
	}

	parent := filepath.Dir(validatedOld)
	cleanName = filepath.Base(cleanName)
	newFolderPath := filepath.ToSlash(filepath.Join(parent, cleanName))

	if strings.EqualFold(validatedOld, newFolderPath) {
		return nil
	}

	if _, err := os.Stat(newFolderPath); err == nil {
		return fmt.Errorf("thư mục '%s' đã tồn tại", cleanName)
	}

	if err := u.repo.RenameFolder(validatedOld, newFolderPath); err != nil {
		return fmt.Errorf("lỗi khi đổi tên thư mục: %w", err)
	}

	_, _ = u.RefreshCache()
	return nil
}
