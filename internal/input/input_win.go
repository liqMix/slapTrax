//go:build windows
// +build windows

package input

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"slices"
	"sync"
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
	WM_HOTKEY      = 0x0312

	registryPath        = `Software\Microsoft\Windows\CurrentVersion\Policies\System`
	gameBarRegistryPath = `SOFTWARE\Microsoft\Windows\CurrentVersion\GameDVR`
	gameConfigStorePath = `System\GameConfigStore`

	// RegisterHotKey modifiers
	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004
	MOD_WIN     = 0x0008

	// Hotkey IDs for our registered hotkeys
	HOTKEY_WIN_ALT_G = 1001
	HOTKEY_WIN_ALT_W = 1002
	HOTKEY_WIN_ALT_B = 1003
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

	// GameBar registry backup values
	originalGameDVREnabled    uint32
	originalAppCaptureEnabled uint32
	originalGameConfigEnabled uint32
	gameBarRegistryModified   bool

	// Hotkey registration tracking
	registeredHotkeys []int
	hotkeyMutex       sync.Mutex
}

type KeyEvent struct {
	VkCode uint32
	IsDown bool
}

var (
	globalKeyboard *keyboard
	globalManager  *KeyboardManager
	isShuttingDown int32

	// Modifier key states for combo detection
	winKeyDown  int32
	altKeyDown  int32
	ctrlKeyDown int32

	// BlockInput state
	inputBlocked int32
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

func (km *KeyboardManager) DisableGameBar() error {
	km.registryMutex.Lock()
	defer km.registryMutex.Unlock()

	var errors []error

	// 1. Disable GameDVR in main GameDVR registry path
	key1, _, err := registry.CreateKey(registry.CURRENT_USER, gameBarRegistryPath, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to open GameDVR registry: %w", err))
	} else {
		defer key1.Close()

		// Backup and set GameDVR_Enabled
		originalValue, _, err := key1.GetIntegerValue("GameDVR_Enabled")
		if err != nil {
			originalValue = 1 // Default enabled
		}
		km.originalGameDVREnabled = uint32(originalValue)

		err = key1.SetDWordValue("GameDVR_Enabled", 0)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to set GameDVR_Enabled: %w", err))
		}

		// Backup and set AppCaptureEnabled
		originalValue, _, err = key1.GetIntegerValue("AppCaptureEnabled")
		if err != nil {
			originalValue = 1 // Default enabled
		}
		km.originalAppCaptureEnabled = uint32(originalValue)

		err = key1.SetDWordValue("AppCaptureEnabled", 0)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to set AppCaptureEnabled: %w", err))
		}
	}

	// 2. Disable GameDVR in GameConfigStore
	key2, _, err := registry.CreateKey(registry.CURRENT_USER, gameConfigStorePath, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to open GameConfigStore registry: %w", err))
	} else {
		defer key2.Close()

		// Backup and set GameDVR_Enabled
		originalValue, _, err := key2.GetIntegerValue("GameDVR_Enabled")
		if err != nil {
			originalValue = 1 // Default enabled
		}
		km.originalGameConfigEnabled = uint32(originalValue)

		err = key2.SetDWordValue("GameDVR_Enabled", 0)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to set GameConfigStore GameDVR_Enabled: %w", err))
		}
	}

	if len(errors) == 0 {
		km.gameBarRegistryModified = true
		logger.Debug("Successfully disabled GameBar via registry")
		return nil
	}

	// Log all errors but continue - hook will provide fallback
	for _, err := range errors {
		logger.Warn("GameBar registry error: %v", err)
	}
	return fmt.Errorf("GameBar disable had %d errors", len(errors))
}

func (km *KeyboardManager) RestoreGameBar() error {
	if !km.gameBarRegistryModified {
		return nil
	}

	km.registryMutex.Lock()
	defer km.registryMutex.Unlock()

	var errors []error

	// 1. Restore GameDVR registry values
	key1, err := registry.OpenKey(registry.CURRENT_USER, gameBarRegistryPath, registry.SET_VALUE)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to open GameDVR registry for restore: %w", err))
	} else {
		defer key1.Close()

		if km.originalGameDVREnabled == 0 {
			key1.DeleteValue("GameDVR_Enabled")
		} else {
			key1.SetDWordValue("GameDVR_Enabled", km.originalGameDVREnabled)
		}

		if km.originalAppCaptureEnabled == 0 {
			key1.DeleteValue("AppCaptureEnabled")
		} else {
			key1.SetDWordValue("AppCaptureEnabled", km.originalAppCaptureEnabled)
		}
	}

	// 2. Restore GameConfigStore registry values
	key2, err := registry.OpenKey(registry.CURRENT_USER, gameConfigStorePath, registry.SET_VALUE)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to open GameConfigStore registry for restore: %w", err))
	} else {
		defer key2.Close()

		if km.originalGameConfigEnabled == 0 {
			key2.DeleteValue("GameDVR_Enabled")
		} else {
			key2.SetDWordValue("GameDVR_Enabled", km.originalGameConfigEnabled)
		}
	}

	if len(errors) > 0 {
		for _, err := range errors {
			logger.Error("GameBar restore error: %v", err)
		}
		return fmt.Errorf("GameBar restore had %d errors", len(errors))
	}

	logger.Debug("Successfully restored GameBar registry settings")
	return nil
}

