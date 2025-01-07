package play

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

// CacheKey represents a unique identifier for cached vector paths
type CacheKey struct {
	track          types.TrackName
	progress       int
	pressed        bool
	solo           bool
	hitEffectTrail bool
}

// CachedPath stores the pre-calculated vector path data
type CachedPath struct {
	vertices []ebiten.Vertex
	indices  []uint16
}

// VectorCache manages pre-calculated vector paths for improved performance
type VectorCache struct {
	renderWidth  int
	renderHeight int
	renderScale  float64
	resolution   int // How many steps to cache (e.g., 100 for 1% increments)

	pathCache map[CacheKey]*CachedPath

	disabled bool
}

// NewVectorCache creates a new vector cache instance
func NewVectorCache() *VectorCache {
	cache := &VectorCache{
		disabled: false,
	}
	cache.Rebuild()
	return cache
}

// getMaxResolution calculates the maximum resolution needed based on track distances
func (c *VectorCache) getMaxResolution() int {
	maxDistance := 0
	centerX := float64(c.renderWidth) / 2

	for _, trackName := range types.TrackNames() {
		point := notePoints[trackName][0]
		endPixel := int(point.X * float64(c.renderWidth))
		distance := int(math.Abs(float64(endPixel) - centerX))
		maxDistance = max(maxDistance, distance)
	}
	return maxDistance
}

// Rebuild reconstructs the entire cache when display parameters change
func (c *VectorCache) Rebuild() {
	if c.disabled {
		return
	}
	c.pathCache = make(map[CacheKey]*CachedPath)
	c.renderWidth, c.renderHeight = display.Window.RenderSize()
	c.renderScale = display.Window.RenderScale()
	c.resolution = c.getMaxResolution()

	if c.disabled {
		return
	}

	logger.Info("Rebuilding vector cache at %dx%d", c.renderWidth, c.renderHeight)
	startTime := time.Now()

	c.buildNoteCache()
	c.buildJudgementLineCache()

	logger.Info("Rebuilt vector cache in %s", time.Since(startTime))
}

// buildNoteCache pre-calculates all possible note paths
func (c *VectorCache) buildNoteCache() {
	for _, trackName := range types.TrackNames() {
		notePts := notePoints[trackName]

		for i := 0; i <= c.resolution; i++ {
			progress := float32(i) / float32(c.resolution)

			for _, solo := range []bool{true, false} {
				for _, hitEffectTrail := range []bool{true, false} {
					key := CacheKey{
						track:          trackName,
						progress:       i,
						solo:           solo,
						hitEffectTrail: hitEffectTrail,
					}
					alpha := GetNoteFadeAlpha(progress)
					if hitEffectTrail {
						alpha = uint8(200 * progress)
					}
					opts := &NotePathOpts{
						lineWidth:       getNoteWidth(),
						isLarge:         !solo,
						largeWidthRatio: noteComboRatio,
						color:           trackName.NoteColor(),
						alpha:           alpha,
						solo:            solo,
					}

					c.pathCache[key] = CreateNotePathFromPoints(notePts, progress, opts)
				}
			}
		}
	}
}

// getOrCreatePath retrieves a cached path or creates a new one if needed
func (c *VectorCache) getOrCreatePath(key CacheKey) *CachedPath {
	path, exists := c.pathCache[key]

	if !exists {
		// If not in cache, create new path
		progress := float32(key.progress) / float32(c.resolution)
		opts := &NotePathOpts{
			lineWidth:       getJudgementWidth(),
			isLarge:         !key.solo,
			largeWidthRatio: noteComboRatio,
			color:           key.track.NoteColor(),
			alpha:           GetNoteFadeAlpha(progress),
			solo:            key.solo,
		}

		path = CreateNotePathFromPoints(notePoints[key.track], progress, opts)

		// Store in cache if not disabled
		if !c.disabled {
			c.pathCache[key] = path
		}
	}

	return path
}

// GetNotePath retrieves the cached path for a regular note
func (c *VectorCache) GetNotePath(track types.TrackName, note *types.Note, hitEffect bool) *CachedPath {
	if note == nil || (note.Progress < 0 || note.Progress > 1) {
		return nil
	}

	progress := SmoothProgress(note.Progress)
	quantizedProgress := int(progress * float32(c.resolution))
	key := CacheKey{
		track:          track,
		progress:       quantizedProgress,
		solo:           note.Solo,
		hitEffectTrail: hitEffect,
	}

	return c.getOrCreatePath(key)
}

func (c *VectorCache) buildJudgementLineCache() {
	for _, trackName := range types.TrackNames() {
		for _, pressed := range []bool{true, false} {
			key := CacheKey{
				track:   trackName,
				pressed: pressed,
			}
			c.pathCache[key] = CreateJudgementPath(trackName, pressed)
		}
	}
}

// GetJudgementLinePath retrieves the cached path for a judgement line
func (c *VectorCache) GetJudgementLinePath(track types.TrackName, pressed bool) *CachedPath {
	if c.disabled {
		return CreateJudgementPath(track, pressed)
	}

	key := CacheKey{
		track:   track,
		pressed: pressed,
	}
	return c.getOrCreatePath(key)
}
