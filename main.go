package main

import (
	"embed"
	"log"

	"ImageMaster/core/downloader"
	"ImageMaster/core/storage"
	"ImageMaster/core/viewer"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

const AppName = "imagemaster"

func main() {
	// 创建存储API
	storageAPI := storage.NewStorageAPI(AppName)

	// 获取配置管理器
	configManager := storageAPI.GetStorage().GetConfigManager()

	// 创建新的查看器实例，传入配置管理器
	v := viewer.NewViewer(configManager)

	// 创建下载管理器API
	downloaderAPI := downloader.NewDownloaderAPI(configManager)

	// 设置存储API到下载器
	downloaderAPI.SetStorageAPI(storageAPI)

	// 创建应用
	err := wails.Run(&options.App{
		Title:  "漫画查看器",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        v.Startup,
		Bind: []interface{}{
			v,
			downloaderAPI,
			storageAPI, // 注册存储API，可以从前端直接调用
		},
		LogLevel:           logger.ERROR,
		LogLevelProduction: logger.ERROR,
	})

	if err != nil {
		log.Fatal(err)
	}

	// 启动查看器
	// v.Startup(context.Background())
}
