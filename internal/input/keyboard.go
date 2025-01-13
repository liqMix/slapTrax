package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	bitsPerWord = 64
	numWords    = 4 // Supports up to 256 keys
)

type bimts [numWords]uint64
type keyboard struct {
	pressedBits  bimts
	releasedBits bimts
	heldBits     bimts

	watchedKeys map[ebiten.Key]int64

	runes []rune
}

func getBitPosition(key ebiten.Key) (wordIndex, bitOffset int) {
	k := uint64(key)
	return int(k / bitsPerWord), int(k % bitsPerWord)
}

func newKeyboard() *keyboard {
	return &keyboard{
		watchedKeys: make(map[ebiten.Key]int64),
		runes:       make([]rune, 16),
	}
}

func (k *keyboard) update() {
	// Clear all bits
	for j := range k.pressedBits {
		k.pressedBits[j] = 0
		k.releasedBits[j] = 0
	}

	pressed := inpututil.AppendJustPressedKeys(nil)
	for _, key := range pressed {
		wordIdx, bitOff := getBitPosition(key)
		if wordIdx < numWords {
			k.pressedBits[wordIdx] |= 1 << bitOff
		}
	}

	released := inpututil.AppendJustReleasedKeys(nil)
	for _, key := range released {
		wordIdx, bitOff := getBitPosition(key)
		if wordIdx < numWords {
			k.releasedBits[wordIdx] |= 1 << bitOff
		}
	}

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
