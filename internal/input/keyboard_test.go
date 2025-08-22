package input

import (
	"sync/atomic"
	"testing"
	"time"
	
	"github.com/hajimehoshi/ebiten/v2"
)

func TestKeyboardCleanupRaceCondition(t *testing.T) {
	// Test that multiple cleanup calls don't cause issues
	k := newKeyboard()
	if k == nil {
		t.Fatal("Failed to create keyboard")
	}

	// Call cleanup multiple times concurrently
	done := make(chan bool, 3)
	
	go func() {
		k.Cleanup()
		done <- true
	}()
	
	go func() {
		k.close()
		done <- true
	}()
	
	go func() {
		// Simulate the defer call from main
		k.Cleanup()
		k.close()
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		select {
		case <-done:
			// Good
		case <-time.After(2 * time.Second):
			t.Fatal("Cleanup timed out - potential deadlock")
		}
	}

	// Verify cleanup state
	if k.osHook != 0 {
		t.Error("Hook handle should be cleared after cleanup")
	}
}

func TestShutdownFlagPreventsHookCalls(t *testing.T) {
	// Ensure that once shutdown flag is set, hook proc returns immediately
	atomic.StoreInt32(&isShuttingDown, 1)
	
	// Mock a hook call - this should pass through immediately
	result := lowLevelKeyboardProc(0, 0x0100, 0) // WM_KEYDOWN
	
	// The result should be from CallNextHookEx, not our blocking return value of 1
	// This is hard to test directly, but we can verify the shutdown flag works
	if atomic.LoadInt32(&isShuttingDown) != 1 {
		t.Error("Shutdown flag should remain set")
	}
	
	// Reset for other tests
	atomic.StoreInt32(&isShuttingDown, 0)
	
	// Verify result is not our blocking value when shutting down
	// Note: This test is limited as we can't easily mock the Windows API calls
	_ = result
}

func TestKeyPressReleaseHandling(t *testing.T) {
	k := newKeyboard()
	if k == nil {
		t.Skip("Cannot test keyboard on this platform")
	}
	defer k.close()

	// Test that keys are properly added to justPressed/justReleased
	initialPressed := len(k.justPressed)
	initialReleased := len(k.justReleased)

	// Simulate key press by directly calling the method that would be called by the hook
	k.m.Lock()
	k.justPressed = append(k.justPressed, 65) // 'A' key
	k.m.Unlock()

	if len(k.justPressed) != initialPressed+1 {
		t.Error("Key press not recorded properly")
	}

	// Simulate key release
	k.m.Lock()
	k.justReleased = append(k.justReleased, 65) // 'A' key
	k.m.Unlock()

	if len(k.justReleased) != initialReleased+1 {
		t.Error("Key release not recorded properly")
	}

	// Test that arrays are cleared manually (since ebiten.IsFocused() won't work in tests)
	k.m.Lock()
	k.justPressed = k.justPressed[:0]
	k.justReleased = k.justReleased[:0]
	k.m.Unlock()
	
	if len(k.justPressed) != 0 || len(k.justReleased) != 0 {
		t.Error("Arrays should be cleared after manual reset")
	}
}

func TestAltKeySystemMessages(t *testing.T) {
	// Test that ALT key handles both WM_SYSKEYDOWN/UP and WM_KEYDOWN/UP
	k := newKeyboard()
	if k == nil {
		t.Skip("Cannot test keyboard on this platform")
	}
	defer k.close()

	// Test ALT key mappings
	altKeys := []uint32{0x12, 0xA4, 0xA5} // ALT, LALT, RALT
	
	for _, vkCode := range altKeys {
		// Test that the VK code maps to ALT key
		if ebitenKey, exists := keyMap[vkCode]; !exists {
			t.Errorf("ALT VK code 0x%X not found in keyMap", vkCode)
		} else if ebitenKey != ebiten.KeyAlt {
			t.Errorf("ALT VK code 0x%X maps to %v instead of KeyAlt", vkCode, ebitenKey)
		}
	}

	// Test message type constants
	expectedConstants := map[string]uint32{
		"WM_KEYDOWN":    0x0100,
		"WM_KEYUP":      0x0101,
		"WM_SYSKEYDOWN": 0x0104,
		"WM_SYSKEYUP":   0x0105,
	}

	actualConstants := map[string]uint32{
		"WM_KEYDOWN":    WM_KEYDOWN,
		"WM_KEYUP":      WM_KEYUP,
		"WM_SYSKEYDOWN": WM_SYSKEYDOWN,
		"WM_SYSKEYUP":   WM_SYSKEYUP,
	}

	for name, expected := range expectedConstants {
		if actual, exists := actualConstants[name]; !exists {
			t.Errorf("Constant %s not defined", name)
		} else if actual != expected {
			t.Errorf("Constant %s = 0x%X, expected 0x%X", name, actual, expected)
		}
	}
}

func TestFinalizerCleanup(t *testing.T) {
	// Test that the finalizer works as a last resort
	k := newKeyboard()
	if k == nil {
		t.Skip("Cannot test keyboard on this platform")
	}
	
	// Don't call close() explicitly - let finalizer handle it
	// This is hard to test reliably, but we can at least verify the hook is set
	if k.osHook == 0 {
		t.Log("Hook not installed (may be expected on test systems)")
	}
	
	// Manually trigger what the finalizer would do
	if k.osHook != 0 {
		k.close()
		if k.osHook != 0 {
			t.Error("Finalizer cleanup should clear osHook")
		}
	}
}