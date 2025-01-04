package play

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

// Vector cache for precalculated note paths
type NoteCacheKey struct {
	track    types.TrackName
	progress int
	solo     bool
}

type JudgementCacheKey struct {
	track   types.TrackName
	pressed bool
}

type VectorCache struct {
	renderWidth  int
	renderHeight int

	resolution int // How many steps to cache (e.g., 100 for 1% increments)

	noteCache          map[NoteCacheKey]*CachedPath
	judgementLineCache map[JudgementCacheKey]*CachedPath
}

type CachedPath struct {
	vertices []ebiten.Vertex
	indices  []uint16
}

func NewVectorCache() *VectorCache {
	cache := &VectorCache{
		noteCache: make(map[NoteCacheKey]*CachedPath),
	}
	cache.Rebuild()
	return cache
}

func (c *VectorCache) getMaxResolution() int {
	// Calculate the maximum distance any note travels
	maxDistance := 0
	for _, trackName := range types.TrackNames() {
		point := notePoints[trackName][0]
		endPixel := int(point.X * float64(c.renderWidth))
		startPixel := int(float64(c.renderWidth) / 2)
		distance := int(math.Abs(float64(endPixel - startPixel)))
		if distance > maxDistance {
			maxDistance = distance
		}
	}
	return maxDistance
}

func (c *VectorCache) Rebuild() {
	c.renderWidth, c.renderHeight = types.Window.RenderSize()
	c.resolution = c.getMaxResolution()
	logger.Info("Rebuilding vector cache at %dx%d", c.renderWidth, c.renderHeight)
	startTime := time.Now()

	c.buildNoteCache()
	c.buildJudgementLineCache()

	logger.Info("Rebuilt vector cache in %s", time.Since(startTime))
}

func (c *VectorCache) buildNoteCache() {
	// Clear existing cache
	c.noteCache = make(map[NoteCacheKey]*CachedPath)

	// For each track
	for _, trackName := range types.TrackNames() {
		notePts := notePoints[trackName]

		// For each progress step
		for i := 0; i <= c.resolution; i++ {
			progress := float32(i) / float32(c.resolution)

			// Cache both solo and non-solo versions
			for _, solo := range []bool{true, false} {
				key := NoteCacheKey{
					track:    trackName,
					progress: i,
					solo:     solo,
				}

				// Create the path for this specific progress
				path := c.createNotePath(notePts, progress, !solo, 3.0, trackName.NoteColor())
				c.noteCache[key] = path
			}
		}
	}
}

func (c *VectorCache) buildJudgementLineCache() {
	// Clear existing cache
	c.judgementLineCache = make(map[JudgementCacheKey]*CachedPath)

	// For each track
	for _, trackName := range types.TrackNames() {
		// Cache both solo and non-solo versions
		for _, pressed := range []bool{true, false} {
			key := JudgementCacheKey{
				track:   trackName,
				pressed: pressed,
			}

			// Create the path for this specific progress
			path := c.createNotePath(judgementLinePoints[trackName], 1, pressed, 1.5, trackName.NoteColor())
			c.judgementLineCache[key] = path
		}
	}

}

func (c *VectorCache) createNotePath(pts []*ui.Point, progress float32, large bool, largeRatio float32, color color.RGBA) *CachedPath {
	notePath := vector.Path{}
	cX, cY := playCenterPoint.ToRender32()

	// Build the path similar to your original draw code
	startX, startY := pts[0].ToRender32()
	x, y := cX+(startX-cX)*progress, cY+(startY-cY)*progress
	notePath.MoveTo(x, y)

	for i := 1; i < len(pts); i++ {
		if pts[i] == nil {
			continue
		}
		x, y = pts[i].ToRender32()
		x, y = cX+(x-cX)*progress, cY+(y-cY)*progress
		notePath.LineTo(x, y)
	}

	width := noteWidth * progress
	if large {
		width *= largeRatio
	}

	alpha := GetFadeAlpha(progress)
	color.A = alpha

	// Pre-compute vertices and indices
	vertices, indices := notePath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:    width,
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	})

	// Apply color and alpha
	ui.ColorVertices(vertices, color)

	return &CachedPath{
		vertices: vertices,
		indices:  indices,
	}
}

func (c *VectorCache) GetNotePath(track types.TrackName, note *types.Note) *CachedPath {
	if note == nil || note.Progress < 0 || note.Progress > 1 {
		return nil
	}

	progress := SmoothProgress(note.Progress)
	// Quantize progress to nearest cached value
	quantizedProgress := int(progress * float32(c.resolution))
	key := NoteCacheKey{
		track:    track,
		progress: quantizedProgress,
		solo:     note.Solo,
	}

	return c.noteCache[key]
}

func (c *VectorCache) GetJudgementLinePath(track types.TrackName, pressed bool) *CachedPath {
	key := NoteCacheKey{
		track:    track,
		progress: c.resolution,
		solo:     !pressed,
	}
	return c.noteCache[key]
}
