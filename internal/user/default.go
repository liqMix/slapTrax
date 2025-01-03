package user

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

// DefaultSettings defines the initial state for new user settings
var DefaultSettings = UserSettings{
	Graphics: Graphics{
		Fullscreen:   false,
		VSync:        true,
		ScreenSizeX:  1280,
		ScreenSizeY:  720,
		RenderWidth:  1280,
		RenderHeight: 720,
	},
	Audio: Audio{
		BGMVolume:  0.5,
		SFXVolume:  0.5,
		SongVolume: 0.5,
	},
	Gameplay: Gameplay{
		Locale:      config.DEFAULT_LOCALE,
		Theme:       types.ThemeStandard,
		AudioOffset: -30,
		InputOffset: 25,
		// AudioOffset: -235,
		// InputOffset: 35,
		NoteSpeed: 0.25,
	},
	Accessibility: Accessibility{
		NoHoldNotes:  true,
		NoEdgeTracks: false,
	},
}

// NewUserSettings creates a new UserSettings instance initialized with default values
func NewUserSettings() *UserSettings {
	settings := DefaultSettings
	return &settings
}
