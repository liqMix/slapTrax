package cache

import "github.com/hajimehoshi/ebiten/v2"

type ImageCache struct {
	cache map[string]*ebiten.Image
}

var Image *ImageCache

func NewImageCache() *ImageCache {
	return &ImageCache{
		cache: make(map[string]*ebiten.Image),
	}
}

func (c *ImageCache) Get(path string) (*ebiten.Image, bool) {
	img, ok := c.cache[path]
	return img, ok
}

func (c *ImageCache) Set(path string, img *ebiten.Image) {
	c.cache[path] = img
}

func (c *ImageCache) Clear() {
	if c.cache != nil {
		for _, img := range c.cache {
			img.Deallocate()
		}
	}
	c.cache = make(map[string]*ebiten.Image)
}
