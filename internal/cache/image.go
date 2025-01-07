package cache

import "github.com/hajimehoshi/ebiten/v2"

type ImageCache struct {
	renderWidth  int
	renderHeight int
	cache        map[string]*ebiten.Image
}

var Image *ImageCache

func InitImageCache(renderWidth, renderHeight int) {
	Image = NewImageCache(renderWidth, renderHeight)
}

func NewImageCache(renderWidth, renderHeight int) *ImageCache {
	return &ImageCache{
		renderWidth:  renderWidth,
		renderHeight: renderHeight,
		cache:        make(map[string]*ebiten.Image),
	}
}

func (c *ImageCache) Get(path string) (*ebiten.Image, bool) {
	img, ok := c.cache[path]
	return img, ok
}

func (c *ImageCache) Set(path string, img *ebiten.Image) {
	c.cache[path] = img
}

func (c *ImageCache) Clear(renderWidth, renderHeight int) bool {
	// If same resolution or empty cache, no need to clear
	if c.renderWidth == renderWidth && c.renderHeight == renderHeight {
		return false
	} else if len(c.cache) == 0 {
		return false
	}

	if c.cache != nil {
		for _, img := range c.cache {
			img.Deallocate()
		}
	}
	c.cache = make(map[string]*ebiten.Image)
	return true
}
