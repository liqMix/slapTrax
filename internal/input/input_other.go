//go:build !windows
// +build !windows

package input

func applyOSHook(k *keyboard) error {
	// noops
	return nil
}

func removeOSHook(k *keyboard) {
	// noops
}

func (k *keyboard) Cleanup() {
	// noops
}
