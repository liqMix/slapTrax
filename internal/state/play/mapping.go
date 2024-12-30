package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type PlayAction string

const (
	RestartSongAction PlayAction = types.L_STATE_PLAY_RESTART
)

var PlayActions = map[PlayAction][]ebiten.Key{
	RestartSongAction: {
		ebiten.KeyF5,
	},
}
var TrackNameToKeys = map[song.TrackName][]ebiten.Key{
	song.LeftBottom: {
		ebiten.KeyControlLeft,
		ebiten.KeyMetaLeft,
		ebiten.KeyAltLeft,

		ebiten.KeyShiftLeft,
		ebiten.KeyZ,
		ebiten.KeyX,
		ebiten.KeyC,
		ebiten.KeyV,

		ebiten.KeyCapsLock,
		ebiten.KeyA,
		ebiten.KeyS,
		ebiten.KeyD,
		ebiten.KeyF,
		ebiten.KeyG,
	},
	song.LeftTop: {
		ebiten.KeyTab,
		ebiten.KeyQ,
		ebiten.KeyW,
		ebiten.KeyE,
		ebiten.KeyR,
		ebiten.KeyT,
		ebiten.KeyY,

		ebiten.KeyBackquote,
		ebiten.Key1,
		ebiten.Key2,
		ebiten.Key3,
		ebiten.Key4,
		ebiten.Key5,
		ebiten.Key6,
		ebiten.Key7,

		// ebiten.KeyEscape,
		// ebiten.KeyF1,
		// ebiten.KeyF2,
		// ebiten.KeyF3,
		// ebiten.KeyF4,
		// ebiten.KeyF5,
		// ebiten.KeyF6,
	},
	song.RightTop: {
		ebiten.KeyU,
		ebiten.KeyI,
		ebiten.KeyO,
		ebiten.KeyP,
		ebiten.KeyBracketLeft,
		ebiten.KeyBracketRight,
		ebiten.KeyBackslash,

		ebiten.Key8,
		ebiten.Key9,
		ebiten.Key0,
		ebiten.KeyMinus,
		ebiten.KeyEqual,
		ebiten.KeyBackspace,

		// ebiten.KeyF7,
		// ebiten.KeyF8,
		// ebiten.KeyF9,
		// ebiten.KeyF10,
		// ebiten.KeyF11,
		// ebiten.KeyF12,
	},
	song.RightBottom: {
		ebiten.KeyAltRight,
		ebiten.KeyMetaRight,
		ebiten.KeyControlRight,

		ebiten.KeyB,
		ebiten.KeyN,
		ebiten.KeyM,
		ebiten.KeyComma,
		ebiten.KeyPeriod,
		ebiten.KeySlash,
		ebiten.KeyShiftRight,

		ebiten.KeyH,
		ebiten.KeyJ,
		ebiten.KeyK,
		ebiten.KeyL,
		ebiten.KeySemicolon,
		ebiten.KeyApostrophe,
		ebiten.KeyEnter,
	},
	song.Center: {
		ebiten.KeySpace,
	},
	song.EdgeTop: {
		ebiten.KeyInsert,
		ebiten.KeyDelete,
		ebiten.KeyHome,
		ebiten.KeyEnd,
		ebiten.KeyPageUp,
		ebiten.KeyPageDown,
	},
	song.EdgeTap1: {
		ebiten.KeyArrowLeft,
	},
	song.EdgeTap2: {
		ebiten.KeyArrowDown,
	},
	song.EdgeTap3: {
		ebiten.KeyArrowRight,
	},
}
