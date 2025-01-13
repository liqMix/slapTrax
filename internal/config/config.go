package config

// Properties not meant to be customizable
const (
	// System
	TITLE = "slapTrax"

	// UI
	FONT_SCALE = 2.5

	// Audio
	SAMPLE_RATE                 = 48000
	SONG_PREVIEW_LENGTH         = 20000
	AUDIO_FADE_MS       float64 = 2000
	AUDIO_FADE_S        float64 = AUDIO_FADE_MS / 1000
	DISABLE_LOGGER              = false

	// Server
	SERVER_ENDPOINT = "https://liq.mx/slapi/v1"
)
