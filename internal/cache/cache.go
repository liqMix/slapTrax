package cache

import "github.com/hajimehoshi/ebiten/v2"

// Little shared cache for static images that
// are created at runtime and only depend on rendering res
var imageCache = make(map[string]*ebiten.Image)
var audioCache = make(map[string][]byte)

func GetAudio(path string) ([]byte, bool) {
	audio, ok := audioCache[path]
	return audio, ok
}

func SetAudio(path string, audio []byte) {
	audioCache[path] = audio
}

func ClearAudioCache() {
	audioCache = make(map[string][]byte)
}

func GetImage(path string) (*ebiten.Image, bool) {
	img, ok := imageCache[path]
	return img, ok
}

func SetImage(path string, img *ebiten.Image) {
	imageCache[path] = img
}

func ClearImageCache() {
	for _, img := range imageCache {
		img.Deallocate()
	}
	imageCache = make(map[string]*ebiten.Image)
}
