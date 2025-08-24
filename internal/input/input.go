package input

import (
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/user"
)

var (
	M      = newMouse()
	K      = newKeyboard()
	Update = func() {
		M.update()
		K.update()
	}
	Close = func() {
		K.close() // This now uses sync.Once, so calling Cleanup() separately is redundant
	}
	JustActioned = func(a Action) bool {
		keys, ok := actionToKey[a]
		if !ok {
			return false
		}
		return K.AreAny(keys, JustPressed)
	}
	IsActioned = func(a Action) bool {
		keys, ok := actionToKey[a]
		if !ok {
			return false
		}
		return K.AreAny(keys, Held)
	}
	NotActioned = func(a Action) bool {
		keys, ok := actionToKey[a]
		if !ok {
			return false
		}
		return !K.AreAny(keys, Held)
	}
)

func InitInput() {
	SetTrackKeys(TrackKeyConfig(user.S().KeyConfig))
}

// SetAllowTextInput enables or disables text input passthrough for login screens
func SetAllowTextInput(allow bool) {
	logger.Debug("Setting allow text input: %v", allow)
	K.SetAllowTextInput(allow)
}

// GetAllowTextInput returns whether text input passthrough is enabled
func GetAllowTextInput() bool {
	return K.GetAllowTextInput()
}
