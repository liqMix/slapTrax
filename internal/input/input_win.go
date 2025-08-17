//go:build windows
// +build windows

package input

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/logger"
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
	0x10: ebiten.KeyShift,   // Left Shift
	0xA0: ebiten.KeyShift,   // Left Shift (specific)
	0xA1: ebiten.KeyShift,   // Right Shift
	0x11: ebiten.KeyControl, // Left Control
	0xA2: ebiten.KeyControl, // Left Control (specific)
	0xA3: ebiten.KeyControl, // Right Control
	0x12: ebiten.KeyAlt,     // Left Alt
	0xA4: ebiten.KeyAlt,     // Left Alt (specific)
	0xA5: ebiten.KeyAlt,     // Right Alt
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
	0xBA: ebiten.KeySemicolon,    // ;
	0xBB: ebiten.KeyEqual,        // =
	0xBC: ebiten.KeyComma,        // ,
	0xBD: ebiten.KeyMinus,        // -
	0xBE: ebiten.KeyPeriod,       // .
	0xBF: ebiten.KeySlash,        // /
	0xC0: ebiten.KeyGraveAccent,  // `
	0xDB: ebiten.KeyLeftBracket,  // [
	0xDC: ebiten.KeyBackslash,    // \
	0xDD: ebiten.KeyRightBracket, // ]
	0xDE: ebiten.KeyApostrophe,   // '

	// Windows/Meta keys
	0x5B: ebiten.KeyMeta, // Left Windows key
	0x5C: ebiten.KeyMeta, // Right Windows key
}

var (
	// Global reference to prevent GC issues
	globalKeyboard *keyboard
	hookCallback   uintptr
)

func applyOSHook(k *keyboard) error {
	logger.Debug("Initializing Windows keyboard hook")

	// Store global reference
	globalKeyboard = k

	// Lock thread for Windows message processing
	runtime.LockOSThread()

	user32 := syscall.NewLazyDLL("user32.dll")
	kernel32 := syscall.NewLazyDLL("kernel32.dll")

	procSetWindowsHookEx := user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx := user32.NewProc("CallNextHookEx")
	procGetModuleHandle := kernel32.NewProc("GetModuleHandleW")

	// Get module handle
	moduleHandle, _, _ := procGetModuleHandle.Call(0)

	next := func(nCode int, wParam, lParam uintptr) uintptr {
		ret, _, _ := procCallNextHookEx.Call(k.osHook, uintptr(nCode), wParam, lParam)
		return ret
	}

	// Create callback and store reference
	hookCallback = syscall.NewCallback(func(nCode int, wParam, lParam uintptr) uintptr {
		// Safety check
		if globalKeyboard == nil {
			return next(nCode, wParam, lParam)
		}

		logger.Debug("Keyboard hook called: nCode=%d, wParam=%d", nCode, wParam)

		if !k.isFocused || nCode < 0 {
			return next(nCode, wParam, lParam)
		}

		kb := (*struct {
			VkCode      uint32
			ScanCode    uint32
			Flags       uint32
			Time        uint32
			DwExtraInfo uintptr
		})(unsafe.Pointer(lParam))

		k.m.Lock()
		defer k.m.Unlock()

		switch wParam {
		case 256: // WM_KEYDOWN
			key := keyMap[kb.VkCode]
			logger.Debug("Key pressed: %d %s", kb.VkCode, key)
			k.justPressed = append(k.justPressed, key)
		case 257: // WM_KEYUP
			key := keyMap[kb.VkCode]
			logger.Debug("Key released: %d %s", kb.VkCode, key)
			k.justReleased = append(k.justReleased, key)
		default:
			return next(nCode, wParam, lParam)
		}

		// Block the key event
		return 1
	})

	// Install hook
	hook, _, err := procSetWindowsHookEx.Call(
		13, // WH_KEYBOARD_LL
		hookCallback,
		moduleHandle,
		0,
	)

	if hook == 0 {
		return err
	}

	k.osHook = hook

	// Setup cleanup handlers
	setupCleanupHandlers(k)

	// Start message pump in separate goroutine
	go messageLoop(k)

	return nil
}

func messageLoop(k *keyboard) {
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
	for k.osHook != 0 {
		ret, _, _ := procGetMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			0, 0, 0,
		)

		if ret == 0 || ret == uintptr(0xFFFFFFFF) {
			break
		}

		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func removeOSHook(k *keyboard) {
	k.cleanup.Do(func() {
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

			k.osHook = 0
			globalKeyboard = nil
			hookCallback = 0

			// Unlock OS thread
			runtime.UnlockOSThread()
		}
	})
}

func setupCleanupHandlers(k *keyboard) {
	// Handle SIGINT/SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Debug("Received termination signal, cleaning up hook...")
		removeOSHook(k)
		os.Exit(0)
	}()

	// Windows-specific console handler for Ctrl+C, Ctrl+Break, close button
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleCtrlHandler := kernel32.NewProc("SetConsoleCtrlHandler")

	consoleHandler := syscall.NewCallback(func(ctrlType uint32) uintptr {
		logger.Debug("Console control event: %d", ctrlType)
		removeOSHook(k)
		return 1 // Handled
	})

	procSetConsoleCtrlHandler.Call(consoleHandler, 1)
}

// Cleanup function to call from main game loop
func (k *keyboard) Cleanup() {
	removeOSHook(k)
}
