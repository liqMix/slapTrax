//go:build windows
// +build windows

package input

import (
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/logger"
	"golang.org/x/sys/windows/registry"
)

var keyMap = map[uint32]ebiten.Key{
	// Letters
	0x41: ebiten.KeyA,
	0x42: ebiten.KeyB,
	0x43: ebiten.KeyC,
	0x44: ebiten.KeyD,
	0x45: ebiten.KeyE,
	0x46: ebiten.KeyF,
	0x47: ebiten.KeyG,
	0x48: ebiten.KeyH,
	0x49: ebiten.KeyI,
	0x4A: ebiten.KeyJ,
	0x4B: ebiten.KeyK,
	0x4C: ebiten.KeyL,
	0x4D: ebiten.KeyM,
	0x4E: ebiten.KeyN,
	0x4F: ebiten.KeyO,
	0x50: ebiten.KeyP,
	0x51: ebiten.KeyQ,
	0x52: ebiten.KeyR,
	0x53: ebiten.KeyS,
	0x54: ebiten.KeyT,
	0x55: ebiten.KeyU,
	0x56: ebiten.KeyV,
	0x57: ebiten.KeyW,
	0x58: ebiten.KeyX,
	0x59: ebiten.KeyY,
	0x5A: ebiten.KeyZ,

	// Numbers
	0x30: ebiten.Key0,
	0x31: ebiten.Key1,
	0x32: ebiten.Key2,
	0x33: ebiten.Key3,
	0x34: ebiten.Key4,
	0x35: ebiten.Key5,
	0x36: ebiten.Key6,
	0x37: ebiten.Key7,
	0x38: ebiten.Key8,
	0x39: ebiten.Key9,

	// Function keys
	0x70: ebiten.KeyF1,
	0x71: ebiten.KeyF2,
	0x72: ebiten.KeyF3,
	0x73: ebiten.KeyF4,
	0x74: ebiten.KeyF5,
	0x75: ebiten.KeyF6,
	0x76: ebiten.KeyF7,
	0x77: ebiten.KeyF8,
	0x78: ebiten.KeyF9,
	0x79: ebiten.KeyF10,
	0x7A: ebiten.KeyF11,
	0x7B: ebiten.KeyF12,

	// Special keys
	0x08: ebiten.KeyBackspace,
	0x09: ebiten.KeyTab,
	0x0D: ebiten.KeyEnter,
	0x10: ebiten.KeyShift,
	0xA0: ebiten.KeyShift,
	0xA1: ebiten.KeyShift,
	0x11: ebiten.KeyControl,
	0xA2: ebiten.KeyControl,
	0xA3: ebiten.KeyControl,
	0x12: ebiten.KeyAlt,
	0xA4: ebiten.KeyAlt,
	0xA5: ebiten.KeyAlt,
	0x14: ebiten.KeyCapsLock,
	0x1B: ebiten.KeyEscape,
	0x20: ebiten.KeySpace,
	0x21: ebiten.KeyPageUp,
	0x22: ebiten.KeyPageDown,
	0x23: ebiten.KeyEnd,
	0x24: ebiten.KeyHome,
	0x25: ebiten.KeyLeft,
	0x26: ebiten.KeyUp,
	0x27: ebiten.KeyRight,
	0x28: ebiten.KeyDown,
	0x2C: ebiten.KeyPrintScreen,
	0x2D: ebiten.KeyInsert,
	0x2E: ebiten.KeyDelete,
	0x13: ebiten.KeyPause,
	0x91: ebiten.KeyScrollLock,

	// Numpad
	0x60: ebiten.KeyNumpad0,
	0x61: ebiten.KeyNumpad1,
	0x62: ebiten.KeyNumpad2,
	0x63: ebiten.KeyNumpad3,
	0x64: ebiten.KeyNumpad4,
	0x65: ebiten.KeyNumpad5,
	0x66: ebiten.KeyNumpad6,
	0x67: ebiten.KeyNumpad7,
	0x68: ebiten.KeyNumpad8,
	0x69: ebiten.KeyNumpad9,
	0x6A: ebiten.KeyNumpadMultiply,
	0x6B: ebiten.KeyNumpadAdd,
	0x6D: ebiten.KeyNumpadSubtract,
	0x6E: ebiten.KeyNumpadDecimal,
	0x6F: ebiten.KeyNumpadDivide,
	0x90: ebiten.KeyNumLock,

	// Special characters
	0xBA: ebiten.KeySemicolon,
	0xBB: ebiten.KeyEqual,
	0xBC: ebiten.KeyComma,
	0xBD: ebiten.KeyMinus,
	0xBE: ebiten.KeyPeriod,
	0xBF: ebiten.KeySlash,
	0xC0: ebiten.KeyGraveAccent,
	0xDB: ebiten.KeyLeftBracket,
	0xDC: ebiten.KeyBackslash,
	0xDD: ebiten.KeyRightBracket,
	0xDE: ebiten.KeyApostrophe,

	// Windows/Meta keys
	0x5B: ebiten.KeyMeta, // Left Windows
	0x5C: ebiten.KeyMeta, // Right Windows
	0x5D: ebiten.KeyMenu, // Application/Menu key
}

