package play

import (
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
func (r *Play) addHitEffects() {
	for _, hit := range r.state.Score.HitRecords[r.hitRecordIdx:] {
		hitTime := hit.Note.HitTime
		now := r.state.CurrentTime()
		speed := r.state.GetTravelTime()

		progress := types.GetTrackProgress(hitTime, now, -speed/4)
		if progress > 0 {
			if !hit.Note.IsHoldNote() {
				op := hit.Note.Progress
				hit.Note.Progress = progress
				path := GetNotePath(hit.Note.TrackName, hit.Note, true)
				hit.Note.Progress = op
				if path != nil {
					r.vectorCollection.AddPath(path)
				}
			}
		} else {
			r.hitRecordIdx++
		}
	}
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
