package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/locale"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

func main() {
	// TODO: Read in user settings file if it exists to pick up existing settings?
	s := user.Settings()

	// Ebiten junk
	ebiten.SetWindowSize(s.RenderWidth, s.RenderHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)

	// Set locale
	locale.Change(config.DEFAULT_LOCALE)
	ebiten.SetWindowTitle(locale.String(types.L_TITLE))

	// Init audio
	audio.InitAudioManager()
	song.InitSongs()

	// Do the game
	game := internal.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