// isTextInputKey checks if a virtual key code represents an alphanumeric or common input key
func isTextInputKey(vkCode uint32) bool {
	// Letters (A-Z)
	if vkCode >= 0x41 && vkCode <= 0x5A {
		return true
	}
	// Numbers (0-9)
	if vkCode >= 0x30 && vkCode <= 0x39 {
		return true
	}
	// Numpad numbers (0-9)
	if vkCode >= 0x60 && vkCode <= 0x69 {
		return true
	}
	// Common punctuation and editing keys
	switch vkCode {
	case 0x08: // Backspace
		return true
	case 0x09: // Tab
		return true
	case 0x0D: // Enter
		return true
	case 0x20: // Space
		return true
	case 0x2E: // Delete
		return true
	case 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF: // ; = , - . /
		return true
	case 0xC0: // `
		return true
	case 0xDB, 0xDC, 0xDD, 0xDE: // [ \ ] '
		return true
	}
	return false
}

// Additional virtual key codes for reference (not mapped to ebiten keys)
// You might want to track these separately if needed:
var systemKeys = map[uint32]string{
	0x2A: "Print",
	0x2F: "Help",
	0x5F: "Sleep",
	0xA6: "BrowserBack",
	0xA7: "BrowserForward",
	0xA8: "BrowserRefresh",
	0xA9: "BrowserStop",
	0xAA: "BrowserSearch",
	0xAB: "BrowserFavorites",
	0xAC: "BrowserHome",
	0xAD: "VolumeMute",
	0xAE: "VolumeDown",
	0xAF: "VolumeUp",
	0xB0: "MediaNext",
	0xB1: "MediaPrev",
	0xB2: "MediaStop",
	0xB3: "MediaPlay",
}

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
	WM_SYSKEYDOWN  = 0x0104
	WM_SYSKEYUP    = 0x0105
	WM_QUIT        = 0x0012
)

var (
	globalKeyboard   *keyboard
	hookCallback     uintptr
	messageThreadID  uintptr
	hookHandle       uintptr
	isShuttingDown   int32
	hookReady        chan bool
	hookShutdownDone chan bool
)

func applyOSHook(k *keyboard) error {
	logger.Debug("Initializing Windows keyboard hook")

	// Disable Windows+L shortcut via registry
	if err := disableWinLShortcut(); err != nil {
		logger.Warn("Failed to disable Windows+L shortcut: %v", err)
		// Continue anyway - the keyboard hook might still help
	}

	// Initialize channels
	hookReady = make(chan bool, 1)
	hookShutdownDone = make(chan bool, 1)

	// Store global reference
	globalKeyboard = k
	atomic.StoreInt32(&isShuttingDown, 0)

	// Start hook in a separate goroutine to avoid blocking
	go hookThread(k)

	// Wait for hook to be ready or timeout
	select {
	case success := <-hookReady:
		if !success {
			// Restore Windows+L shortcut on failure
			enableWinLShortcut()
			return syscall.Errno(1) // Generic error
		}
		logger.Debug("Hook installed successfully")
	case <-time.After(5 * time.Second):
		logger.Error("Timeout waiting for hook installation")
		// Restore Windows+L shortcut on timeout
		enableWinLShortcut()
		return syscall.Errno(1)
	}

	// Setup cleanup handlers
	setupCleanupHandlers(k)

	return nil
}

