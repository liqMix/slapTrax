package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/play"
)

func main() {
	// Ebiten junk
	ebiten.SetWindowTitle(config.TITLE)
	ebiten.SetWindowSize(config.CANVAS_WIDTH, config.CANVAS_HEIGHT)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Set locale
	l.Change(config.DEFAULT_LOCALE)

	// Do the game
	game := play.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
