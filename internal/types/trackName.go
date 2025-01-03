package types

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type TrackName string

const (
	LeftBottom   TrackName = "track.leftbottom"
	LeftTop      TrackName = "track.lefttop"
	CenterBottom TrackName = "track.centerbottom"
	CenterTop    TrackName = "track.topcenter"
	RightBottom  TrackName = "track.rightbottom"
	RightTop     TrackName = "track.righttop"

	Center   TrackName = "track.center"
	EdgeTop  TrackName = "track.edgetop"
	EdgeTap1 TrackName = "track.edgetap1"
	EdgeTap2 TrackName = "track.edgetap2"
	EdgeTap3 TrackName = "track.edgetap3"
)

func (t TrackName) String() string {
	return string(t)
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
	case Center:
		return Yellow
	case CenterBottom:
		return Yellow
	case CenterTop:
		return Yellow
	}
	return White
}
func (t TrackName) NotePairColor() color.RGBA {
	switch t {
	case LeftBottom:
		return Blue
	case LeftTop:
		return Blue
	case RightBottom:
		return Blue
	case RightTop:
		return Blue
	case Center:
		return LightBlue
	case CenterBottom:
		return LightBlue
	case CenterTop:
		return LightBlue
	}
	return White
}

var TrackNameToKeys = map[TrackName][]ebiten.Key{
	LeftBottom: {
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
	LeftTop: {
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
	CenterTop: {
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
	CenterBottom: {
		ebiten.KeySpace,

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
	// Center: {
	// 	ebiten.KeySpace,
	// },
	// EdgeTop: {
	// 	ebiten.KeyInsert,
	// 	ebiten.KeyDelete,
	// 	ebiten.KeyHome,
	// 	ebiten.KeyEnd,
	// 	ebiten.KeyPageUp,
	// 	ebiten.KeyPageDown,
	// },
	// EdgeTap1: {
	// 	ebiten.KeyArrowLeft,
	// },
	// EdgeTap2: {
	// 	ebiten.KeyArrowDown,
	// },
	// EdgeTap3: {
	// 	ebiten.KeyArrowRight,
	// },
}

func TrackNames() []TrackName {
	return []TrackName{
		LeftBottom,
		LeftTop,
		RightBottom,
		RightTop,
		CenterBottom,
		CenterTop,

		// Center,
		// EdgeTop,
		// EdgeTap1,
		// EdgeTap2,
		// EdgeTap3,
	}
}

var MainTracks = []TrackName{
	LeftBottom,
	LeftTop,
	RightBottom,
	RightTop,
	CenterBottom,
	CenterTop,

	// Center,
}

// var EdgeTracks = []TrackName{
// 	EdgeTop,
// 	EdgeTap1,
// 	EdgeTap2,
// 	EdgeTap3,
// }
