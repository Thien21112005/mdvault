package controller

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	return c.useCase.GetNote(filePath)
}

func (c *VaultController) SaveNoteContent(filePath string, content string, title string, tags []string) error {
	return c.useCase.SaveNote(filePath, content, title, tags)
}

func (c *VaultController) CreateFolder(parentPath string, folderName string) error {
	return c.useCase.CreateFolder(parentPath, folderName)
}

func (c *VaultController) CreateNote(folderPath string, noteTitle string) (*model.Note, error) {
	return c.useCase.CreateNote(folderPath, noteTitle)
}

func (c *VaultController) SearchNotes(query string) ([]model.Note, error) {
	return c.useCase.SearchNotes(query)
}
