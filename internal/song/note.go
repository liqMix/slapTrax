package song

import "github.com/liqmix/ebiten-holiday-2024/internal/config"

type Note struct {
	Target      int64 // ms from start of song it should be played
	Hold        int64 // ms the note should be held until
	HitTime     int64 // ms the note was hit
	ReleaseTime int64 // ms the note was released

	Progress float64 // the note's progress towards the target down the track
}

func NewNote(target, hold int64) *Note {
	return &Note{
		Target: target,
		Hold:   hold,
	}
}

func (n *Note) Reset() {
	n.HitTime = 0
	n.ReleaseTime = 0
	n.Progress = 0
}

func (n *Note) Hit(currentTime int64) {
	if !n.IsHit() {
		n.HitTime = currentTime
	}
}

func (n *Note) Release(currentTime int64) {
	if n.Hold > 0 {
		n.ReleaseTime = currentTime
	}
}

func (n *Note) IsHit() bool {
	return n.HitTime > 0
}

// Updates note's progress towards the target
func (n *Note) Update(currentTime int64) {
	travelTime := config.TRAVEL_TIME / int64(config.NOTE_SPEED)
	progress := float64(currentTime-n.Target) / float64(travelTime)
	if progress < 0 {
		progress = 0
	}
	n.Progress = float64(currentTime-n.Target) / float64(travelTime)
}
