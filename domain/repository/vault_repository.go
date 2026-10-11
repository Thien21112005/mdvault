package repository

import (
	"mdvault/domain/model"
)

// VaultRepository định nghĩa giao thức quản lý dữ liệu tệp tin trong Clean Architecture
type VaultRepository interface {
	ScanWorkspace(rootPath string) (*model.WorkspaceData, error)
	ParseNote(filePath string, folderID string) (*model.Note, error)
	SaveNote(filePath string, content string, title string, tags []string) error
	CreateFolder(parentPath string, folderName string) error
	CreateNote(folderPath string, noteTitle string) (*model.Note, error)
	DeleteNote(filePath string) error
	DeleteFolder(folderPath string) error
	RenameNote(oldPath string, newPath string) error
	RenameFolder(oldPath string, newPath string) error
	StartWatcher(rootPath string, onChange func(path string)) error
	StopWatcher() error
	UpdateWatchPath(dirPath string) error
}
