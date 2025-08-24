//go:build !windows

package system

import "errors"

// OpenAudioFileDialog is not implemented on non-Windows systems
func OpenAudioFileDialog() (string, error) {
	return "", errors.New("file dialogs not implemented on this platform")
}

// SaveJSONFileDialog is not implemented on non-Windows systems  
func SaveJSONFileDialog(defaultName string) (string, error) {
	return "", errors.New("file dialogs not implemented on this platform")
}

// OpenJSONFileDialog is not implemented on non-Windows systems
func OpenJSONFileDialog() (string, error) {
	return "", errors.New("file dialogs not implemented on this platform")
}