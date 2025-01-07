package user

import (
	"encoding/json"
	"fmt"
	"os"
)

const settingsFilename = "settings.json"

// Settings contains all customizable properties
type Settings struct {
	Locale              string  `json:"locale"`
	Version             string  `json:"version"`
	Fullscreen          bool    `json:"fullscreen"`
	ScreenWidth         int     `json:"screenWidth"`
	ScreenHeight        int     `json:"screenHeight"`
	RenderWidth         int     `json:"renderWidth"`
	RenderHeight        int     `json:"renderHeight"`
	FixedRenderScale    bool    `json:"fixedRenderScale"`
	BGMVolume           float64 `json:"bgmVolume"`
	SFXVolume           float64 `json:"sfxVolume"`
	SongVolume          float64 `json:"songVolume"`
	LaneSpeed           float64 `json:"laneSpeed"`
	AudioOffset         int64   `json:"audioOffset"`
	InputOffset         int64   `json:"inputOffset"`
	WaveringLane        bool    `json:"waveringLane"`
	NoteColorTheme      string  `json:"noteColorTheme"`
	CenterNoteColor     string  `json:"centerNoteColor"`
	CornerNoteColor     string  `json:"cornerNoteColor"`
	DisableHoldNotes    bool    `json:"disableHoldNotes"`
	DisableHitEffects   bool    `json:"disableHitEffects"`
	DisableLaneEffects  bool    `json:"disableLaneEffects"`
	PromptedOffsetCheck bool    `json:"promptedOffsetCheck"`
}

// Default values
func NewSettings() *Settings {
	return &Settings{

		Locale:              "en-us",
		Version:             "0.0.1",
		Fullscreen:          false,
		ScreenWidth:         1280,
		ScreenHeight:        720,
		RenderWidth:         1280,
		RenderHeight:        720,
		FixedRenderScale:    false,
		BGMVolume:           0.7,
		SFXVolume:           0.7,
		SongVolume:          0.7,
		LaneSpeed:           1.0,
		AudioOffset:         0,
		InputOffset:         0,
		WaveringLane:        false,
		NoteColorTheme:      "note.color.default",
		CenterNoteColor:     "#e6e600ff",
		CornerNoteColor:     "#e68200ff",
		DisableHoldNotes:    false,
		DisableHitEffects:   false,
		DisableLaneEffects:  false,
		PromptedOffsetCheck: false,
	}
}

func (s *Settings) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	return os.WriteFile(settingsFilename, data, 0644)
}

func LoadSettings() (*Settings, error) {
	data, err := os.ReadFile(settingsFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	settings := NewSettings()
	if err := json.Unmarshal(data, settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return settings, nil
}

func (s *Settings) MergeFrom(other *Settings) {
	if other == nil {
		return
	}

	if other.Locale != "" {
		s.Locale = other.Locale
	}
	if other.Version != "" {
		s.Version = other.Version
	}
	if other.Fullscreen {
		s.Fullscreen = other.Fullscreen
	}
	if other.ScreenWidth > 0 {
		s.ScreenWidth = other.ScreenWidth
	}
	if other.ScreenHeight > 0 {
		s.ScreenHeight = other.ScreenHeight
	}
	if other.RenderWidth > 0 {
		s.RenderWidth = other.RenderWidth
	}
	if other.RenderHeight > 0 {
		s.RenderHeight = other.RenderHeight
	}
	if other.FixedRenderScale {
		s.FixedRenderScale = other.FixedRenderScale
	}
	if other.BGMVolume > 0 {
		s.BGMVolume = other.BGMVolume
	}
	if other.SFXVolume > 0 {
		s.SFXVolume = other.SFXVolume
	}
	if other.SongVolume > 0 {
		s.SongVolume = other.SongVolume
	}
	if other.LaneSpeed > 0 {
		s.LaneSpeed = other.LaneSpeed
	}
	if other.AudioOffset != 0 {
		s.AudioOffset = other.AudioOffset
	}
	if other.InputOffset != 0 {
		s.InputOffset = other.InputOffset
	}
	if other.WaveringLane {
		s.WaveringLane = other.WaveringLane
	}
	if other.NoteColorTheme != "" {
		s.NoteColorTheme = other.NoteColorTheme
	}
	if other.CenterNoteColor != "" {
		s.CenterNoteColor = other.CenterNoteColor
	}
	if other.CornerNoteColor != "" {
		s.CornerNoteColor = other.CornerNoteColor
	}
	if other.DisableHoldNotes {
		s.DisableHoldNotes = other.DisableHoldNotes
	}
	if other.DisableHitEffects {
		s.DisableHitEffects = other.DisableHitEffects
	}
	if other.DisableLaneEffects {
		s.DisableLaneEffects = other.DisableLaneEffects
	}
}
