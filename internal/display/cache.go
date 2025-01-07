package display

import "github.com/hajimehoshi/ebiten/v2"

// Little shared cache for static images that
// are created at runtime and only depend on rendering res

type Cache interface {
	Rebuild()
}

var caches = []Cache{}

func AttachCache(c Cache) {
	caches = append(caches, c)
}

func RebuildCaches() {
	for _, c := range caches {
		c.Rebuild()
	}
}

func ResetCaches() {
	caches = []Cache{}
}

var imageCache = make(map[string]*ebiten.Image)

func GetCachedImage(path string) (*ebiten.Image, bool) {
	img, ok := imageCache[path]
	return img, ok
}

func SetCachedImage(path string, img *ebiten.Image) {
	imageCache[path] = img
}

func clearImageCache() {
	for _, img := range imageCache {
		img.Deallocate()
	}
	imageCache = make(map[string]*ebiten.Image)
}
