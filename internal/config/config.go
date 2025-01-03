package config

// Properties not meant to be customizable
const (
	// System
	TITLE = "Ebiten Holiday Jam 2024"
	// Assets
	SONG_DIR        = "assets/songs"
	SONG_META_NAME  = "meta.yaml"
	SONG_AUDIO_NAME = "audio"
	LOCALE_DIR      = "assets/locales"
	SFX_DIR         = "assets/sfx"
	BGM_DIR         = "assets/bgm"

	// Locale
	DEFAULT_LOCALE             = "en-us"
	FALLBACK_TO_DEFAULT_LOCALE = true

	// UI
	FONT_SCALE = 2.5

	// Audio
	SAMPLE_RATE                 = 48000
	INHERENT_OFFSET             = -200
	SONG_PREVIEW_LENGTH         = 20000
	AUDIO_FADE_MS       float64 = 2000
	AUDIO_FADE_S        float64 = AUDIO_FADE_MS / 1000

	TRAVEL_TIME int64 = 2500 // The amount of time it takes for a note to travel from it's spawn point to the hit zone.
)
