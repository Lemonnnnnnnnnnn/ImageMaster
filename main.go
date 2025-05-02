package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 创建应用实例
	viewer := NewViewer()

	// 创建应用
	err := wails.Run(&options.App{
		Title:  "漫画查看器",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        viewer.Startup,
		Bind: []interface{}{
			viewer,
		},
		LogLevel:           logger.ERROR,
		LogLevelProduction: logger.ERROR,
	})

	if err != nil {
		log.Fatal(err)
	}
}
