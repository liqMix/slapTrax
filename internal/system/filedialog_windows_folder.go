//go:build windows
// +build windows

package system

import (
	"syscall"
	"unsafe"
)

// Folder selection dialog structures and constants
type browseInfo struct {
	hwndOwner      uintptr
	pidlRoot       uintptr
	pszDisplayName uintptr
	lpszTitle      *uint16
	ulFlags        uint32
	lpfn           uintptr
	lParam         uintptr
	iImage         int
}

const (
	BIF_RETURNONLYFSDIRS = 0x00000001
	BIF_NEWDIALOGSTYLE   = 0x00000040
	BIF_NONEWFOLDERBUTTON = 0x00000200
)

var (
	shell32           = syscall.NewLazyDLL("shell32.dll")
	ole32             = syscall.NewLazyDLL("ole32.dll")
	shBrowseForFolder = shell32.NewProc("SHBrowseForFolderW")
	shGetPathFromIDList = shell32.NewProc("SHGetPathFromIDListW")
	coTaskMemFree     = ole32.NewProc("CoTaskMemFree")
)

// SelectFolderDialog opens a Windows folder selection dialog
func SelectFolderDialog(title string) (string, error) {
	titlePtr := syscall.StringToUTF16Ptr(title)
	displayName := make([]uint16, 260)
	
	bi := browseInfo{
		hwndOwner:      0,
		pidlRoot:       0,
		pszDisplayName: uintptr(unsafe.Pointer(&displayName[0])),
		lpszTitle:      titlePtr,
		ulFlags:        BIF_RETURNONLYFSDIRS | BIF_NEWDIALOGSTYLE | BIF_NONEWFOLDERBUTTON,
		lpfn:           0,
		lParam:         0,
		iImage:         0,
	}
	
	pidl, _, _ := shBrowseForFolder.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl == 0 {
		return "", nil // User cancelled
	}
	defer coTaskMemFree.Call(pidl)
	
	pathBuffer := make([]uint16, 260)
	ret, _, _ := shGetPathFromIDList.Call(pidl, uintptr(unsafe.Pointer(&pathBuffer[0])))
	if ret == 0 {
		return "", syscall.GetLastError()
	}
	
	return syscall.UTF16ToString(pathBuffer), nil
}