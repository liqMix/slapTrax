package user

import "github.com/liqmix/ebiten-holiday-2024/internal/types"

// Properties meant to be customizable
// TODO: store in user profile for persistence
type UserSettings struct {
	// System/Graphics
	Fullscreen bool
	VSync      bool

	ScreenWidth  int
	ScreenHeight int
	RenderWidth  int
	RenderHeight int

	// Game
	Theme       types.Theme
	AudioOffset float64
	InputOffset float64
	NoteSpeed   float64

	// Audio
	BGMVolume         float64
	SFXVolume         float64
	SongVolume        float64
	SongPreviewVolume float64
}

var DefaultSettings = UserSettings{
	Theme:             types.ThemeDefault,
	BGMVolume:         0.5,
	SFXVolume:         0.5,
	SongVolume:        0.5,
	SongPreviewVolume: 0.5,
	RenderWidth:       1280,
	RenderHeight:      720,
	NoteSpeed:         1.0,
}