func (km *KeyboardManager) RegisterBlockingHotkeys() error {
	km.hotkeyMutex.Lock()
	defer km.hotkeyMutex.Unlock()

	user32 := syscall.NewLazyDLL("user32.dll")
	procRegisterHotKey := user32.NewProc("RegisterHotKey")
	procGetForegroundWindow := user32.NewProc("GetForegroundWindow")

	// Get current window handle - this might help RegisterHotKey work better
	hwnd, _, _ := procGetForegroundWindow.Call()

	// Register Win+Alt+G (GameBar record)
	ret, _, err := procRegisterHotKey.Call(hwnd, HOTKEY_WIN_ALT_G, MOD_WIN|MOD_ALT, 0x47) // G key
	if ret != 0 {
		km.registeredHotkeys = append(km.registeredHotkeys, HOTKEY_WIN_ALT_G)
		logger.Debug("Registered Win+Alt+G hotkey with HWND")
	} else {
		logger.Warn("Failed to register Win+Alt+G hotkey: %v", err)
		// Try without HWND as fallback
		ret, _, _ = procRegisterHotKey.Call(0, HOTKEY_WIN_ALT_G, MOD_WIN|MOD_ALT, 0x47)
		if ret != 0 {
			km.registeredHotkeys = append(km.registeredHotkeys, HOTKEY_WIN_ALT_G)
			logger.Debug("Registered Win+Alt+G hotkey without HWND")
		}
	}

	// Register Win+Alt+W (GameBar widget)
	ret, _, err = procRegisterHotKey.Call(hwnd, HOTKEY_WIN_ALT_W, MOD_WIN|MOD_ALT, 0x57) // W key
	if ret != 0 {
		km.registeredHotkeys = append(km.registeredHotkeys, HOTKEY_WIN_ALT_W)
		logger.Debug("Registered Win+Alt+W hotkey with HWND")
	} else {
		logger.Warn("Failed to register Win+Alt+W hotkey: %v", err)
		ret, _, _ = procRegisterHotKey.Call(0, HOTKEY_WIN_ALT_W, MOD_WIN|MOD_ALT, 0x57)
		if ret != 0 {
			km.registeredHotkeys = append(km.registeredHotkeys, HOTKEY_WIN_ALT_W)
			logger.Debug("Registered Win+Alt+W hotkey without HWND")
		}
	}

	// Register Win+Alt+B (HDR toggle)
	ret, _, err = procRegisterHotKey.Call(hwnd, HOTKEY_WIN_ALT_B, MOD_WIN|MOD_ALT, 0x42) // B key
	if ret != 0 {
		km.registeredHotkeys = append(km.registeredHotkeys, HOTKEY_WIN_ALT_B)
		logger.Debug("Registered Win+Alt+B hotkey with HWND")
	} else {
		logger.Warn("Failed to register Win+Alt+B hotkey: %v", err)
		ret, _, _ = procRegisterHotKey.Call(0, HOTKEY_WIN_ALT_B, MOD_WIN|MOD_ALT, 0x42)
		if ret != 0 {
			km.registeredHotkeys = append(km.registeredHotkeys, HOTKEY_WIN_ALT_B)
			logger.Debug("Registered Win+Alt+B hotkey without HWND")
		}
	}

	if len(km.registeredHotkeys) > 0 {
		logger.Debug("Successfully registered %d blocking hotkeys", len(km.registeredHotkeys))
	}

	return nil
}

func (km *KeyboardManager) UnregisterBlockingHotkeys() error {
	km.hotkeyMutex.Lock()
	defer km.hotkeyMutex.Unlock()

	if len(km.registeredHotkeys) == 0 {
		return nil
	}

	user32 := syscall.NewLazyDLL("user32.dll")
	procUnregisterHotKey := user32.NewProc("UnregisterHotKey")

	for _, hotkeyID := range km.registeredHotkeys {
		ret, _, err := procUnregisterHotKey.Call(0, uintptr(hotkeyID))
		if ret != 0 {
			logger.Debug("Unregistered hotkey ID %d", hotkeyID)
		} else {
			logger.Warn("Failed to unregister hotkey ID %d: %v", hotkeyID, err)
		}
	}

	km.registeredHotkeys = km.registeredHotkeys[:0]
	logger.Debug("Unregistered all blocking hotkeys")
	return nil
}

