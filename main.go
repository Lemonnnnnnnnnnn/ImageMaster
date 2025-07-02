package main

import (
	"context"
	"embed"
	"log"

	"ImageMaster/core/crawler/api"
	"ImageMaster/core/library"
	"ImageMaster/core/storage"

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
	storageAPI := storage.NewAPI(AppName)

	// 获取配置管理器
	configManager := storageAPI.GetStorage().GetConfigManager()

	// 创建图书馆API
	libraryAPI := library.NewAPI(configManager)

	// 创建爬虫API
	crawlerAPI := api.NewCrawlerAPI(configManager)

	// 设置存储API到爬虫
	crawlerAPI.SetStorageAPI(storageAPI)

	// 创建应用
	err := wails.Run(&options.App{
		Title:  "漫画查看器",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			// 设置context到各个组件
			libraryAPI.SetContext(ctx)
			libraryAPI.InitializeLibraryManager()
			crawlerAPI.SetContext(ctx)
		},
		Bind: []interface{}{
			libraryAPI,
			crawlerAPI,
			storageAPI, // 注册存储API，可以从前端直接调用
		},
		LogLevel:                 logger.ERROR,
		LogLevelProduction:       logger.ERROR,
		EnableDefaultContextMenu: true,
	})

	if err != nil {
		log.Fatal(err)
	}
}
