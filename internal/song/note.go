package song

import (
	"math"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
)

type MarkerType int

const (
	MarkerTypeNone MarkerType = iota
	MarkerTypeBeat
	MarkerTypeMeasure
)

type Note struct {
	Target        int64 // ms from start of song it should be played
	TargetRelease int64 // ms the note should be held until

	HitTime     int64 // ms the note was hit
	ReleaseTime int64 // ms the note was released

	Progress   float64    // the note's progress towards the target down the track
	MarkerType MarkerType // Allows for special markers to be in the track
}

func NewNote(target, targetRelease int64) *Note {
	return &Note{
		Target:        target,
		TargetRelease: targetRelease,
	}
}

func (n *Note) Reset() {
	n.HitTime = 0
	n.ReleaseTime = 0
	n.Progress = 0
}

func (n *Note) WasHit() bool {
	return n.HitTime > 0
}

func (n *Note) Hit(currentTime int64) {
	if n.MarkerType != MarkerTypeNone {
		return
	}
	if !n.WasHit() {
		n.HitTime = currentTime
	}
}

func (n *Note) Release(currentTime int64) {
	if n.MarkerType != MarkerTypeNone {
		return
	}
	if n.WasHit() && n.ReleaseTime == 0 {
		n.ReleaseTime = currentTime
	}
}

// Updates note's progress towards the target
// 0 = not started, 1 = at target
func (n *Note) Update(currentTime int64) {
	n.Progress = math.Max(0, 1-(float64(n.Target-currentTime)/(float64(config.TRAVEL_TIME)/config.NOTE_SPEED)))
}
