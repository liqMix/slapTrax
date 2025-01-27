package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/assets"
	"github.com/liqmix/slaptrax/internal/l"
)

type Action int

const (
	ActionBack Action = iota
	ActionSelect
	ActionUp
	ActionDown
	ActionLeft
	ActionRight
	ActionToggleDebug

	// Track activations
	ActionLeftBottom
	ActionLeftTop
	ActionCenterBottom
	ActionCenterTop
	ActionRightBottom
	ActionRightTop

	ActionUnknown
)

var actionToKey = map[Action][]ebiten.Key{
	ActionBack:        {ebiten.KeyEscape, ebiten.KeyF1},
	ActionSelect:      {ebiten.KeyEnter},
	ActionUp:          {ebiten.KeyArrowUp},
	ActionDown:        {ebiten.KeyArrowDown},
	ActionLeft:        {ebiten.KeyArrowLeft},
	ActionRight:       {ebiten.KeyArrowRight},
	ActionToggleDebug: {ebiten.KeyF2},
}

func (a Action) String() string {
	switch a {
	case ActionBack:
		return l.ACTION_BACK
	case ActionSelect:
		return l.ACTION_SELECT
	case ActionUp:
		return l.ACTION_UP
	case ActionDown:
		return l.ACTION_DOWN
	case ActionLeft:
		return l.ACTION_LEFT
	case ActionRight:
		return l.ACTION_RIGHT
	}
	return l.UNKNOWN
}

func (a Action) Key() []ebiten.Key {
	return actionToKey[a]
}

type TrackKeyConfig int

const (
	TrackKeyConfigDefault TrackKeyConfig = iota
	TrackKeyConfigReduced                // Only utilizes the main keyboard, no arrow or nav keys
)

func (t TrackKeyConfig) String() string {
	switch t {
	case TrackKeyConfigDefault:
		return l.KEY_CONFIG_DEFAULT
	case TrackKeyConfigReduced:
		return l.KEY_CONFIG_REDUCED
	}
	return l.UNKNOWN
}

func (t TrackKeyConfig) Image() *ebiten.Image {
	switch t {
	case TrackKeyConfigDefault:
		return assets.GetImage("key_default.png")
	case TrackKeyConfigReduced:
		return assets.GetImage("key_reduced.png")
	}
	return nil
}

func SetTrackKeys(config TrackKeyConfig) {
	switch config {
	case TrackKeyConfigDefault:
		actionToKey[ActionLeftBottom] = []ebiten.Key{
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
		}
		actionToKey[ActionLeftTop] = []ebiten.Key{
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

			ebiten.KeyF2,
			ebiten.KeyF3,
			ebiten.KeyF4,
		}
		actionToKey[ActionCenterTop] = []ebiten.Key{
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

			ebiten.KeyF5,
			ebiten.KeyF6,
			ebiten.KeyF7,
			ebiten.KeyF8,
			ebiten.KeyF9,
			ebiten.KeyF10,
			ebiten.KeyF11,
			ebiten.KeyF12,
		}
		actionToKey[ActionCenterBottom] = []ebiten.Key{
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
		}
		actionToKey[ActionRightTop] = []ebiten.Key{
			ebiten.KeyInsert,
			ebiten.KeyDelete,
			ebiten.KeyHome,
			ebiten.KeyEnd,
			ebiten.KeyPageUp,
			ebiten.KeyPageDown,
		}
		actionToKey[ActionRightBottom] = []ebiten.Key{
			ebiten.KeyArrowLeft,
			ebiten.KeyArrowDown,
			ebiten.KeyArrowRight,
			ebiten.KeyArrowUp,
		}
	case TrackKeyConfigReduced:
		actionToKey[ActionLeftBottom] = []ebiten.Key{
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
		}
		actionToKey[ActionLeftTop] = []ebiten.Key{
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

			ebiten.KeyF2,
			ebiten.KeyF3,
			ebiten.KeyF4,
		}
		actionToKey[ActionCenterTop] = []ebiten.Key{
			ebiten.KeyT,
			ebiten.KeyY,
			ebiten.KeyU,
			ebiten.KeyI,
			ebiten.KeyO,

			ebiten.Key5,
			ebiten.Key6,
			ebiten.Key7,
			ebiten.Key8,
			ebiten.Key9,

			ebiten.KeyF5,
			ebiten.KeyF6,
			ebiten.KeyF7,
			ebiten.KeyF8,
		}
		actionToKey[ActionCenterBottom] = []ebiten.Key{
			ebiten.KeySpace,

			ebiten.KeyV,
			ebiten.KeyB,
			ebiten.KeyN,
			ebiten.KeyM,
			ebiten.KeyComma,

			ebiten.KeyF,
			ebiten.KeyG,
			ebiten.KeyH,
			ebiten.KeyJ,
			ebiten.KeyK,
			ebiten.KeyL,
		}
		actionToKey[ActionRightTop] = []ebiten.Key{
			ebiten.KeyP,
			ebiten.KeyBracketLeft,
			ebiten.KeyBracketRight,
			ebiten.KeyBackslash,

			ebiten.Key0,
			ebiten.KeyMinus,
			ebiten.KeyEqual,
			ebiten.KeyBackspace,

			ebiten.KeyF9,
			ebiten.KeyF10,
			ebiten.KeyF11,
			ebiten.KeyF12,
		}
		actionToKey[ActionRightBottom] = []ebiten.Key{
			ebiten.KeyAltRight,
			ebiten.KeyMetaRight,
			ebiten.KeyControlRight,

			ebiten.KeyPeriod,
			ebiten.KeySlash,
			ebiten.KeyShiftRight,

			ebiten.KeySemicolon,
			ebiten.KeyApostrophe,
			ebiten.KeyEnter,
		}
	}
}
