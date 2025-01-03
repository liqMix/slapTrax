package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/locale"
	"github.com/liqmix/ebiten-holiday-2024/internal/resource"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

func main() {
	user.Init()

	ebiten.SetWindowSize(user.S().Graphics.ScreenSizeX, user.S().Graphics.ScreenSizeY)
	ebiten.SetWindowTitle(locale.String(types.L_TITLE))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(ebiten.SyncWithFPS)

	audio.InitAudioManager()
	resource.InitSongs()
	locale.Change(config.DEFAULT_LOCALE)

	// Do the game
	game := internal.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
