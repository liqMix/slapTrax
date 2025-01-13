package l

import "github.com/liqmix/ebiten-holiday-2024/internal/assets"

const (
	LOCALE = "locale" // The current locale name in locale language

	TITLE        = "title"
	EXIT         = "exit"
	BACK         = "back"
	SAVE         = "save"
	OK           = "ok"
	CANCEL       = "cancel"
	CONTINUE     = "continue"
	OFF          = "off"
	ON           = "on"
	GUEST        = "guest"
	LOADING      = "loading"
	CHART        = "chart"
	CENTER       = "center"
	CORNER       = "corner"
	DIFFICULTIES = "difficulties"
	WELCOME      = "welcome"
	LEADERBOARD  = "leaderboard"

	// Actions
	ACTION_BACK   = "action.back"
	ACTION_UP     = "action.up"
	ACTION_DOWN   = "action.down"
	ACTION_LEFT   = "action.left"
	ACTION_RIGHT  = "action.right"
	ACTION_SELECT = "action.select"

	// Song
	SONG_ARTIST = "song.artist"
	SONG_ALBUM  = "song.album"

	// States
	STATE_TITLE                = "state.greeting"
	STATE_EDITOR               = "state.editor"
	STATE_SETTINGS             = "state.settings"
	STATE_PROFILE              = "state.profile"
	STATE_PLAY                 = "state.play"
	STATE_RESULT               = "state.result"
	STATE_OFFSET               = "state.offset"
	STATE_SONG_SELECTION       = "state.song.selection"
	STATE_LOGIN                = "state.login"
	STATE_DIFFICULTY_SELECTION = "state.difficulty.selection"
	STATE_HOW_TO_PLAY          = "state.howtoplay"

	// Login
	LOGIN_TEXT_OFFLINE = "login.text.offline"
	LOGIN_TEXT_ONLINE  = "login.text.online"
	LOGIN_SAVE_LOCAL   = "login.save.local"
	LOGIN_CONTINUE     = "login.continue"
	LOGIN_LOGIN        = "login.login"
	LOGIN_LOGOUT       = "login.logout"
	LOGIN_USERNAME     = "login.username"
	LOGIN_PASSWORD     = "login.password"

	// Settings
	//// System/Graphics
	SETTINGS_GFX             = "settings.gfx"
	SETTINGS_GFX_FULLSCREEN  = "settings.gfx.fullscreen"
	SETTINGS_GFX_VSYNC       = "settings.gfx.vsync"
	SETTINGS_GFX_SCREENSIZE  = "settings.gfx.screensize"
	SETTINGS_GFX_RENDERSIZE  = "settings.gfx.rendersize"
	SETTINGS_GFX_FIXEDRENDER = "settings.gfx.fixedrender"
	SETTINGS_GFX_NOTECOLOR   = "settings.gfx.notecolor"

	RENDERSIZE_TINY   = "rendersize.tiny"
	RENDERSIZE_SMALL  = "rendersize.small"
	RENDERSIZE_MEDIUM = "rendersize.medium"
	RENDERSIZE_LARGE  = "rendersize.large"
	RENDERSIZE_MAX    = "rendersize.max"

	//// Game
	SETTINGS_GAME              = "settings.game"
	SETTINGS_GAME_KEY_CONFIG   = "settings.game.keyconfig"
	SETTINGS_GAME_NOTEWIDTH    = "settings.game.notewidth"
	SETTINGS_GAME_LOCALE       = "settings.game.locale"
	SETTINGS_GAME_THEME        = "settings.game.theme"
	SETTINGS_GAME_AUDIOOFFSET  = "settings.game.audiooffset"
	SETTINGS_GAME_INPUTOFFSET  = "settings.game.inputoffset"
	SETTINGS_GAME_LANESPEED    = "settings.game.lanespeed"
	SETTINGS_GAME_EDGEPLAYAREA = "settings.game.edgeplayarea"

	//// Audio
	SETTINGS_AUDIO                   = "settings.audio"
	SETTINGS_AUDIO_BGMVOLUME         = "settings.audio.bgmvolume"
	SETTINGS_AUDIO_SFXVOLUME         = "settings.audio.sfxvolume"
	SETTINGS_AUDIO_SONGVOLUME        = "settings.audio.songvolume"
	SETTINGS_AUDIO_SONGPREVIEWVOLUME = "settings.audio.songpreviewvolume"

	// Accessibility
	SETTINGS_ACCESS              = "settings.access"
	SETTINGS_ACCESS_NOHOLDNOTES  = "settings.access.noholdnotes"
	SETTINGS_ACCESS_NOHITEFFECT  = "settings.access.nohiteffect"
	SETTINGS_ACCESS_NOLANEEFFECT = "settings.access.nolaneeffect"
	SETTINGS_ACCESS_MIRROR       = "settings.access.mirror"

	KEY_CONFIG_DEFAULT      = "keyconfig.default"
	KEY_CONFIG_DEFAULT_DESC = "keyconfig.default.desc"
	KEY_CONFIG_REDUCED      = "keyconfig.reduced"
	KEY_CONFIG_REDUCED_DESC = "keyconfig.reduced.desc"

	// Dialog
	DIALOG_NEW_PLAYER        = "dialog.newplayer"
	DIALOG_HOW_TO_PLAY       = "dialog.howtoplay"
	DIALOG_BE_SURE_TO_LOGIN  = "dialog.besuretologin"
	DIALOG_BE_SURE_TO_OFFSET = "dialog.besuretooffset"

	// Offset
	OFFSET_INSTRUCTIONS = "offset.instructions"
	OFFSET_INPUT        = "offset.input"

	// States
	//// Play
	STATE_PLAY_RESTART = "state.play.restart"
	STATE_PLAY_PAUSE   = "state.play.pause"

	// Themes
	THEME_STANDARD   = "theme.standard"
	THEME_LEFTBEHIND = "theme.leftbehind"

	// Colors
	COLOR_WHITE = "color.white"
	COLOR_BLACK = "color.black"

	// Note Colors
	NOTE_COLOR_DEFAULT   = "note.color.default"
	NOTE_COLOR_MONO      = "note.color.mono"
	NOTE_COLOR_DUSK      = "note.color.dusk"
	NOTE_COLOR_DAWN      = "note.color.dawn"
	NOTE_COLOR_CUSTOM    = "note.color.custom"
	NOTE_COLOR_AURORA    = "note.color.aurora"
	NOTE_COLOR_ARORUA    = "note.color.arorua"
	NOTE_COLOR_HAMBURGER = "note.color.hamburger"
	NOTE_COLOR_CLASSIC   = "note.color.classic"

	// Erors
	ERROR_USERNAME_REQUIRED   = "error.login.username.required"
	ERROR_PASSWORD_REQUIRED   = "error.login.password.required"
	ERROR_REGISTER_FAIL       = "error.register.fail"
	ERROR_LOGIN_REGISTER_FAIL = "error.login.register.fail"
	ERROR_LOGIN_FAILED        = "error.login.failed"

	// Etc
	UNKNOWN = "unknown"
)

func String(key string) string {
	return assets.GetLocaleString(key)
}