func isGameActive() bool {
	if globalKeyboard == nil || !globalKeyboard.isFocused {
		return false
	}
	currentState := globalKeyboard.GetCurrentState()
	return currentState == "state.play" || currentState == "state.play.pause"
}

func updateModifierStates(vkCode uint32, isKeyDown bool) {
	state := int32(0)
	if isKeyDown {
		state = 1
	}

	switch vkCode {
	case 0x5B, 0x5C: // Left/Right Windows keys
		atomic.StoreInt32(&winKeyDown, state)
	case 0xA4, 0xA5: // Left/Right Alt keys
		atomic.StoreInt32(&altKeyDown, state)
	case 0xA2, 0xA3: // Left/Right Ctrl keys
		atomic.StoreInt32(&ctrlKeyDown, state)
	}
}

func isWinAltComboPressed() bool {
	return atomic.LoadInt32(&winKeyDown) == 1 && atomic.LoadInt32(&altKeyDown) == 1
}

func isDangerousWinAltCombo(vkCode uint32) bool {
	if !isWinAltComboPressed() {
		return false
	}

	switch vkCode {
	case 0x47: // G key - GameBar record
		return true
	case 0x57: // W key - GameBar widget
		return true
	case 0x42: // B key - HDR toggle
		return true
	}
	return false
}

