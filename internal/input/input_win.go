//go:build windows
// +build windows

package input

import (
	"errors"
	"os"
	"os/signal"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"syscall"
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
	0x10: ebiten.KeyShift,        // Generic shift
	0xA0: ebiten.KeyShiftLeft,    // Left Shift
	0xA1: ebiten.KeyShiftRight,   // Right Shift
	0x11: ebiten.KeyControl,      // Generic control
	0xA2: ebiten.KeyControlLeft,  // Left Control
	0xA3: ebiten.KeyControlRight, // Right Control
	0x12: ebiten.KeyAlt,          // Generic alt
	0xA4: ebiten.KeyAltLeft,      // Left Alt
	0xA5: ebiten.KeyAltRight,     // Right Alt
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
	0x5B: ebiten.KeyMetaLeft,  // Left Windows
	0x5C: ebiten.KeyMetaRight, // Right Windows
	0x5D: ebiten.KeyMenu,      // Application/Menu key
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

	registryPath = `Software\Microsoft\Windows\CurrentVersion\Policies\System`
)

type HHOOK uintptr

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

type KeyboardManager struct {
	hook              HHOOK
	msgLoopDone       chan bool
	cleanupOnce       sync.Once
	registryMutex     sync.Mutex
	originalWinLValue uint32
	threadID          uintptr
}

type KeyEvent struct {
	VkCode uint32
	IsDown bool
}

var (
	globalKeyboard *keyboard
	globalManager  *KeyboardManager
	isShuttingDown int32
)

