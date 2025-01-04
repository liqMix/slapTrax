package types

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type TrackName int

// Order is critical here
const (
	LeftBottom TrackName = iota
	LeftTop
	CenterBottom
	CenterTop
	RightBottom
	RightTop
)

func TrackNames() []TrackName {
	return []TrackName{
		LeftBottom,
		LeftTop,
		RightBottom,
		RightTop,
		CenterBottom,
		CenterTop,
	}
}

func (t TrackName) String() string {
	switch t {
	case LeftBottom:
		return "LeftBottom"
	case LeftTop:
		return "LeftTop"
	case RightBottom:
		return "RightBottom"
	case RightTop:
		return "RightTop"
	case CenterBottom:
		return "CenterBottom"
	case CenterTop:
		return "CenterTop"
	}
	return "Unknown"
}

func (t TrackName) NoteColor() color.RGBA {
	switch t {
	case LeftBottom:
		return Orange
	case LeftTop:
		return Orange
	case RightBottom:
		return Orange
	case RightTop:
		return Orange
	case CenterBottom:
		return Yellow
	case CenterTop:
		return Yellow
	}
	return White
}

// Hmm..
// func (t TrackName) NotePairColor() color.RGBA {
// 	switch t {
// 	case LeftBottom:
// 		return Blue
// 	case LeftTop:
// 		return Blue
// 	case RightBottom:
// 		return Blue
// 	case RightTop:
// 		return Blue
// 	case CenterBottom:
// 		return LightBlue
// 	case CenterTop:
// 		return LightBlue
// 	}
// 	return White
// }

// This one has larger center tracks
var TrackNameToKeys = map[TrackName][]ebiten.Key{
	LeftBottom: {
		ebiten.KeyControlLeft,
		ebiten.KeyMetaLeft,
		ebiten.KeyAltLeft,

		ebiten.KeyShiftLeft,
		ebiten.KeyZ,
		ebiten.KeyX,
		ebiten.KeyC,

		ebiten.KeyCapsLock,
		ebiten.KeyA,
		ebiten.KeyS,
		ebiten.KeyD,
	},
	LeftTop: {
		ebiten.KeyTab,
		ebiten.KeyQ,
		ebiten.KeyW,
		ebiten.KeyE,
		ebiten.KeyR,

		ebiten.KeyBackquote,
		ebiten.Key1,
		ebiten.Key2,
		ebiten.Key3,
		ebiten.Key4,
	},
	CenterTop: {
		ebiten.KeyT,
		ebiten.KeyY,
		ebiten.KeyU,
		ebiten.KeyI,
		ebiten.KeyO,
		ebiten.KeyP,
		ebiten.KeyBracketLeft,
		ebiten.KeyBracketRight,
		ebiten.KeyBackslash,

		ebiten.Key5,
		ebiten.Key6,
		ebiten.Key7,
		ebiten.Key8,
		ebiten.Key9,
		ebiten.Key0,
		ebiten.KeyMinus,
		ebiten.KeyEqual,
		ebiten.KeyBackspace,
	},
	CenterBottom: {
		ebiten.KeySpace,

		ebiten.KeyAltRight,
		ebiten.KeyMetaRight,
		ebiten.KeyControlRight,

		ebiten.KeyV,
		ebiten.KeyB,
		ebiten.KeyN,
		ebiten.KeyM,
		ebiten.KeyComma,
		ebiten.KeyPeriod,
		ebiten.KeySlash,
		ebiten.KeyShiftRight,

		ebiten.KeyF,
		ebiten.KeyG,
		ebiten.KeyH,
		ebiten.KeyJ,
		ebiten.KeyK,
		ebiten.KeyL,
		ebiten.KeySemicolon,
		ebiten.KeyApostrophe,
		ebiten.KeyEnter,
	},

	RightTop: {
		ebiten.KeyInsert,
		ebiten.KeyDelete,
		ebiten.KeyHome,
		ebiten.KeyEnd,
		ebiten.KeyPageUp,
		ebiten.KeyPageDown,
	},

	RightBottom: {
		ebiten.KeyArrowLeft,
		ebiten.KeyArrowDown,
		ebiten.KeyArrowRight,
		ebiten.KeyArrowUp,
	},
}

//// Keys set from old style
//// Left side is larger
// var TrackNameToKeys = map[TrackName][]ebiten.Key{
// 	LeftBottom: {
// 		ebiten.KeyControlLeft,
// 		ebiten.KeyMetaLeft,
// 		ebiten.KeyAltLeft,

// 		ebiten.KeyShiftLeft,
// 		ebiten.KeyZ,
// 		ebiten.KeyX,
// 		ebiten.KeyC,
// 		ebiten.KeyV,

// 		ebiten.KeyCapsLock,
// 		ebiten.KeyA,
// 		ebiten.KeyS,
// 		ebiten.KeyD,
// 		ebiten.KeyF,
// 		ebiten.KeyG,
// 	},
// 	LeftTop: {
// 		ebiten.KeyTab,
// 		ebiten.KeyQ,
// 		ebiten.KeyW,
// 		ebiten.KeyE,
// 		ebiten.KeyR,
// 		ebiten.KeyT,
// 		ebiten.KeyY,

// 		ebiten.KeyBackquote,
// 		ebiten.Key1,
// 		ebiten.Key2,
// 		ebiten.Key3,
// 		ebiten.Key4,
// 		ebiten.Key5,
// 		ebiten.Key6,
// 		ebiten.Key7,
// 	},
// 	CenterTop: {
// 		ebiten.KeyU,
// 		ebiten.KeyI,
// 		ebiten.KeyO,
// 		ebiten.KeyP,
// 		ebiten.KeyBracketLeft,
// 		ebiten.KeyBracketRight,
// 		ebiten.KeyBackslash,

// 		ebiten.Key8,
// 		ebiten.Key9,
// 		ebiten.Key0,
// 		ebiten.KeyMinus,
// 		ebiten.KeyEqual,
// 		ebiten.KeyBackspace,
// 	},
// 	CenterBottom: {
// 		ebiten.KeySpace,

// 		ebiten.KeyAltRight,
// 		ebiten.KeyMetaRight,
// 		ebiten.KeyControlRight,

// 		ebiten.KeyB,
// 		ebiten.KeyN,
// 		ebiten.KeyM,
// 		ebiten.KeyComma,
// 		ebiten.KeyPeriod,
// 		ebiten.KeySlash,
// 		ebiten.KeyShiftRight,

// 		ebiten.KeyH,
// 		ebiten.KeyJ,
// 		ebiten.KeyK,
// 		ebiten.KeyL,
// 		ebiten.KeySemicolon,
// 		ebiten.KeyApostrophe,
// 		ebiten.KeyEnter,
// 	},

// 	RightTop: {
// 		ebiten.KeyInsert,
// 		ebiten.KeyDelete,
// 		ebiten.KeyHome,
// 		ebiten.KeyEnd,
// 		ebiten.KeyPageUp,
// 		ebiten.KeyPageDown,
// 	},

// 	RightBottom: {
// 		ebiten.KeyArrowLeft,
// 		ebiten.KeyArrowDown,
// 		ebiten.KeyArrowRight,
// 		ebiten.KeyArrowUp,
// 	},
// }
