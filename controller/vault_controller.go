package controller

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mdvault/domain/model"
	"mdvault/usecase"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type VaultController struct {
	ctx     context.Context
	useCase usecase.VaultUseCase
}

func NewVaultController(uc usecase.VaultUseCase) *VaultController {
	return &VaultController{
		useCase: uc,
	}
}

func (c *VaultController) Startup(ctx context.Context) {
	c.ctx = ctx
	c.useCase.SetOnChangeCallback(func(path string) {
		runtime.EventsEmit(c.ctx, "on-file-change", path)
	})

	cwd, err := os.Getwd()
	if err == nil {
		defaultWs := filepath.Join(cwd, "markdown")
		if info, err := os.Stat(defaultWs); err == nil && info.IsDir() {
			_, _ = c.useCase.OpenWorkspace(defaultWs)
		}
	}
}

func (c *VaultController) SelectWorkspace() (*model.WorkspaceData, error) {
	selectedDir, err := runtime.OpenDirectoryDialog(c.ctx, runtime.OpenDialogOptions{
		Title: "Chọn thư mục Workspace ghi chú Markdown",
	})
	if err != nil {
		return nil, err
	}
	if selectedDir == "" {
		return nil, fmt.Errorf("không có thư mục nào được chọn")
	}
	return c.useCase.OpenWorkspace(selectedDir)
}

func (c *VaultController) OpenWorkspacePath(dirPath string) (*model.WorkspaceData, error) {
	return c.useCase.OpenWorkspace(dirPath)
}

func (c *VaultController) ScanFolder() (*model.WorkspaceData, error) {
	return c.useCase.GetCurrentWorkspace()
}

func (c *VaultController) GetNoteContent(filePath string) (*model.Note, error) {
	cleanPath := filepath.Clean(strings.TrimSpace(filePath))
	return c.useCase.GetNote(cleanPath)
}

func (c *VaultController) SaveNoteContent(filePath string, content string, title string, tags []string) error {
	cleanPath := filepath.Clean(strings.TrimSpace(filePath))
	return c.useCase.SaveNote(cleanPath, content, title, tags)
}

func (c *VaultController) CreateFolder(parentPath string, folderName string) error {
	cleanParent := filepath.Clean(strings.TrimSpace(parentPath))
	cleanFolder := strings.TrimSpace(folderName)
	return c.useCase.CreateFolder(cleanParent, cleanFolder)
}

func (c *VaultController) CreateNote(folderPath string, noteTitle string) (*model.Note, error) {
	cleanFolder := filepath.Clean(strings.TrimSpace(folderPath))
	cleanTitle := strings.TrimSpace(noteTitle)
	return c.useCase.CreateNote(cleanFolder, cleanTitle)
}

func (c *VaultController) SearchNotes(query string) ([]model.Note, error) {
	return c.useCase.SearchNotes(query)
}

func (c *VaultController) DeleteNote(filePath string) error {
	cleanPath := filepath.Clean(strings.TrimSpace(filePath))
	return c.useCase.DeleteNote(cleanPath)
}

func (c *VaultController) DeleteFolder(folderPath string) error {
	cleanPath := filepath.Clean(strings.TrimSpace(folderPath))
	return c.useCase.DeleteFolder(cleanPath)
}

func (c *VaultController) RenameNote(oldPath string, newTitle string) (*model.Note, error) {
	cleanPath := filepath.Clean(strings.TrimSpace(oldPath))
	cleanTitle := strings.TrimSpace(newTitle)
	return c.useCase.RenameNote(cleanPath, cleanTitle)
}

func (c *VaultController) RenameFolder(oldPath string, newName string) error {
	cleanPath := filepath.Clean(strings.TrimSpace(oldPath))
	cleanName := strings.TrimSpace(newName)
	return c.useCase.RenameFolder(cleanPath, cleanName)
}
