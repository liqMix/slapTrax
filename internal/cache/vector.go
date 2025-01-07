package cache

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var Path *PathCache

type CachedPath struct {
	Vertices []ebiten.Vertex
	Indices  []uint16
}

type PathCacheKey struct {
	TrackName        int // types.TrackName
	Progress         int
	Pressed          bool
	Solo             bool
	IsHitEffectTrail bool
}

type PathCache struct {
	renderWidth  int
	renderHeight int
	resolution   int // How many steps to cache (e.g., 100 for 1% increments)

	cache map[PathCacheKey]*CachedPath
}

func InitPathCache(renderWidth, renderHeight, resolution int) {
	Path = NewPathCache(renderWidth, renderHeight, resolution)
}

func NewPathCache(renderWidth, renderHeight, resolution int) *PathCache {
	return &PathCache{
		renderWidth:  renderWidth,
		renderHeight: renderHeight,
		resolution:   resolution,
		cache:        make(map[PathCacheKey]*CachedPath),
	}
}

func (c *PathCache) GetResolution() int {
	return c.resolution
}

func (c *PathCache) SetResolution(resolution int) {
	c.resolution = resolution
	c.Clear(c.renderWidth, c.renderHeight)
}

// for if we're changing something fundamental about the paths
// but render size is the same:
//   - note color
func (c *PathCache) ForceClear() {
	// ok fine fine, I'll clear it
	c.cache = make(map[PathCacheKey]*CachedPath)
}

// If same render size or empty cache, no need to clear
// Resolution is currently tied directly to render size,
// so we don't need to clear if resolution changes
func (c *PathCache) Clear(renderWidth, renderHeight int) bool {
	if len(c.cache) == 0 {
		return true
	} else if c.renderWidth == renderWidth && c.renderHeight == renderHeight {
		return false
	}
	c.cache = make(map[PathCacheKey]*CachedPath)
	return true
}

func (c *PathCache) Get(key *PathCacheKey) *CachedPath {
	if key == nil {
		return nil
	}
	return c.cache[*key]
}

func (c *PathCache) Set(key *PathCacheKey, path *CachedPath) {
	if key == nil || path == nil {
		return
	}
	c.cache[*key] = path
}