func (km *KeyboardManager) DisableWinL() error {
	km.registryMutex.Lock()
	defer km.registryMutex.Unlock()

	// Try multiple registry approaches for Win+L blocking
	approaches := []struct {
		root  registry.Key
		path  string
		value string
	}{
		{registry.CURRENT_USER, registryPath, "DisableLockWorkstation"},
		{registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Policies\Explorer`, "NoWinKeys"},
		{registry.LOCAL_MACHINE, registryPath, "DisableLockWorkstation"},
	}

	var lastErr error
	for _, approach := range approaches {
		key, _, err := registry.CreateKey(approach.root, approach.path, registry.SET_VALUE)
		if err != nil {
			lastErr = err
			logger.Debug("Failed to open registry %s\\%s: %v", approach.root, approach.path, err)
			continue
		}

		// Store original value first
		originalValue, _, err := key.GetIntegerValue(approach.value)
		if err != nil {
			originalValue = 0 // Default if not exists
		}
		km.originalWinLValue = uint32(originalValue)

		// Set to 1 to disable
		err = key.SetDWordValue(approach.value, 1)
		key.Close()

		if err == nil {
			logger.Debug("Successfully disabled Win+L using %s\\%s\\%s", approach.root, approach.path, approach.value)
			return nil
		}

		lastErr = err
		logger.Debug("Failed to set registry value %s\\%s\\%s: %v", approach.root, approach.path, approach.value, err)
	}

	// If registry approach fails, we'll rely on the hook to block Win+L
	logger.Warn("Could not disable Win+L via registry (tried all approaches): %v", lastErr)
	return lastErr
}

func (km *KeyboardManager) RestoreWinL() error {
	km.registryMutex.Lock()
	defer km.registryMutex.Unlock()

	key, err := registry.OpenKey(registry.CURRENT_USER, registryPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	if km.originalWinLValue == 0 {
		// Delete the value to restore default
		return key.DeleteValue("DisableLockWorkstation")
	}
	return key.SetDWordValue("DisableLockWorkstation", km.originalWinLValue)
}

func isGameActive() bool {
	if globalKeyboard == nil || !globalKeyboard.isFocused {
		return false
	}
	currentState := globalKeyboard.GetCurrentState()
	return currentState == "state.play" || currentState == "state.play.pause"
}

func applyOSHook(k *keyboard) error {
	logger.Debug("Initializing Windows keyboard hook")

	// Store global reference
	globalKeyboard = k
	atomic.StoreInt32(&isShuttingDown, 0)

	// Create keyboard manager
	globalManager = &KeyboardManager{
		msgLoopDone: make(chan bool, 1),
	}

	// Initialize the keyboard manager
	err := globalManager.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (km *KeyboardManager) Initialize() error {
	// 1. Set up registry modification
	if err := km.DisableWinL(); err != nil {
		// Log warning but continue - not critical
		logger.Warn("Could not disable Win+L: %v", err)
	}

	// 2. Install hook in dedicated thread
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		// Initialize DLLs locally in this goroutine to avoid race conditions
		user32 := syscall.NewLazyDLL("user32.dll")
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		
		procGetCurrentThreadId := kernel32.NewProc("GetCurrentThreadId")
		procGetModuleHandle := kernel32.NewProc("GetModuleHandleW")
		procSetWindowsHookExW := user32.NewProc("SetWindowsHookExW")
		procGetMessage := user32.NewProc("GetMessageW")
		procTranslateMessage := user32.NewProc("TranslateMessage")
		procDispatchMessage := user32.NewProc("DispatchMessageW")

		// Get current thread ID
		threadID, _, _ := procGetCurrentThreadId.Call()
		km.threadID = threadID

		// Get module handle
		moduleHandle, _, _ := procGetModuleHandle.Call(0)

		callback := syscall.NewCallback(lowLevelKeyboardProc)
		hook, _, err := procSetWindowsHookExW.Call(
			WH_KEYBOARD_LL,
			callback,
			moduleHandle,
			0) // Global hook

		if hook == 0 {
			logger.Error("Failed to install hook: %v", err)
			km.msgLoopDone <- false
			return
		}

		km.hook = HHOOK(hook)
		globalKeyboard.osHook = hook
		km.msgLoopDone <- true

		// Message loop
		var msg MSG
		for {
			ret, _, _ := procGetMessage.Call(
				uintptr(unsafe.Pointer(&msg)),
				0, 0, 0)

			if ret == 0 || ret == ^uintptr(0) { // WM_QUIT or error
				break
			}

			procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}()

	// 3. Wait for initialization
	success := <-km.msgLoopDone
	if !success {
		return errors.New("failed to install hook")
	}

	// 4. Set up cleanup handlers
	km.setupCleanupHandlers()

	return nil
}

func (km *KeyboardManager) setupCleanupHandlers() {
	// Cleanup on signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Debug("Received termination signal")
		km.Cleanup()
		os.Exit(0)
	}()

	// Windows console handler
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleCtrlHandler := kernel32.NewProc("SetConsoleCtrlHandler")

	consoleHandler := syscall.NewCallback(func(ctrlType uint32) uintptr {
		logger.Debug("Console control event: %d", ctrlType)
		go km.Cleanup() // Don't block the console handler
		return 1
	})

	procSetConsoleCtrlHandler.Call(consoleHandler, 1)

	// Set finalizer as last resort
	runtime.SetFinalizer(km, func(km *KeyboardManager) {
		km.Cleanup()
	})
}

func (km *KeyboardManager) Cleanup() {
	km.cleanupOnce.Do(func() {
		// Initialize DLLs locally to avoid race conditions
		user32 := syscall.NewLazyDLL("user32.dll")
		procUnhookWindowsHookEx := user32.NewProc("UnhookWindowsHookEx")
		procPostThreadMessage := user32.NewProc("PostThreadMessageW")

		// 1. Remove hook
		if km.hook != 0 {
			procUnhookWindowsHookEx.Call(uintptr(km.hook))
			logger.Debug("Hook removed")
		}

		// 2. Restore registry
		if err := km.RestoreWinL(); err != nil {
			logger.Error("Failed to restore Win+L: %v", err)
		}

		// 3. Post quit message to stop message loop
		if km.threadID != 0 {
			procPostThreadMessage.Call(km.threadID, WM_QUIT, 0, 0)
		}

		// 4. Clear global state
		if globalKeyboard != nil {
			globalKeyboard.osHook = 0
		}
		km.hook = 0

		logger.Debug("Keyboard manager cleanup completed")
	})
}

func lowLevelKeyboardProc(nCode int, wParam, lParam uintptr) uintptr {
	// Initialize DLLs locally to avoid race conditions
	user32 := syscall.NewLazyDLL("user32.dll")
	procCallNextHookEx := user32.NewProc("CallNextHookEx")

	// Fast path: pass through if not in game or shutting down
	if atomic.LoadInt32(&isShuttingDown) == 1 || nCode < 0 {
		ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
		return ret
	}

	kb := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

	// Always track keys synchronously to prevent stuck keys
	if globalKeyboard != nil {
		trackKeyForInputSystem(kb.VkCode, wParam)
	}

	// Special handling for Win+L regardless of game state
	if isGameActive() {
		// Use the same user32 reference for GetAsyncKeyState
		procGetAsyncKeyState := user32.NewProc("GetAsyncKeyState")
		
		// Check for Windows key + L combination
		if kb.VkCode == 0x4C { // L key
			leftWin, _, _ := procGetAsyncKeyState.Call(uintptr(0x5B))  // Left Windows
			rightWin, _, _ := procGetAsyncKeyState.Call(uintptr(0x5C)) // Right Windows

			if (leftWin&0x8000) != 0 || (rightWin&0x8000) != 0 {
				logger.Debug("Blocking Win+L during gameplay")
				return 1 // Block Win+L
			}
		}

		// Check for Windows key when L is already pressed
		if kb.VkCode == 0x5B || kb.VkCode == 0x5C { // Windows keys
			lKey, _, _ := procGetAsyncKeyState.Call(uintptr(0x4C)) // L key

			if (lKey & 0x8000) != 0 {
				logger.Debug("Blocking Windows key when L is pressed during gameplay")
				return 1 // Block Windows key when L is pressed
			}
		}
	}

	// Check if we should be blocking keys
	if !isGameActive() {
		// Not in game - allow keys through
		ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
		return ret
	}

	// Allow text input when enabled (for login screen)
	if globalKeyboard != nil && globalKeyboard.allowTextInput && isTextInputKey(kb.VkCode) {
		ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
		return ret
	}

	// During gameplay, block ALL keys except Ctrl+Alt+Del
	// (Ctrl+Alt+Del cannot be blocked and will pass through automatically)
	return 1 // Block the key
}

func trackKeyForInputSystem(vkCode uint32, wParam uintptr) {
	key, exists := keyMap[vkCode]
	if !exists {
		return // Don't track unmapped keys
	}

	globalKeyboard.m.Lock()
	defer globalKeyboard.m.Unlock()

	switch wParam {
	case WM_KEYDOWN, WM_SYSKEYDOWN:
		// Add to justPressed if not already pressed
		wordIdx, bitOff := getBitPosition(key)
		if wordIdx < numWords {
			wasHeld := (globalKeyboard.heldBits[wordIdx] & (1 << bitOff)) != 0
			if !wasHeld {
				// Check for duplicates before adding
				isDuplicate := slices.Contains(globalKeyboard.justPressed, key)
				if !isDuplicate {
					globalKeyboard.justPressed = append(globalKeyboard.justPressed, key)
					globalKeyboard.heldBits[wordIdx] |= (1 << bitOff)
				}
			}
		}
	case WM_KEYUP, WM_SYSKEYUP:
		// Add to justReleased if was held
		wordIdx, bitOff := getBitPosition(key)
		if wordIdx < numWords {
			wasHeld := (globalKeyboard.heldBits[wordIdx] & (1 << bitOff)) != 0
			if wasHeld {
				// Check for duplicates before adding
				isDuplicate := slices.Contains(globalKeyboard.justReleased, key)
				if !isDuplicate {
					globalKeyboard.justReleased = append(globalKeyboard.justReleased, key)
					globalKeyboard.heldBits[wordIdx] &^= (1 << bitOff)
				}
			}
		}
	}
}

func removeOSHook(k *keyboard) {
	// Set shutdown flag
	atomic.StoreInt32(&isShuttingDown, 1)

	if globalManager != nil {
		globalManager.Cleanup()
	}
}

// Cleanup function to call from main game loop
func (k *keyboard) Cleanup() {
	k.cleanup.Do(func() {
		removeOSHook(k)
	})
}
