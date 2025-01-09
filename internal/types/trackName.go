package types

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

type TrackName int

// Order is critical here
const (
	TrackLeftBottom TrackName = iota
	TrackLeftTop
	TrackCenterBottom
	TrackCenterTop
	TrackRightBottom
	TrackRightTop
)

func TrackNames() []TrackName {
	return []TrackName{
		TrackLeftBottom,
		TrackLeftTop,
		TrackRightBottom,
		TrackRightTop,
		TrackCenterBottom,
		TrackCenterTop,
	}
}

func (t TrackName) String() string {
	switch t {
	case TrackLeftBottom:
		return "LeftBottom"
	case TrackLeftTop:
		return "LeftTop"
	case TrackRightBottom:
		return "RightBottom"
	case TrackRightTop:
		return "RightTop"
	case TrackCenterBottom:
		return "CenterBottom"
	case TrackCenterTop:
		return "CenterTop"
	}
	return "Unknown"
}

func (t TrackName) NoteColor() color.RGBA {
	return TrackTypeFromName(t).Color()
}

type TrackType int

const (
	TrackTypeCenter TrackType = iota
	TrackTypeCorner
)

func TrackTypeFromName(n TrackName) TrackType {
	switch n {
	case TrackCenterBottom, TrackCenterTop:
		return TrackTypeCenter
	}
	return TrackTypeCorner
}

func (t TrackType) Color() color.RGBA {
	theme := NoteColorTheme(user.S().NoteColorTheme)
	switch t {
	case TrackTypeCenter:
		return theme.CenterColor()
	case TrackTypeCorner:
		return theme.CornerColor()
	}
	return White.C()
}

// This one has larger center tracks
var TrackNameToKeys = map[TrackName][]ebiten.Key{
	TrackLeftBottom: {
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
	TrackLeftTop: {
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
	TrackCenterTop: {
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
	TrackCenterBottom: {
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

	TrackRightTop: {
		ebiten.KeyInsert,
		ebiten.KeyDelete,
		ebiten.KeyHome,
		ebiten.KeyEnd,
		ebiten.KeyPageUp,
		ebiten.KeyPageDown,
	},

	TrackRightBottom: {
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
