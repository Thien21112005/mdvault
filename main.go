package main

import (
	"embed"
	"log"

	"mdvault/controller"
	"mdvault/infrastructure/filesystem"
	"mdvault/usecase"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Khởi tạo các phân lớp theo Clean Architecture
	repo := filesystem.NewVaultRepository()
	uc := usecase.NewVaultUseCase(repo)
	vaultCtrl := controller.NewVaultController(uc)

	err := wails.Run(&options.App{
		Title:  "MDVault - Desktop Markdown Note Taking App",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 255},
		OnStartup:        vaultCtrl.Startup,
		Bind: []interface{}{
			vaultCtrl,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
