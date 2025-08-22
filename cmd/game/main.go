package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal"
	"github.com/liqmix/slaptrax/internal/assets"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/cache"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/l"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/user"
)

func main() {
	err := user.Init()
	if err != nil {
		logger.Warn("Failed to initialize user: %v", err)
	}
	display.InitWindow()
	cache.InitCaches()

	assets.Init(
		assets.AssetInit{
			Locale: user.S().Locale,
		})
	audio.InitAudioManager(&audio.Volume{
		BGM:  user.S().BGMVolume,
		SFX:  user.S().SFXVolume,
		Song: user.S().SongVolume,
	})

	input.InitInput()
	defer input.Close()

	// Ebiten setup
	ebiten.SetWindowSize(user.S().ScreenWidth, user.S().ScreenHeight)
	ebiten.SetFullscreen(user.S().Fullscreen)
	ebiten.SetWindowTitle(l.String(l.TITLE))
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)

	// Do the game thing
	game := internal.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		logger.Error("Game error: %v", err)
		return
	}
}
