package types

const (
	L_TITLE              = "title"
	L_EXIT               = "exit"
	L_BACK               = "back"
	L_SAVE               = "save"
	L_CANCEL             = "cancel"
	L_OFF                = "off"
	L_ON                 = "on"
	L_CHART              = "chart"
	L_DIFFICULTIES       = "difficulties"
	L_DIFFICULTY_EASY    = "difficulty.easy"
	L_DIFFICULTY_MEDIUM  = "difficulty.medium"
	L_DIFFICULTY_HARD    = "difficulty.hard"
	L_DIFFICULTY_UNKNOWN = "difficulty.unknown"
	L_HIT_PERFECT        = "hit.perfect"
	L_HIT_GOOD           = "hit.good"
	L_HIT_BAD            = "hit.bad"
	L_HIT_MISS           = "hit.miss"

	// Song
	L_SONG_ARTIST = "song.artist"
	L_SONG_ALBUM  = "song.album"

	// States
	L_STATE_TITLE                = "state.greeting"
	L_STATE_EDITOR               = "state.editor"
	L_STATE_SETTINGS             = "state.settings"
	L_STATE_PROFILE              = "state.profile"
	L_STATE_PLAY                 = "state.play"
	L_STATE_OFFSET               = "state.offset"
	L_STATE_SONG_SELECTION       = "state.song.selection"
	L_STATE_DIFFICULTY_SELECTION = "state.difficulty.selection"

	// Settings
	//// System/Graphics
	L_SETTINGS_GFX            = "settings.gfx"
	L_SETTINGS_GFX_FULLSCREEN = "settings.gfx.fullscreen"
	L_SETTINGS_GFX_VSYNC      = "settings.gfx.vsync"
	L_SETTINGS_GFX_SCREENSIZE = "settings.gfx.screensize"
	L_SETTINGS_GFX_RENDERSIZE = "settings.gfx.rendersize"
	L_RENDERSIZE_TINY         = "rendersize.tiny"
	L_RENDERSIZE_SMALL        = "rendersize.small"
	L_RENDERSIZE_MEDIUM       = "rendersize.medium"
	L_RENDERSIZE_LARGE        = "rendersize.large"
	L_RENDERSIZE_MAX          = "rendersize.max"

	//// Game
	L_SETTINGS_GAME             = "settings.game"
	L_SETTINGS_GAME_LOCALE      = "settings.game.locale"
	L_SETTINGS_GAME_THEME       = "settings.game.theme"
	L_SETTINGS_GAME_AUDIOOFFSET = "settings.game.audiooffset"
	L_SETTINGS_GAME_INPUTOFFSET = "settings.game.inputoffset"
	L_SETTINGS_GAME_NOTESPEED   = "settings.game.notespeed"

	//// Audio
	L_SETTINGS_AUDIO                   = "settings.audio"
	L_SETTINGS_AUDIO_BGMVOLUME         = "settings.audio.bgmvolume"
	L_SETTINGS_AUDIO_SFXVOLUME         = "settings.audio.sfxvolume"
	L_SETTINGS_AUDIO_SONGVOLUME        = "settings.audio.songvolume"
	L_SETTINGS_AUDIO_SONGPREVIEWVOLUME = "settings.audio.songpreviewvolume"

	// Accessibility
	L_SETTINGS_ACCESS              = "settings.access"
	L_SETTINGS_ACCESS_NOHOLDNOTES  = "settings.access.noholdnotes"
	L_SETTINGS_ACCESS_NOEDGETRACKS = "settings.access.noedgetracks"

	// States
	//// Play
	L_STATE_PLAY_RESTART = "state.play.restart"
	L_STATE_PLAY_PAUSE   = "state.play.pause"

	// Themes
	L_THEME_STANDARD   = "theme.standard"
	L_THEME_LEFTBEHIND = "theme.leftbehind"

	// Etc
	L_UNKNOWN = "unknown"
)

var AllLocaleKeys = []string{}
