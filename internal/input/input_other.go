//go:build !windows
// +build !windows

package input

func applyOSHook (k *keyboard) {
	// noops
}
func removeOSHook(k *keyboard) {
	// noops
}