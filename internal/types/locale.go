package types

const (
	L_TITLE = "title"

	L_DIFFICULTY_EASY    = "difficulty.easy"
	L_DIFFICULTY_MEDIUM  = "difficulty.medium"
	L_DIFFICULTY_HARD    = "difficulty.hard"
	L_DIFFICULTY_UNKNOWN = "difficulty.unknown"
	L_HIT_PERFECT        = "hit.perfect"
	L_HIT_GOOD           = "hit.good"
	L_HIT_BAD            = "hit.bad"
	L_HIT_MISS           = "hit.miss"

	// States
	L_STATE_TITLE    = "state.greeting"
	L_STATE_MENU     = "state.menu"
	L_STATE_EDITOR   = "state.editor"
	L_STATE_SETTINGS = "state.settings"
	L_STATE_PROFILE  = "state.profile"
	L_STATE_PLAY     = "state.play"

	// Settings
	//// System/Graphics
	L_SETTINGS_GFX_FULLSCREEN = "settings.gfx.fullscreen"
	L_SETTINGS_GFX_VSYNC      = "settings.gfx.vsync"
	L_SETTINGS_GFX_WIDTH      = "settings.gfx.width"
	L_SETTINGS_GFX_HEIGHT     = "settings.gfx.height"
	L_SETTINGS_GFX_RENDERW    = "settings.gfx.renderw"
	L_SETTINGS_GFX_RENDERH    = "settings.gfx.renderh"

	//// Game
	L_SETTINGS_GAME_THEME       = "settings.game.theme"
	L_SETTINGS_GAME_AUDIOOFFSET = "settings.game.audiooffset"
	L_SETTINGS_GAME_INPUTOFFSET = "settings.game.inputoffset"
	L_SETTINGS_GAME_NOTESPEED   = "settings.game.notespeed"

	//// Audio
	L_SETTINGS_AUDIO_BGMVOLUME         = "settings.audio.bgmvolume"
	L_SETTINGS_AUDIO_SFXVOLUME         = "settings.audio.sfxvolume"
	L_SETTINGS_AUDIO_SONGVOLUME        = "settings.audio.songvolume"
	L_SETTINGS_AUDIO_SONGPREVIEWVOLUME = "settings.audio.songpreviewvolume"

	// States
	//// Play
	L_STATE_PLAY_RESTART = "state.play.restart"

	// Themes
	L_THEME_STANDARD   = "theme.standard"
	L_THEME_LEFTBEHIND = "theme.leftbehind"

	// Etc
	L_UNKNOWN = "unknown"
)

var AllLocaleKeys = []string{
	L_TITLE,
	L_DIFFICULTY_EASY,
	L_DIFFICULTY_MEDIUM,
	L_DIFFICULTY_HARD,
	L_DIFFICULTY_UNKNOWN,
	L_STATE_TITLE,
	L_STATE_MENU,
	L_STATE_EDITOR,
	L_STATE_SETTINGS,
	L_STATE_PROFILE,
	L_STATE_PLAY,
	L_STATE_PLAY_RESTART,
	L_SETTINGS_GFX_FULLSCREEN,
	L_SETTINGS_GFX_VSYNC,
	L_SETTINGS_GFX_WIDTH,
	L_SETTINGS_GFX_HEIGHT,
	L_SETTINGS_GFX_RENDERW,
	L_SETTINGS_GFX_RENDERH,
	L_SETTINGS_GAME_THEME,
	L_SETTINGS_GAME_AUDIOOFFSET,
	L_SETTINGS_GAME_INPUTOFFSET,
	L_SETTINGS_GAME_NOTESPEED,
	L_SETTINGS_AUDIO_BGMVOLUME,
	L_SETTINGS_AUDIO_SFXVOLUME,
	L_SETTINGS_AUDIO_SONGVOLUME,
	L_SETTINGS_AUDIO_SONGPREVIEWVOLUME,
	L_HIT_PERFECT,
	L_HIT_GOOD,
	L_HIT_BAD,
	L_HIT_MISS,
	L_THEME_STANDARD,
	L_THEME_LEFTBEHIND,
	L_UNKNOWN,
}
