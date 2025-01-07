package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

func main() {
	err := user.Init()
	if err != nil {
		logger.Warn("Failed to initialize user: %v", err)
	}
	display.InitWindow()

	renderWidth, renderHeight := display.Window.RenderSize()
	cache.InitCaches(renderWidth, renderHeight)

	// Refresh window after cache init
	display.Window.Refresh()

	assets.Init(
		assets.AssetInit{
			Locale: user.S.Locale,
		})
	audio.InitAudioManager(&audio.Volume{
		BGM:  user.S.BGMVolume,
		SFX:  user.S.SFXVolume,
		Song: user.S.SongVolume,
	})

	ebiten.SetWindowSize(user.S.ScreenWidth, user.S.ScreenHeight)
	ebiten.SetFullscreen(user.S.Fullscreen)
	ebiten.SetWindowTitle(l.String(l.TITLE))
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)

	// Do the game thing
	game := internal.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
