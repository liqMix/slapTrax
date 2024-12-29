package song

import (
	"fmt"
	"math"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/score/hit"
)

type MarkerType int

const (
	MarkerTypeNone MarkerType = iota
	// MarkerTypeBeat
	// MarkerTypeMeasure
)

type Note struct {
	Target        int64 // ms from start of song it should be played
	TargetRelease int64 // ms the note should be held until

	HitTime     int64 // ms the note was hit
	ReleaseTime int64 // ms the note was released

	Progress   float64    // the note's progress towards the target down the track
	MarkerType MarkerType // Allows for special markers to be in the track ?

	HitRating hit.HitRating // The rating of the hit
}

func NewNote(target, targetRelease int64) *Note {
	return &Note{
		Target:        target,
		TargetRelease: targetRelease,
	}
}
func NewMarker(target int64, markerType MarkerType) *Note {
	return &Note{
		Target:     target,
		MarkerType: markerType,
	}
}

// String overload for fmt
func (n *Note) String() string {
	return fmt.Sprintf("Note{Target: %d, TargetRelease: %d, HitTime: %d, ReleaseTime: %d, Progress: %f, MarkerType: %d, HitRating: %v}",
		n.Target, n.TargetRelease, n.HitTime, n.ReleaseTime, n.Progress, n.MarkerType, n.HitRating)
}

func (n *Note) Reset() {
	n.HitTime = 0
	n.ReleaseTime = 0
	n.Progress = 0
	n.HitRating = hit.Rating.None
}

func (n *Note) IsHoldNote() bool {
	return n.TargetRelease > 0
}

func (n *Note) WasHit() bool {
	return n.HitTime > 0
}

func (n *Note) Hit(currentTime int64) hit.HitRating {
	if n.MarkerType != MarkerTypeNone || n.WasHit() {
		return hit.Rating.None
	}

	n.HitTime = currentTime
	n.HitRating = hit.GetHitRating(n.Target - n.HitTime)
	if n.HitRating != hit.Rating.Miss {
		fmt.Println(n.HitRating, n.Target, n.HitTime, n.Target-n.HitTime)
	}
	return n.HitRating
}

func (n *Note) Release(currentTime int64) hit.HitRating {
	if n.MarkerType != MarkerTypeNone || !n.IsHoldNote() {
		return hit.Rating.None
	}

	if n.WasHit() && n.ReleaseTime == 0 {
		n.ReleaseTime = currentTime
		if n.IsHoldNote() {
			// force rating to miss if the note was released early
			// *use more generous window for release
			diff := n.TargetRelease - currentTime
			if diff > (hit.Window.Bad * 3) {
				n.HitRating = hit.Rating.Miss
			}
		}
	}
	return n.HitRating
}

// Updates note's progress towards the target
// 0 = not started, 1 = at target
func (n *Note) Update(currentTime int64) {
	n.Progress = math.Max(0, 1-(float64(n.Target-currentTime)/config.ActualTravelTimeFloat64))
}
