package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/play"
)

func main() {
	game := play.NewGame()
	ebiten.SetWindowSize(config.SCREEN_WIDTH, config.SCREEN_HEIGHT)
	ebiten.SetWindowTitle(config.TITLE)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
