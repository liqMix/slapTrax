package config

// Properties not meant to be customizable
const (
	// System
	TITLE = "Ebiten Holiday Jam 2024"

	// UI
	FONT_SCALE = 2.5

	// Audio
	GRACE_PERIOD                = 4000
	SAMPLE_RATE                 = 48000
	INHERENT_OFFSET             = -200
	SONG_PREVIEW_LENGTH         = 20000
	AUDIO_FADE_MS       float64 = 2000
	AUDIO_FADE_S        float64 = AUDIO_FADE_MS / 1000
	TRAVEL_TIME         int64   = 2500 // The amount of time it takes for a note to travel from it's spawn point to the hit zone.

	// Server
	SERVER_ENDPOINT = "http://localhost:8080"
)
