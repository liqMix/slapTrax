package config

// Properties not meant to be customizable
const (
	// DEBUG
	DEBUG = true

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
	SFX_VOLUME                  = 0.5
	SONG_PREVIEW_LENGTH         = 10000
	SONG_PREVIEW_FADE   float64 = 1000.0

	// Game
	AUDIO_OFFSET int64 = -235 // The amount of time the audio is ahead of the notes.
	INPUT_OFFSET int64 = -5   // The amount of time the input is ahead of the notes.
	// GRACE_PERIOD int64   = 2000 // The amount of time before the song starts that the player can get ready.
	TRAVEL_TIME int64   = 2500 // The amount of time it takes for a note to travel from it's spawn point to the hit zone.
	NOTE_SPEED  float64 = 0.25 // The speed at which notes travel. total travel time = TRAVEL_TIME / NOTE_SPEED

	// Accessibility
	NO_EDGE_TRACKS bool = true // If true, the edge tracks will not be used.
)

var (
	GracePeriod                     = ActualTravelTimeInt64 / 2
	ActualTravelTimeInt64   int64   = int64(float64(TRAVEL_TIME) / NOTE_SPEED)
	ActualTravelTimeFloat64 float64 = float64(TRAVEL_TIME) / NOTE_SPEED
)