func blockInputTemporarily() {
	if atomic.CompareAndSwapInt32(&inputBlocked, 0, 1) {
		user32 := syscall.NewLazyDLL("user32.dll")
		procBlockInput := user32.NewProc("BlockInput")

		// Block input for a very short time
		procBlockInput.Call(1) // TRUE = block
		logger.Debug("Temporarily blocked all input")

		// Unblock after a short delay in a goroutine
		go func() {
			// Wait a tiny bit to let the shortcut attempt fail
			time.Sleep(100 * time.Millisecond)
			procBlockInput.Call(0) // FALSE = unblock
			atomic.StoreInt32(&inputBlocked, 0)
			logger.Debug("Unblocked input")
		}()
	}
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
	// 1. Set up registry modifications
	if err := km.DisableWinL(); err != nil {
		// Log warning but continue - not critical
		logger.Warn("Could not disable Win+L: %v", err)
	}

	if err := km.DisableGameBar(); err != nil {
		// Log warning but continue - not critical
		logger.Warn("Could not disable GameBar: %v", err)
	}

	if err := km.RegisterBlockingHotkeys(); err != nil {
		// Log warning but continue - not critical
		logger.Warn("Could not register blocking hotkeys: %v", err)
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

			// Consume WM_HOTKEY messages to prevent Windows from processing them
			if msg.Message == WM_HOTKEY {
				switch msg.WParam {
				case HOTKEY_WIN_ALT_G, HOTKEY_WIN_ALT_W, HOTKEY_WIN_ALT_B:
					logger.Debug("Consumed blocked hotkey message: %d", msg.WParam)
					continue // Don't dispatch this message
				}
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

func (km *KeyboardManager) sendWindowsKeyUpEvents() {
	// Try multiple approaches to clear Windows key state
	user32 := syscall.NewLazyDLL("user32.dll")
	
	// Method 1: Use keybd_event (older but sometimes more reliable)
	procKeybdEvent := user32.NewProc("keybd_event")
	
	const KEYEVENTF_KEYUP = 0x0002
	
	// Send key up events using keybd_event
	procKeybdEvent.Call(0x5B, 0, KEYEVENTF_KEYUP, 0) // Left Windows key up
	procKeybdEvent.Call(0x5C, 0, KEYEVENTF_KEYUP, 0) // Right Windows key up
	
	logger.Debug("Sent Windows key up events using keybd_event")
	
	// Method 2: Also try SendInput as backup
	procSendInput := user32.NewProc("SendInput")

	// INPUT structure for keyboard input
	type INPUT struct {
		Type uint32
		Ki   struct {
			VkCode      uint16
			ScanCode    uint16
			Flags       uint32
			Time        uint32
			ExtraInfo   uintptr
		}
	}

	const INPUT_KEYBOARD = 1

	// Create key up events for both left and right Windows keys with proper scan codes
	inputs := [2]INPUT{
		{
			Type: INPUT_KEYBOARD,
			Ki: struct {
				VkCode      uint16
				ScanCode    uint16
				Flags       uint32
				Time        uint32
				ExtraInfo   uintptr
			}{
				VkCode:   0x5B, // Left Windows key
				ScanCode: 0x5B, // Left Windows scan code
				Flags:    KEYEVENTF_KEYUP,
			},
		},
		{
			Type: INPUT_KEYBOARD,
			Ki: struct {
				VkCode      uint16
				ScanCode    uint16
				Flags       uint32
				Time        uint32
				ExtraInfo   uintptr
			}{
				VkCode:   0x5C, // Right Windows key
				ScanCode: 0x5C, // Right Windows scan code
				Flags:    KEYEVENTF_KEYUP,
			},
		},
	}

	// Send the key up events
	ret, _, _ := procSendInput.Call(
		2, // number of inputs
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]))

	if ret == 2 {
		logger.Debug("Sent Windows key up events using SendInput")
	} else {
		logger.Warn("SendInput failed, but keybd_event was attempted")
	}
}

func (km *KeyboardManager) Cleanup() {
	km.cleanupOnce.Do(func() {
		// Initialize DLLs locally to avoid race conditions
		user32 := syscall.NewLazyDLL("user32.dll")
		procUnhookWindowsHookEx := user32.NewProc("UnhookWindowsHookEx")
		procPostThreadMessage := user32.NewProc("PostThreadMessageW")

		// 1. Send synthetic Windows key up events to clear any stuck key state BEFORE removing hook
		km.sendWindowsKeyUpEvents()

		// 1.5. Clear internal modifier states to ensure no stuck key tracking
		atomic.StoreInt32(&winKeyDown, 0)
		atomic.StoreInt32(&altKeyDown, 0)
		atomic.StoreInt32(&ctrlKeyDown, 0)

		// 2. Remove hook
		if km.hook != 0 {
			procUnhookWindowsHookEx.Call(uintptr(km.hook))
			logger.Debug("Hook removed")
		}

		// 3. Restore registry
		if err := km.RestoreWinL(); err != nil {
			logger.Error("Failed to restore Win+L: %v", err)
		}

		if err := km.RestoreGameBar(); err != nil {
			logger.Error("Failed to restore GameBar: %v", err)
		}

		if err := km.UnregisterBlockingHotkeys(); err != nil {
			logger.Error("Failed to unregister hotkeys: %v", err)
		}

		// 4. Post quit message to stop message loop
		if km.threadID != 0 {
			procPostThreadMessage.Call(km.threadID, WM_QUIT, 0, 0)
		}

		// 5. Clear global state
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
	isKeyDown := wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN

	// Always update modifier states for combo detection
	updateModifierStates(kb.VkCode, isKeyDown)

	// Allow Windows keys to pass through to PowerToys and other system handlers
	// Do this BEFORE tracking them to avoid state inconsistencies
	if kb.VkCode == 0x5B || kb.VkCode == 0x5C { // Left/Right Windows keys
		// Let PowerToys and other system handlers process Windows keys
		ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
		return ret
	}

	// Always track keys synchronously to prevent stuck keys
	if globalKeyboard != nil {
		trackKeyForInputSystem(kb.VkCode, wParam)
	}

	// Check if we should be blocking keys
	if !isGameActive() {
		// Not in game - allow keys through
		ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
		return ret
	}

	// Detect dangerous Win+Alt combinations and use BlockInput to prevent them
	isWinPressed := atomic.LoadInt32(&winKeyDown) == 1
	isAltPressed := atomic.LoadInt32(&altKeyDown) == 1

	// If we detect a dangerous Win+Alt combination, block all input temporarily
	if isWinPressed && isAltPressed && isKeyDown {
		switch kb.VkCode {
		case 0x47, 0x57, 0x42: // G, W, B keys
			logger.Debug("Detected dangerous Win+Alt+%c - blocking all input", kb.VkCode)
			blockInputTemporarily()
			return 1 // Also block this specific key
		}
	}

	// Also specifically block the dangerous combinations even if our state tracking is off
	if isDangerousWinAltCombo(kb.VkCode) {
		logger.Debug("Blocked dangerous Win+Alt+%c combination", kb.VkCode)
		blockInputTemporarily() // Use BlockInput as additional protection
		return 1                // Block the key
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

	// Map Windows keys to their associated shift keys
	if vkCode == 0x5B { // Left Windows key
		key = ebiten.KeyShiftLeft
	} else if vkCode == 0x5C { // Right Windows key
		key = ebiten.KeyShiftRight
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
