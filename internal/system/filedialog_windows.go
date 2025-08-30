//go:build windows

package system

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	comdlg32          = syscall.NewLazyDLL("comdlg32.dll")
	getOpenFileNameW  = comdlg32.NewProc("GetOpenFileNameW")
	getSaveFileNameW  = comdlg32.NewProc("GetSaveFileNameW")
)

type openFileName struct {
	lStructSize       uint32
	hwndOwner         uintptr
	hInstance         uintptr
	lpstrFilter       *uint16
	lpstrCustomFilter *uint16
	nMaxCustFilter    uint32
	nFilterIndex      uint32
	lpstrFile         *uint16
	nMaxFile          uint32
	lpstrFileTitle    *uint16
	nMaxFileTitle     uint32
	lpstrInitialDir   *uint16
	lpstrTitle        *uint16
	flags             uint32
	nFileOffset       uint16
	nFileExtension    uint16
	lpstrDefExt       *uint16
	lCustData         uintptr
	lpfnHook          uintptr
	lpTemplateName    *uint16
}

const (
	OFN_ALLOWMULTISELECT     = 0x00000200
	OFN_CREATEPROMPT         = 0x00002000
	OFN_ENABLEHOOK           = 0x00000020
	OFN_ENABLETEMPLATE       = 0x00000040
	OFN_ENABLETEMPLATEHANDLE = 0x00000080
	OFN_EXPLORER             = 0x00080000
	OFN_EXTENSIONDIFFERENT   = 0x00000400
	OFN_FILEMUSTEXIST        = 0x00001000
	OFN_HIDEREADONLY         = 0x00000004
	OFN_LONGNAMES            = 0x00200000
	OFN_NOCHANGEDIR          = 0x00000008
	OFN_NODEREFERENCELINKS   = 0x00100000
	OFN_NOLONGNAMES          = 0x00040000
	OFN_NONETWORKBUTTON      = 0x00020000
	OFN_NOREADONLYRETURN     = 0x00008000
	OFN_NOTESTFILECREATE     = 0x00010000
	OFN_NOVALIDATE           = 0x00000100
	OFN_OVERWRITEPROMPT      = 0x00000002
	OFN_PATHMUSTEXIST        = 0x00000800
	OFN_READONLY             = 0x00000001
	OFN_SHAREAWARE           = 0x00004000
	OFN_SHOWHELP             = 0x00000010
)

// OpenAudioFileDialog opens a Windows file dialog for selecting audio files
func OpenAudioFileDialog() (string, error) {
	title := syscall.StringToUTF16Ptr("Select Audio File")
	// Manually create filter with NUL separators for Windows file dialog
	// We can't use syscall.StringToUTF16 because it doesn't handle embedded NULs
	filter := utf16.Encode([]rune("Audio Files\x00*.ogg;*.mp3;*.wav;*.flac\x00OGG Files\x00*.ogg\x00MP3 Files\x00*.mp3\x00WAV Files\x00*.wav\x00FLAC Files\x00*.flac\x00All Files\x00*.*\x00\x00"))
	
	filename := make([]uint16, 260)
	
	ofn := openFileName{
		lStructSize:     uint32(unsafe.Sizeof(openFileName{})),
		lpstrFilter:     &filter[0],
		nFilterIndex:    1,
		lpstrFile:       &filename[0],
		nMaxFile:        260,
		lpstrTitle:      title,
		flags:          OFN_FILEMUSTEXIST | OFN_PATHMUSTEXIST | OFN_NOCHANGEDIR | OFN_EXPLORER,
	}
	
	ret, _, _ := getOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	if ret == 0 {
		return "", syscall.GetLastError()
	}
	
	return syscall.UTF16ToString(filename), nil
}

// SaveJSONFileDialog opens a Windows file dialog for saving JSON files
func SaveJSONFileDialog(defaultName string) (string, error) {
	title := syscall.StringToUTF16Ptr("Save Chart")
	filter := utf16.Encode([]rune("JSON Files\x00*.json\x00All Files\x00*.*\x00\x00"))
	
	filename := make([]uint16, 260)
	if defaultName != "" {
		copy(filename, syscall.StringToUTF16(defaultName))
	}
	
	ofn := openFileName{
		lStructSize:     uint32(unsafe.Sizeof(openFileName{})),
		lpstrFilter:     &filter[0],
		nFilterIndex:    1,
		lpstrFile:       &filename[0],
		nMaxFile:        260,
		lpstrTitle:      title,
		lpstrDefExt:     syscall.StringToUTF16Ptr("json"),
		flags:          OFN_OVERWRITEPROMPT | OFN_NOCHANGEDIR | OFN_EXPLORER,
	}
	
	ret, _, _ := getSaveFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	if ret == 0 {
		return "", syscall.GetLastError()
	}
	
	return syscall.UTF16ToString(filename), nil
}

// OpenJSONFileDialog opens a Windows file dialog for selecting JSON chart files
func OpenJSONFileDialog() (string, error) {
	title := syscall.StringToUTF16Ptr("Open Chart")
	filter := utf16.Encode([]rune("JSON Files\x00*.json\x00All Files\x00*.*\x00\x00"))
	
	filename := make([]uint16, 260)
	
	ofn := openFileName{
		lStructSize:     uint32(unsafe.Sizeof(openFileName{})),
		lpstrFilter:     &filter[0],
		nFilterIndex:    1,
		lpstrFile:       &filename[0],
		nMaxFile:        260,
		lpstrTitle:      title,
		flags:          OFN_FILEMUSTEXIST | OFN_PATHMUSTEXIST | OFN_NOCHANGEDIR | OFN_EXPLORER,
	}
	
	ret, _, _ := getOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	if ret == 0 {
		return "", syscall.GetLastError()
	}
	
	return syscall.UTF16ToString(filename), nil
}