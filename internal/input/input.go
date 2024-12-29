package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	bitsPerWord = 64
	numWords    = 4 // Supports up to 256 keys
)

var pressedBits [numWords]uint64
var releasedBits [numWords]uint64
var heldBits [numWords]uint64

func getBitPosition(key ebiten.Key) (wordIndex, bitOffset int) {
	k := uint64(key)
	return int(k / bitsPerWord), int(k % bitsPerWord)
}

func Update() {
	// Clear all bits
	for j := range pressedBits {
		pressedBits[j] = 0
		releasedBits[j] = 0
	}

	pressed := inpututil.AppendJustPressedKeys(nil)
	for _, key := range pressed {
		wordIdx, bitOff := getBitPosition(key)
		if wordIdx < numWords {
			pressedBits[wordIdx] |= 1 << bitOff
		}
	}

	released := inpututil.AppendJustReleasedKeys(nil)
	for _, key := range released {
		wordIdx, bitOff := getBitPosition(key)
		if wordIdx < numWords {
			releasedBits[wordIdx] |= 1 << bitOff
		}
	}

	// Update held bits
	for j := range heldBits {
		heldBits[j] = (heldBits[j] | pressedBits[j]) &^ releasedBits[j]
	}
}

func IsKeyJustPressed(key ebiten.Key) bool {
	wordIdx, bitOff := getBitPosition(key)
	if wordIdx >= numWords {
		return false
	}
	return (pressedBits[wordIdx] & (1 << bitOff)) != 0
}

func IsKeyJustReleased(key ebiten.Key) bool {
	wordIdx, bitOff := getBitPosition(key)
	if wordIdx >= numWords {
		return false
	}
	return (releasedBits[wordIdx] & (1 << bitOff)) != 0
}

func IsKeyHeld(key ebiten.Key) bool {
	wordIdx, bitOff := getBitPosition(key)
	if wordIdx >= numWords {
		return false
	}
	return (heldBits[wordIdx] & (1 << bitOff)) != 0
}

func AnyKeysJustPressed(keys []ebiten.Key) bool {
	if len(keys) == 0 {
		return false
	}

	for _, key := range keys {
		if IsKeyJustPressed(key) {
			return true
		}
	}
	return false
}

func AnyKeysJustReleased(keys []ebiten.Key) bool {
	if len(keys) == 0 {
		return false
	}

	for _, key := range keys {
		if IsKeyJustReleased(key) {
			return true
		}
	}
	return false
}

func AllKeysReleased(keys []ebiten.Key) bool {
	if len(keys) == 0 {
		return false
	}

	for _, key := range keys {
		if IsKeyHeld(key) {
			return false
		}
	}
	return true
}

func AnyKeysPressed(keys []ebiten.Key) bool {
	if len(keys) == 0 {
		return false
	}

	for _, key := range keys {
		if IsKeyHeld(key) {
			return true
		}
	}
	return false
}