func hookThread(k *keyboard) {
	// Lock this goroutine to its OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	user32 := syscall.NewLazyDLL("user32.dll")
	kernel32 := syscall.NewLazyDLL("kernel32.dll")

	procSetWindowsHookEx := user32.NewProc("SetWindowsHookExW")
	procGetCurrentThreadId := kernel32.NewProc("GetCurrentThreadId")
	procGetModuleHandle := kernel32.NewProc("GetModuleHandleW")

	// Get current thread ID
	threadID, _, _ := procGetCurrentThreadId.Call()
	messageThreadID = threadID

	// Get module handle
	moduleHandle, _, _ := procGetModuleHandle.Call(0)

	// Create callback
	hookCallback = syscall.NewCallback(lowLevelKeyboardProc)

	// Install hook
	hook, _, err := procSetWindowsHookEx.Call(
		WH_KEYBOARD_LL,
		hookCallback,
		moduleHandle,
		0,
	)

	if hook == 0 {
		logger.Error("Failed to install hook: %v", err)
		hookReady <- false
		return
	}

	hookHandle = hook
	k.osHook = hook

	// Notify that hook is ready
	hookReady <- true

	// Run message loop
	messageLoop()

	// Cleanup after message loop exits
	cleanupHook(k)

	// Notify shutdown is complete
	select {
	case hookShutdownDone <- true:
	default:
	}
}

