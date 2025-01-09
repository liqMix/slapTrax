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
	cbs   []*func()
}

func InitPathCache() {
	Path = NewPathCache()
}

func NewPathCache() *PathCache {
	return &PathCache{
		cache: make(map[PathCacheKey]*CachedPath),
		cbs:   make([]*func(), 0),
	}
}

func (c *PathCache) AddCb(cb *func()) {
	c.cbs = append(c.cbs, cb)
}

func (c *PathCache) RemoveCbs() {
	c.cbs = make([]*func(), 0)
}

func (c *PathCache) RemoveCb(cb *func()) {
	for i, v := range c.cbs {
		if v == cb {
			c.cbs = append(c.cbs[:i], c.cbs[i+1:]...)
			return
		}
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
	if c.cbs != nil {
		for _, cb := range c.cbs {
			(*cb)()
		}
	}
}

// If same render size or empty cache, no need to clear
// Resolution is currently tied directly to render size,
// so we don't need to clear if resolution changes
func (c *PathCache) Clear(renderWidth, renderHeight int) {
	if len(c.cache) == 0 {
		return
	} else if c.renderWidth == renderWidth && c.renderHeight == renderHeight {
		return
	}
	c.cache = make(map[PathCacheKey]*CachedPath)
	if c.cbs != nil {
		for _, cb := range c.cbs {
			(*cb)()
		}
	}
	return
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
