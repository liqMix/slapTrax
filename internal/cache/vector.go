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

	isBuilding bool
	cache      map[PathCacheKey]*CachedPath
	cbs        []*func()
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
}

func (c *PathCache) Clear() {
	c.cache = make(map[PathCacheKey]*CachedPath)
	if c.cbs != nil {
		for _, cb := range c.cbs {
			(*cb)()
		}
	}
}

func (c *PathCache) Get(key *PathCacheKey) *CachedPath {
	if key == nil {
		return nil
	}

	// Can modify all notes here.

	return c.cache[*key]
}

func (c *PathCache) Set(key *PathCacheKey, path *CachedPath) {
	if key == nil || path == nil {
		return
	}
	c.cache[*key] = path
}

func (c *PathCache) SetIsBuilding(isBuilding bool) {
	c.isBuilding = isBuilding
}

func (c *PathCache) IsBuilding() bool {
	return c.isBuilding
}
