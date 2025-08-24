package types

import (
	"sort"

	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/logger"
)

type Track struct {
	Name        TrackName
	AllNotes    []*Note
	ActiveNotes []*Note

	Active      bool
	StaleActive bool

	NextNoteIndex int
}

func NewTrack(name TrackName, notes []*Note, beatInterval int64) *Track {
	// Reset the notes
	for _, n := range notes {
		n.Reset()
	}

	// Sort the notes by target time
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Target < notes[j].Target
	})

	return &Track{
		Name:     name,
		AllNotes: notes,
	}
}

func (t *Track) HasMoreNotes() bool {
	return t.NextNoteIndex < len(t.AllNotes)
}

func (t *Track) Reset() {
	t.ActiveNotes = make([]*Note, 0)
	t.Active = false
	t.StaleActive = false
	t.NextNoteIndex = 0

	for _, n := range t.AllNotes {
		n.Reset()
	}
}

func (t Track) IsPressed() bool {
	return t.Active || t.StaleActive
}

func (t *Track) Update(currentTime int64, travelTime int64, maxTime int64) HitRating {
	if !t.Active && !t.StaleActive {
		if input.JustActioned(t.Name.Action()) {
			t.Active = true
			if t.Name == TrackLeftBottom {
				logger.Debug("Track LeftBottom activated")
			}
		}
	}

	if t.StaleActive || t.Active {
		if input.NotActioned(t.Name.Action()) {
			t.Active = false
			t.StaleActive = false
			if t.Name == TrackLeftBottom {
				logger.Debug("Track LeftBottom deactivated")
			}
		}
	}

	hit := None

	// Reset active notes
	notes := make([]*Note, 0, len(t.ActiveNotes))

	// Only update notes that are currently visible
	for _, n := range t.ActiveNotes {
		n.Update(currentTime, travelTime)

		if n.IsHoldNote() {
			if n.WasReleased() {
				continue
			}

			if !n.WasHit() && t.Active && !t.StaleActive {
				if n.Hit(currentTime, score) {
					hit = n.HitRating
					t.StaleActive = true
					notes = append(notes, n)
					continue
				}
			}

			if !t.Active && !t.StaleActive && !n.WasReleased() {
				n.Release(currentTime)
			}
		} else {
			if n.WasHit() {
				continue
			}

			if t.Active && !t.StaleActive {
				if n.Hit(currentTime, score) {
					hit = n.HitRating
					t.StaleActive = true
					continue
				}
			}
		}

		// Note not yet reached the out of bounds window
		windowEnd := n.Target + LatestWindow
		if n.IsHoldNote() {
			windowEnd = n.TargetRelease + LatestWindow
		}
		if currentTime < windowEnd {
			notes = append(notes, n)
			continue
		}

		// Drop expired notes
		if n.IsHoldNote() && t.Active {
			// u good
			n.Release(currentTime)
			continue
		}
		n.Miss(score)
	}

	// Add new approaching notes
	if t.NextNoteIndex < len(t.AllNotes) {
		for i := t.NextNoteIndex; i < len(t.AllNotes); i++ {
			note := t.AllNotes[i]
			if note.Target > maxTime {
				break
			}
			notes = append(notes, note)
			t.NextNoteIndex = i + 1
		}

	}

	t.ActiveNotes = notes

	// dont stay active if no notes in window
	if t.Active && !t.StaleActive {
		for _, n := range t.ActiveNotes {
			if n.InWindow(currentTime-EarliestWindow, currentTime+LatestWindow) {
				return hit
			}
		}
		t.StaleActive = true
	}
	return hit
}
