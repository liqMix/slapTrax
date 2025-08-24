package input

import (
	"runtime"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/logger"
)

const (
	bitsPerWord = 64
	numWords    = 4
)

type bimts [numWords]uint64

type keyboard struct {
	justPressed  []ebiten.Key
	justReleased []ebiten.Key

	pressedBits  bimts
	releasedBits bimts
	heldBits     bimts

	watchedKeys map[ebiten.Key]int64

	isFocused      bool
	allowTextInput bool // Allow alphanumeric keys to pass through for text input
	runes          []rune
	osHook         uintptr
	m              sync.RWMutex
	cleanup        sync.Once
}

func getBitPosition(key ebiten.Key) (wordIndex, bitOffset int) {
	k := uint64(key)
	return int(k / bitsPerWord), int(k % bitsPerWord)
}

func newKeyboard() *keyboard {
	keys := &keyboard{
		justPressed:  make([]ebiten.Key, 0, 16),
		justReleased: make([]ebiten.Key, 0, 16),
		watchedKeys:  make(map[ebiten.Key]int64),
		runes:        make([]rune, 16),
	}
	if err := applyOSHook(keys); err != nil {
		logger.Error("Failed to apply hook: %v", err)
	}
	// Last-resort cleanup if all else fails
	runtime.SetFinalizer(keys, func(k *keyboard) {
		if k.osHook != 0 {
			logger.Warn("Keyboard hook not cleaned up properly, finalizer triggered")
			k.cleanup.Do(func() {
				removeOSHook(k)
			})
		}
	})
	return keys
}

func (k *keyboard) update() {
	k.isFocused = ebiten.IsFocused()

	if !k.isFocused {
		return
	}

	// Clear all bits
	for j := range k.pressedBits {
		k.pressedBits[j] = 0
		k.releasedBits[j] = 0
	}

	k.m.Lock()
	for _, key := range k.justPressed {
		wordIdx, bitOff := getBitPosition(key)
		if wordIdx < numWords {
			k.pressedBits[wordIdx] |= 1 << bitOff
		}
	}

	for _, key := range k.justReleased {
		wordIdx, bitOff := getBitPosition(key)
		if wordIdx < numWords {
			k.releasedBits[wordIdx] |= 1 << bitOff
		}
	}

	// Clear just pressed and released keys
	k.justPressed = k.justPressed[:0]
	k.justReleased = k.justReleased[:0]
	k.m.Unlock()

	// Update held bits
	for j := range k.heldBits {
		k.heldBits[j] = (k.heldBits[j] | k.pressedBits[j]) &^ k.releasedBits[j]
	}

	// Update watched keys
	for key := range k.watchedKeys {
		if k.Is(key, JustPressed) {
			k.watchedKeys[key]++
		}
	}
}

func (k *keyboard) close() {
	k.cleanup.Do(func() {
		removeOSHook(k)
	})
}

func (k *keyboard) Runes() []rune {
	return ebiten.AppendInputChars(k.runes[:0])
}

func (k *keyboard) Get(s Status) []ebiten.Key {
	var keys []ebiten.Key
	for i := 0; i < numWords; i++ {
		for j := 0; j < bitsPerWord; j++ {
			key := ebiten.Key(i*bitsPerWord + j)
			if k.Is(key, s) {
				keys = append(keys, key)
			}
		}
	}
	return keys
}

func (k *keyboard) ForceReset() {
	for i := range k.pressedBits {
		k.pressedBits[i] = 0
		k.releasedBits[i] = 0
		k.heldBits[i] = 0
	}
}

func (k *keyboard) Is(key ebiten.Key, s Status) bool {
	wordIdx, bitOff := getBitPosition(key)
	if wordIdx >= numWords {
		return false
	}

	var bimts *bimts
	switch s {
	case JustPressed:
		bimts = &k.pressedBits
	case JustReleased:
		bimts = &k.releasedBits
	case Held:
		bimts = &k.heldBits
	}
	if bimts == nil {
		return false
	}

	return (bimts[wordIdx] & (1 << bitOff)) != 0
}

func (k *keyboard) AreAny(key []ebiten.Key, s Status) bool {
	if len(key) == 0 {
		return false
	}

	for _, key := range key {
		if k.Is(key, s) {
			return true
		}
	}
	return false
}

func (k *keyboard) AreAll(key []ebiten.Key, s Status) bool {
	if len(key) == 0 {
		return false
	}

	for _, key := range key {
		if !k.Is(key, s) {
			return false
		}
	}
	return true
}

func (k *keyboard) WatchKeys(keys []ebiten.Key) {
	k.watchedKeys = make(map[ebiten.Key]int64)
	for _, key := range keys {
		k.watchedKeys[key] = 0
	}
}

func (k *keyboard) ClearWatchedKeys() {
	k.watchedKeys = nil
}

func (k *keyboard) IsKeyHeldFor(key ebiten.Key, frames int64) bool {
	if k.watchedKeys == nil {
		return false
	}

	if frames < 0 || frames == 0 {
		return false
	}

	if _, ok := k.watchedKeys[key]; ok {
		return k.watchedKeys[key] >= frames
	}

	return false
}

// SetAllowTextInput enables or disables text input passthrough
func (k *keyboard) SetAllowTextInput(allow bool) {
	k.m.Lock()
	defer k.m.Unlock()
	k.allowTextInput = allow
}

// GetAllowTextInput returns whether text input passthrough is enabled
func (k *keyboard) GetAllowTextInput() bool {
	k.m.RLock()
	defer k.m.RUnlock()
	return k.allowTextInput
}
