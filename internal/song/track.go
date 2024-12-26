package song

import (
	"math"
	"sort"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/score/hit"
)

type TrackName string

// tracks from left to right
const (
	LeftBottom  TrackName = "track.left.bottom"
	LeftTop     TrackName = "track.left.top"
	Tap1        TrackName = "track.tap1"
	Tap2        TrackName = "track.tap2"
	Tap3        TrackName = "track.tap3"
	Space       TrackName = "track.space"
	Tap4        TrackName = "track.tap4"
	Tap5        TrackName = "track.tap5"
	Tap6        TrackName = "track.tap6"
	RightTop    TrackName = "track.right.top"
	RightBottom TrackName = "track.right.bottom"

	// Hmm...
	ExtraTop    TrackName = "track.extra.top"    // nav keys
	ExtraBottom TrackName = "track.extra.bottom" // arrow keys
)

type Track struct {
	Name          TrackName
	Notes         []Note
	notesInPlay   []*Note
	activeStart   int64 // when the track was activated for a hit
	nextNoteIndex int
}

func NewTrack(name TrackName, notes []Note) *Track {
	return &Track{
		Name: name,
	}
}
func (t *Track) Reset() {
	t.notesInPlay = nil
	t.activeStart = 0
	t.nextNoteIndex = 0
	for _, n := range t.Notes {
		n.Reset()
	}
}

func (t *Track) StartHit(currentTime int64) {
	// If the track is already active, ignore
	if t.activeStart > 0 || currentTime < 0 {
		return
	}
	t.activeStart = currentTime + config.INPUT_OFFSET
}

func (t *Track) EndHit(currentTime int64) {
	if t.activeStart > 0 {
		t.activeStart = 0
	}
}

// Allow notes to travel through the judgement line
const MAX_PROGRESS = 1.1

func (t *Track) Update(currentTime int64) *hit.HitRating {
	var notes []*Note
	var rating *hit.HitRating

	for _, n := range t.notesInPlay {
		n.Update(currentTime)

		// If the note is not hittable yet, continue
		if n.HitTime < currentTime-hit.Window.Bad {
			continue
		}

		// If the note has traveled past the judgement line
		if n.Progress > MAX_PROGRESS {
			// If the note has not been hit, mark it as missed
			if !n.IsHit() {
				n.Hit(currentTime)
			}
			continue
		}

		// If the note is beyond the hit window, continue
		if n.HitTime > currentTime+hit.Window.Bad {
			continue
		}

		// If active, check if hit is within window
		// only allow one hit per track at a time
		if t.activeStart > 0 && rating != nil {
			timeDiff := int64(math.Abs(float64(n.Target - t.activeStart)))
			hitRating := hit.GetHitRating(timeDiff)
			if hitRating != hit.Rating.Miss {
				n.Hit(currentTime)
				rating = &hitRating
				break
			}
		}

	}

	// Add notes to play when their travel time starts
	for t.nextNoteIndex < len(t.Notes) {
		note := t.Notes[t.nextNoteIndex]
		if currentTime >= (note.Target - config.GetTravelTime()) {
			notes = append(notes, &note)
			t.nextNoteIndex++
		} else {
			// Since notes are ordered by start time, no need to check further
			break
		}
	}

	t.notesInPlay = notes

	// Sort notes by Z depth (back to front)
	sort.Slice(t.notesInPlay, func(i, j int) bool {
		return t.notesInPlay[i].Progress > t.notesInPlay[j].Progress
	})
	return rating
}
