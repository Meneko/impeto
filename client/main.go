package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "impeto",
		Description: "Impeto",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Services: []application.Service{
			application.NewService(NewAppsService()),
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Impeto",
		Width:            351,
		Height:           200,
		BackgroundColour: application.NewRGBA(6, 7, 15, 0),
		BackgroundType:   application.BackgroundTypeTransparent,
		URL:              "/",
		AlwaysOnTop:      true,
	})

	window.Center()

	x, _ := window.Position()
	window.SetPosition(x, 30)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
