package play

import (
	"math"
	"time"

	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

// Rebuild reconstructs the entire vector path cache
func RebuildVectorCache() {
	renderWidth, renderHeight := display.Window.RenderSize()
	resolution := getMaxResolution(renderWidth)

	cache.Path.SetResolution(resolution)
	logger.Info("Rebuilding vector cache at %dx%d, with %d resolution", renderWidth, renderHeight, resolution)
	startTime := time.Now()

	buildNoteCache(resolution)
	buildJudgementLineCache()

	logger.Info("Rebuilt vector cache in %s", time.Since(startTime))
}

// getMaxResolution calculates the maximum resolution needed based on track distances
func getMaxResolution(renderWidth int) int {
	maxDistance := 0
	centerX := float64(renderWidth) / 2

	for _, trackName := range types.TrackNames() {
		point := notePoints[trackName][0]
		endPixel := int(point.X * float64(renderWidth))
		distance := int(math.Abs(float64(endPixel) - centerX))
		maxDistance = max(maxDistance, distance)
	}
	return maxDistance
}

// buildNoteCache pre-calculates all possible note paths
func buildNoteCache(resolution int) {
	for _, trackName := range types.TrackNames() {
		notePts := notePoints[trackName]
		for i := 0; i <= resolution; i++ {
			progress := float32(i) / float32(resolution)

			for _, solo := range []bool{true, false} {
				for _, hitEffectTrail := range []bool{true, false} {
					key := &cache.PathCacheKey{
						TrackName:        int(trackName),
						Progress:         i,
						Solo:             solo,
						IsHitEffectTrail: hitEffectTrail,
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

					cache.Path.Set(key, CreateNotePathFromPoints(notePts, progress, opts))
				}
			}
		}
	}
}

// getOrCreatePath retrieves a cached path or creates a new one if needed
func getOrCreatePath(key *cache.PathCacheKey) *cache.CachedPath {
	path := cache.Path.Get(key)

	if path == nil {
		// If not in cache, create new path
		res := cache.Path.GetResolution()
		progress := float32(key.Progress) / float32(res)
		opts := &NotePathOpts{
			lineWidth:       getJudgementWidth(),
			isLarge:         !key.Solo,
			largeWidthRatio: noteComboRatio,
			color:           types.TrackName(key.TrackName).NoteColor(),
			alpha:           GetNoteFadeAlpha(progress),
			solo:            key.Solo,
		}

		path = CreateNotePathFromPoints(notePoints[key.TrackName], progress, opts)
		cache.Path.Set(key, path)
	}

	return path
}

// GetNotePath retrieves the cached path for a regular note
func GetNotePath(track types.TrackName, note *types.Note, hitEffect bool) *cache.CachedPath {
	if note == nil || (note.Progress < 0 || note.Progress > 1) {
		return nil
	}

	res := cache.Path.GetResolution()
	progress := SmoothProgress(note.Progress)
	quantizedProgress := int(progress * float32(res))
	key := &cache.PathCacheKey{
		TrackName:        int(track),
		Progress:         quantizedProgress,
		Solo:             note.Solo,
		IsHitEffectTrail: hitEffect,
	}

	return getOrCreatePath(key)
}

func buildJudgementLineCache() {
	for _, trackName := range types.TrackNames() {
		for _, pressed := range []bool{true, false} {
			key := &cache.PathCacheKey{
				TrackName: int(trackName),
				Pressed:   pressed,
			}
			cache.Path.Set(key, CreateJudgementPath(trackName, pressed))
		}
	}
}

// GetJudgementLinePath retrieves the cached path for a judgement line
func GetJudgementLinePath(track types.TrackName, pressed bool) *cache.CachedPath {
	key := &cache.PathCacheKey{
		TrackName: int(track),
		Pressed:   pressed,
	}
	return getOrCreatePath(key)
}
