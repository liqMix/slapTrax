package types

import (
	"sort"

	"github.com/liqmix/ebiten-holiday-2024/internal/input"
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

func (t *Track) CheckInput(currentTime int64) {
	if !t.Active && !t.StaleActive {
		if input.K.AreAny(TrackNameToKeys[t.Name], input.JustPressed) {
			t.Active = true
			return
		}
	}

	if t.StaleActive || t.Active {
		if !input.K.AreAny(TrackNameToKeys[t.Name], input.Held) {
			t.Active = false
			t.StaleActive = false
		}
	}
}
