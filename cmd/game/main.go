package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
)

func main() {
	// Ebiten junk
	ebiten.SetWindowTitle(config.TITLE)
	ebiten.SetWindowSize(config.CANVAS_WIDTH, config.CANVAS_HEIGHT)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)

	// Set locale
	l.Change(config.DEFAULT_LOCALE)

	// Init audio
	audio.InitAudioManager()

	// Do the game
	game := internal.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