func lowLevelKeyboardProc(nCode int, wParam, lParam uintptr) uintptr {
	user32 := syscall.NewLazyDLL("user32.dll")
	procCallNextHookEx := user32.NewProc("CallNextHookEx")

	// Always pass through if shutting down
	if atomic.LoadInt32(&isShuttingDown) == 1 {
		ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
		return ret
	}

	// Pass through if we should
	if globalKeyboard == nil || nCode < 0 || !globalKeyboard.isFocused {
		ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
		return ret
	}

	kb := (*struct {
		VkCode      uint32
		ScanCode    uint32
		Flags       uint32
		Time        uint32
		DwExtraInfo uintptr
	})(unsafe.Pointer(lParam))

	// Block ALL keys when focused, not just mapped ones
	// This prevents any system shortcuts from triggering

	// Special handling for critical system combinations
	// You might want to allow Escape or a specific quit combo
	if kb.VkCode == 0x1B { // Escape key
		// Could implement a "hold Escape for 2 seconds to quit" here
	}

	// Check if we should allow text input (for login screen)
	if globalKeyboard.allowTextInput && isTextInputKey(kb.VkCode) {
		logger.Debug("Allowing text input for key VK=0x%X", kb.VkCode)
		
		// Still track the key for our input system
		key, exists := keyMap[kb.VkCode]
		if exists {
			globalKeyboard.m.Lock()
			switch wParam {
			case WM_KEYDOWN, WM_SYSKEYDOWN:
				// Prevent duplicate entries in justPressed
				found := false
				for _, existingKey := range globalKeyboard.justPressed {
					if existingKey == key {
						found = true
						break
					}
				}
				if !found {
					globalKeyboard.justPressed = append(globalKeyboard.justPressed, key)
				}
			case WM_KEYUP, WM_SYSKEYUP:
				// Prevent duplicate entries in justReleased
				found := false
				for _, existingKey := range globalKeyboard.justReleased {
					if existingKey == key {
						found = true
						break
					}
				}
				if !found {
					globalKeyboard.justReleased = append(globalKeyboard.justReleased, key)
				}
			}
			globalKeyboard.m.Unlock()
		}
		
		// Pass through for text input - this allows Ebiten's AppendInputChars to work
		ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
		return ret
	}

	// Map the key if we recognize it
	key, exists := keyMap[kb.VkCode]
	if !exists {
		// If not recognized, check system keys
		if name, ok := systemKeys[kb.VkCode]; ok {
			logger.Debug("Unmapped system key pressed: %s (VK=0x%X)", name, kb.VkCode)
			return 1 // Block the key
		}
		logger.Debug("Unmapped key pressed: VK=0x%X", kb.VkCode)
		return 1 // Block unmapped keys too
	}

	// Handle the key
	globalKeyboard.m.Lock()
	defer globalKeyboard.m.Unlock()

	switch wParam {
	case WM_KEYDOWN, WM_SYSKEYDOWN:
		logger.Debug("Key pressed: VK=0x%X, Key=%s (msg=0x%X)", kb.VkCode, key, wParam)
		// Prevent duplicate entries in justPressed
		for _, existingKey := range globalKeyboard.justPressed {
			if existingKey == key {
				return 1 // Already in the list, just block
			}
		}
		globalKeyboard.justPressed = append(globalKeyboard.justPressed, key)
		return 1 // Block the key
	case WM_KEYUP, WM_SYSKEYUP:
		logger.Debug("Key released: VK=0x%X, Key=%s (msg=0x%X)", kb.VkCode, key, wParam)
		// Prevent duplicate entries in justReleased
		for _, existingKey := range globalKeyboard.justReleased {
			if existingKey == key {
				return 1 // Already in the list, just block
			}
		}
		globalKeyboard.justReleased = append(globalKeyboard.justReleased, key)
		return 1 // Block the key
	}

	// Pass through other messages
	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

func messageLoop() {
	user32 := syscall.NewLazyDLL("user32.dll")
	procGetMessage := user32.NewProc("GetMessageW")
	procTranslateMessage := user32.NewProc("TranslateMessage")
	procDispatchMessage := user32.NewProc("DispatchMessageW")

	type MSG struct {
		Hwnd    uintptr
		Message uint32
		WParam  uintptr
		LParam  uintptr
		Time    uint32
		Pt      struct{ X, Y int32 }
	}

	var msg MSG
	for {
		ret, _, _ := procGetMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			0, 0, 0,
		)

		if ret == 0 { // WM_QUIT
			logger.Debug("Received WM_QUIT")
			break
		}

		if ret == ^uintptr(0) { // Error
			logger.Error("GetMessage error")
			break
		}

		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func removeOSHook(k *keyboard) {
	// Set shutdown flag
	atomic.StoreInt32(&isShuttingDown, 1)

	// Post quit message to the message thread
	if messageThreadID != 0 {
		user32 := syscall.NewLazyDLL("user32.dll")
		procPostThreadMessage := user32.NewProc("PostThreadMessageW")

		ret, _, err := procPostThreadMessage.Call(
			messageThreadID,
			WM_QUIT,
			0,
			0,
		)

		if ret == 0 {
			logger.Error("Failed to post quit message: %v", err)
			// Force cleanup if we can't post the message
			if hookHandle != 0 {
				cleanupHookForce()
			}
		} else {
			// Wait for clean shutdown with timeout
			select {
			case <-hookShutdownDone:
				logger.Debug("Hook shutdown completed cleanly")
			case <-time.After(1 * time.Second):
				logger.Warn("Hook shutdown timeout, forcing cleanup")
				cleanupHookForce()
			}
		}
	}
}

func cleanupHook(k *keyboard) {
	if k.osHook != 0 {
		logger.Debug("Removing Windows keyboard hook")

		user32 := syscall.NewLazyDLL("user32.dll")
		procUnhookWindowsHookEx := user32.NewProc("UnhookWindowsHookEx")

		ret, _, err := procUnhookWindowsHookEx.Call(k.osHook)
		if ret == 0 {
			logger.Error("Failed to unhook: %v", err)
		} else {
			logger.Debug("Hook removed successfully")
		}

		// Restore Windows+L shortcut
		if err := enableWinLShortcut(); err != nil {
			logger.Warn("Failed to restore Windows+L shortcut: %v", err)
		}

		// Clear all global state
		k.osHook = 0
		hookHandle = 0
		globalKeyboard = nil
		hookCallback = 0
		messageThreadID = 0
		
		// Clear the finalizer since we've cleaned up properly
		runtime.SetFinalizer(k, nil)
	}
}

func cleanupHookForce() {
	if hookHandle != 0 {
		user32 := syscall.NewLazyDLL("user32.dll")
		procUnhookWindowsHookEx := user32.NewProc("UnhookWindowsHookEx")

		ret, _, err := procUnhookWindowsHookEx.Call(hookHandle)
		if ret == 0 {
			logger.Error("Failed to force unhook: %v", err)
		} else {
			logger.Debug("Force unhook successful")
		}

		// Restore Windows+L shortcut even on force cleanup
		if err := enableWinLShortcut(); err != nil {
			logger.Warn("Failed to restore Windows+L shortcut during force cleanup: %v", err)
		}

		hookHandle = 0
		globalKeyboard = nil
		hookCallback = 0
		messageThreadID = 0
	}
}

func setupCleanupHandlers(k *keyboard) {
	// Handle SIGINT/SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Debug("Received termination signal")
		removeOSHook(k)
		os.Exit(0)
	}()

	// Windows console handler
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleCtrlHandler := kernel32.NewProc("SetConsoleCtrlHandler")

	consoleHandler := syscall.NewCallback(func(ctrlType uint32) uintptr {
		logger.Debug("Console control event: %d", ctrlType)
		go removeOSHook(k) // Don't block the console handler
		return 1
	})

	procSetConsoleCtrlHandler.Call(consoleHandler, 1)
}

// Cleanup function to call from main game loop
func (k *keyboard) Cleanup() {
	k.cleanup.Do(func() {
		removeOSHook(k)
	})
}

// Alternative: Consider using Raw Input instead
// This is a cleaner approach that only captures input when focused
func useRawInputInstead(k *keyboard) error {
	// Raw Input is less invasive and doesn't block other applications
	// It automatically only captures when your window has focus
	// See: https://docs.microsoft.com/en-us/windows/win32/inputdev/raw-input

	// This would be a better long-term solution
	// Example implementation would go here

	return nil
}

// Registry key management for Windows+L shortcut
var (
	registryKey     = `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`
	registryValue   = "DisableLockWorkstation"
	originalValue   uint32
	hadOriginalKey  bool
)

// disableWinLShortcut disables the Windows+L shortcut via registry
func disableWinLShortcut() error {
	logger.Debug("Disabling Windows+L shortcut via registry")
	
	// Open or create the registry key
	key, _, err := registry.CreateKey(registry.CURRENT_USER, registryKey, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	
	// Check if the value already exists and store original
	val, _, err := key.GetIntegerValue(registryValue)
	if err == nil {
		originalValue = uint32(val)
		hadOriginalKey = true
		logger.Debug("Found existing DisableLockWorkstation value: %d", originalValue)
	} else {
		hadOriginalKey = false
		logger.Debug("No existing DisableLockWorkstation value found")
	}
	
	// Set the value to 1 to disable Windows+L
	err = key.SetDWordValue(registryValue, 1)
	if err != nil {
		return err
	}
	
	logger.Debug("Successfully disabled Windows+L shortcut")
	return nil
}

// enableWinLShortcut restores the Windows+L shortcut via registry
func enableWinLShortcut() error {
	logger.Debug("Restoring Windows+L shortcut via registry")
	
	key, err := registry.OpenKey(registry.CURRENT_USER, registryKey, registry.ALL_ACCESS)
	if err != nil {
		logger.Debug("Could not open registry key for restoration: %v", err)
		return nil // Don't fail if we can't restore
	}
	defer key.Close()
	
	if hadOriginalKey {
		// Restore original value
		err = key.SetDWordValue(registryValue, originalValue)
		if err != nil {
			logger.Error("Failed to restore original DisableLockWorkstation value: %v", err)
		} else {
			logger.Debug("Restored original DisableLockWorkstation value: %d", originalValue)
		}
	} else {
		// Delete the value since it didn't exist before
		err = key.DeleteValue(registryValue)
		if err != nil {
			logger.Error("Failed to delete DisableLockWorkstation value: %v", err)
		} else {
			logger.Debug("Removed DisableLockWorkstation value")
		}
	}
	
	return nil
}
