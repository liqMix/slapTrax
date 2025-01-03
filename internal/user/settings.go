package user

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/locale"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

// Graphics contains all display and rendering related settings
type Graphics struct {
	Fullscreen   bool
	VSync        bool
	ScreenSizeX  int
	ScreenSizeY  int
	RenderWidth  int
	RenderHeight int
}

// func (g *Graphics) SetScreenSize(width, height int) {
// 	g.ScreenSizeX = width
// 	g.ScreenSizeY = height
// }

// func (g *Graphics) ScreenSize() (int, int) {
// 	return g.ScreenSizeX, g.ScreenSizeY
// }

func (g *Graphics) Apply() {
	// ebiten.SetFullscreen(g.Fullscreen)
	// ebiten.SetVsyncEnabled(g.VSync)
	// ebiten.SetWindowSize(g.ScreenSizeX, g.ScreenSizeY)
	// cache.ClearImageCache()
}

// Audio contains all volume related settings
type Audio struct {
	BGMVolume         float64 // 0.0-1.0
	SFXVolume         float64 // 0.0-1.0
	SongVolume        float64 // 0.0-1.0
	SongPreviewVolume float64 // 0.0-1.0
}

func (a *Audio) Apply() {}

// Gameplay contains core game mechanics settings
type Gameplay struct {
	Locale      string
	Theme       types.Theme
	AudioOffset int64   // Milliseconds ahead of notes (negative = earlier)
	InputOffset int64   // Milliseconds ahead of notes (negative = earlier)
	NoteSpeed   float64 // Travel speed multiplier (0.1-2.0)
}

func (g *Gameplay) Apply() {
	locale.Change(g.Locale)
}

// Accessibility contains settings for game accessibility features
type Accessibility struct {
	NoHoldNotes  bool
	NoEdgeTracks bool
}

// UserSettings contains all customizable properties for a user's game experience
type UserSettings struct {
	Graphics      Graphics
	Audio         Audio
	Gameplay      Gameplay
	Accessibility Accessibility
}

// Setting keys for localization and UI mapping
const (
	SettingsGraphicsFullscreen = types.L_SETTINGS_GFX_FULLSCREEN
	SettingsGraphicsVSync      = types.L_SETTINGS_GFX_VSYNC
	SettingsGraphicsScreen     = types.L_SETTINGS_GFX_SCREENSIZE
	SettingsGraphicsRender     = types.L_SETTINGS_GFX_RENDERSIZE

	SettingsGameLocale      = types.L_SETTINGS_GAME_LOCALE
	SettingsGameTheme       = types.L_SETTINGS_GAME_THEME
	SettingsGameAudioOffset = types.L_SETTINGS_GAME_AUDIOOFFSET
	SettingsGameInputOffset = types.L_SETTINGS_GAME_INPUTOFFSET
	SettingsGameNoteSpeed   = types.L_SETTINGS_GAME_NOTESPEED

	SettingsAudioBGM     = types.L_SETTINGS_AUDIO_BGMVOLUME
	SettingsAudioSFX     = types.L_SETTINGS_AUDIO_SFXVOLUME
	SettingsAudioSong    = types.L_SETTINGS_AUDIO_SONGVOLUME
	SettingsAudioPreview = types.L_SETTINGS_AUDIO_SONGPREVIEWVOLUME

	SettingsAccessHoldNotes  = types.L_SETTINGS_ACCESS_NOHOLDNOTES
	SettingsAccessEdgeTracks = types.L_SETTINGS_ACCESS_NOEDGETRACKS
)

// Updates application with user settings
func (s *UserSettings) Apply() {
	s.Graphics.Apply()
	s.Audio.Apply()
	s.Gameplay.Apply()
}

// Copy creates a deep copy of UserSettings
func (s *UserSettings) Copy() *UserSettings {
	copy := *s
	return &copy
}
