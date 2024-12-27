package song

import (
	"math"
	"sort"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/score/hit"
)

type TrackName string

const (
	LeftBottom  TrackName = "track.leftbottom"
	LeftTop     TrackName = "track.lefttop"
	Center      TrackName = "track.center"
	RightBottom TrackName = "track.rightbottom"
	RightTop    TrackName = "track.righttop"

	EdgeTop  TrackName = "track.edgetop"
	EdgeTap1 TrackName = "track.edgetap1"
	EdgeTap2 TrackName = "track.edgetap2"
	EdgeTap3 TrackName = "track.edgetap3"
)

func TrackNames() []TrackName {
	return []TrackName{
		LeftBottom,
		LeftTop,
		Center,
		RightBottom,
		RightTop,
		EdgeTop,
		EdgeTap1,
		EdgeTap2,
		EdgeTap3,
	}
}

type Track struct {
	Name          TrackName
	Notes         []*Note
	VisibleNotes  []*Note
	activeStart   int64 // when the track was activated for a hit
	activeEnd     int64 // when the track was deactivated for a hit
	isHeld        bool
	nextNoteIndex int
}

func NewTrack(name TrackName, notes []*Note) *Track {
	// Reset the notes
	for _, n := range notes {
		n.Reset()
	}

	// Sort the notes by target time
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Target < notes[j].Target
	})

	return &Track{
		Name:         name,
		Notes:        notes,
		VisibleNotes: make([]*Note, 0),
	}
}

func (t *Track) Reset() {
	t.VisibleNotes = make([]*Note, 0)
	t.activeStart = 0
	t.activeEnd = 0
	t.nextNoteIndex = 0
	t.isHeld = false
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
	t.activeEnd = 0
}

func (t *Track) EndHit(currentTime int64) {
	if t.activeStart > 0 {
		t.activeStart = 0
		t.activeEnd = currentTime + config.INPUT_OFFSET
	}
}

func (t *Track) Update(currentTime int64) *hit.HitRating {
	var notes []*Note
	var rating *hit.HitRating

	for _, n := range t.VisibleNotes {
		n.Update(currentTime)

		// If the note is not hittable yet, keep it, no interaction yet
		if n.Target < currentTime-hit.Window.Bad {
			notes = append(notes, n)
			continue
		}

		// If the note has traveled too far (and it's not a hold note)
		if n.Target > currentTime+hit.Window.Bad && n.TargetRelease == 0 {
			// If the note has not been hit, mark it as missed
			if !n.WasHit() {
				n.Hit(currentTime)
			}
			continue
		}

		// By this time we have a note (or held note) that is within the hit window
		if t.isHeld {
			// If we released or we're past the max hold, release the note
			if t.activeEnd > 0 || n.TargetRelease > currentTime+hit.Window.Bad {
				n.Release(currentTime)
				t.isHeld = false
				t.activeEnd = 0
				continue
			}

			// Keep holding
			notes = append(notes, n)
			continue
		}

		// If we have an active track
		if t.activeStart > 0 {
			timeDiff := int64(math.Abs(float64(n.Target - t.activeStart)))
			hitRating := hit.GetHitRating(timeDiff)
			n.Hit(currentTime)
			rating = &hitRating

			// If the note is a hold note, set the track as active and keep it in play
			if n.TargetRelease > 0 {
				t.isHeld = true
				notes = append(notes, n)
			}

			// Reset the active start time
			t.activeStart = 0
			continue
		}

		notes = append(notes, n)
	}

	// Add notes to play when their travel time starts
	for t.nextNoteIndex < len(t.Notes) {
		note := t.Notes[t.nextNoteIndex]
		if currentTime >= (note.Target - config.GetTravelTime()) {
			notes = append(notes, note)
			t.nextNoteIndex++
		} else {
			// Since notes are ordered by start time, no need to check further
			break
		}
	}

	t.VisibleNotes = notes
	return rating
}
