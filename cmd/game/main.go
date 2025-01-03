package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

func main() {
	user.Init()
	assets.InitLocales()
	assets.SetLocale(user.S().Gameplay.Locale)

	ebiten.SetWindowSize(user.S().Graphics.ScreenSizeX, user.S().Graphics.ScreenSizeY)
	ebiten.SetWindowTitle(assets.String(types.L_TITLE))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(ebiten.SyncWithFPS)

	assets.InitAudioManager(user.Volume())
	assets.InitSongs()
	ui.InitTextRenderer()

	// Do the game
	game := internal.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
