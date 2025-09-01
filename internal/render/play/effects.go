package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/render/shaders"
	"github.com/liqmix/slaptrax/internal/types"
)

const (
	PULSE_LINE_WIDTH = 5.0
	PULSE_AMPLITUDE  = 0.06
	PULSE_MAX_SCALE  = 5.0
	PULSE_MIN_SCALE  = 0.1
	PULSE_SPEED      = 20.0
)

func (r *Play) addTrackEffects(track *types.Track) {
	// for _, note := range track.ActiveNotes {
	// 	r.drawNoteEffect(note)
	// }

}
func (r *Play) addHitEffects(screen *ebiten.Image) {
	for _, hit := range r.state.Score.HitRecords[r.hitRecordIdx:] {
		hitTime := hit.Note.HitTime
		now := r.state.CurrentTime()
		speed := r.state.GetTravelTime()

		// Calculate progress from hit point backwards towards center at 8x speed
		// Normal notes use travelTime = speed
		// For hit effects to travel very fast: travelTime = speed/8  
		// With negative direction: -speed/8 (eighth travel time = 8x speed)
		// This makes them shoot down the lane much faster
		progress := types.GetTrackProgress(hitTime, now, -speed/8)
		if progress > 0 && progress <= 1 {
			if !hit.Note.IsHoldNote() {
				r.addHitEffectShader(screen, hit.Note, float32(progress))
			}
		} else {
			r.hitRecordIdx++
		}
	}
}

// addHitEffectShader renders a single hit effect using shaders
func (r *Play) addHitEffectShader(screen *ebiten.Image, note *types.Note, hitProgress float32) {
	if shaders.Renderer == nil {
		return
	}

	// Get track points and center point for this track
	trackPoints := notePoints[note.TrackName]
	if len(trackPoints) == 0 {
		return
	}

	// Calculate effect opacity - fade out as it approaches the center
	// hitProgress is already correct: starts at 1.0 (hit point) and decreases to 0.0 (center)
	// The negative travel time in GetTrackProgress already handles the backward direction
	effectOpacity := hitProgress // Fade out as it approaches center (hitProgress goes 1.0 → 0.0)

	// Render the hit effect using shaders
	// Pass hitProgress directly - it's already 1.0 at hit point, 0.0 at center
	shaders.Renderer.RenderHitEffect(
		screen,
		note.TrackName,
		note,
		trackPoints,
		&playCenterPoint,
		hitProgress,    // Use hitProgress directly for correct direction
		effectOpacity,
	)
}

// func (r *Play) drawNoteEffect(note *types.Note) {}

// progress := SmoothProgress(p)
// sin, cos := calculatePulseWave(float64(progress))

// // Create local copy of points for pulse modification
// points := notePoints[hit.Note.TrackName]

// opts := &NotePathOpts{
// 	lineWidth: PULSE_LINE_WIDTH,
// 	isLarge:   true,
// 	color:     hit.Note.TrackName.NoteColor(),
// 	solo:      true,
// }
// // Modulate outer points with sine wave
// for i := 0; i < len(pulsePoints); i += 2 {
// 	pulsePoints[i].X = points[i].X * sin
// 	pulsePoints[i].Y = points[i].Y * sin
// }

// // Modulate center point with cosine wave
// pulsePoints[1].X = points[1].X * cos
// pulsePoints[1].Y = points[1].Y * cos

// // Fade out opacity and scale as effect travels
// opts.alpha = uint8(100 * progress)
// opts.largeWidthRatio = float32(PULSE_MAX_SCALE - ((PULSE_MAX_SCALE - PULSE_MIN_SCALE) * progress))

// // Create and add path to vector collection
// // path := r.vectorCache.createNotePath(pulsePoints, float32(progress), opts)
// // r.vectorCollection.Add(path.vertices, path.indices)

// // Add the normal path as well
// path := CreateNotePathFromPoints(points, float32(progress), opts)
// if path != nil {
// 	r.vectorCollection.Add(path.vertices, path.indices)
// }
