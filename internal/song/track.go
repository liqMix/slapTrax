package song

import (
	"fmt"
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
	AllNotes      []*Note
	ActiveNotes   []*Note
	activeTs      int64 // when the track was last activated
	releaseTs     int64 // when the track was last released
	nextNoteIndex int

	newHits []hit.HitRating
}

const minTs = math.MinInt64

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
		Name:        name,
		AllNotes:    notes,
		ActiveNotes: make([]*Note, 0),
		newHits:     make([]hit.HitRating, 0),
		activeTs:    minTs,
		releaseTs:   minTs,
	}
}

func (t *Track) Reset() {
	t.ActiveNotes = make([]*Note, 0)
	t.activeTs = math.MinInt64
	t.releaseTs = math.MinInt64
	t.nextNoteIndex = 0
	for _, n := range t.AllNotes {
		n.Reset()
	}
}

func (t *Track) NewHits() []hit.HitRating {
	return t.newHits
}

// not equal to int minimum 64 bit
func (t *Track) IsActive() bool {
	return t.activeTs > minTs
}

func (t *Track) Activate(currentTime int64) {
	// If the track is already active or held
	if t.IsActive() {
		return
	}
	fmt.Println("Activating track ", t.Name, " at ", currentTime)
	t.activeTs = currentTime + config.INPUT_OFFSET
	t.releaseTs = minTs
}

func (t *Track) Release(currentTime int64) {
	if t.releaseTs > minTs {
		return
	}
	t.releaseTs = currentTime + config.INPUT_OFFSET
	t.activeTs = minTs
}

func (t *Track) Update(currentTime int64) {
	// Reset the new hits
	t.newHits = t.newHits[:0]

	// No notes to update
	if len(t.ActiveNotes) == 0 && t.nextNoteIndex >= len(t.AllNotes) {
		return
	}

	notes := make([]*Note, 0, len(t.ActiveNotes))

	// Only update notes that are currently visible
	for _, n := range t.ActiveNotes {
		n.Update(currentTime)

		// Note is still approaching the hit window
		windowStart := n.Target - hit.Window.Bad
		if currentTime < windowStart {
			notes = append(notes, n)
			continue
		}

		// Drop expired notes
		windowEnd := n.Target + hit.Window.Bad
		releaseWindowEnd := n.TargetRelease + hit.Window.Bad

		if currentTime > windowEnd {
			if !n.WasHit() {
				rating := n.Hit(currentTime)
				if rating != hit.Rating.None {
					t.newHits = append(t.newHits, rating)
				}
			}
			if !n.IsHoldNote() || currentTime > releaseWindowEnd {
				continue // Drop
			}
		}

		// Handle hold notes
		if n.IsHoldNote() && n.WasHit() {
			var hitRating hit.HitRating
			// If player released or the hold window expired while holding
			if t.releaseTs > minTs {
				hitRating = n.Release(t.releaseTs)
			} else if currentTime > releaseWindowEnd {
				hitRating = n.Release(currentTime)
			} else {
				notes = append(notes, n) // Keep held or missed-hit hold notes
			}
			if hitRating != hit.Rating.None {
				// TODO: determine handling release scores
				// t.newHits = append(t.newHits, hitRating)
			}
			continue
		}

		// Process new hits
		if !n.WasHit() {
			if t.IsActive() {
				hitRating := n.Hit(t.activeTs)
				t.activeTs = minTs
				if hitRating != hit.Rating.None {
					t.newHits = append(t.newHits, hitRating)
				}

				// Begin holding the note
				if n.IsHoldNote() {
					notes = append(notes, n)
				}
				continue
			}
		}
		notes = append(notes, n)
	}

	// Add new approaching notes with single bounds check
	spawnTime := currentTime + config.ActualTravelTimeInt64
	if t.nextNoteIndex < len(t.AllNotes) {
		for i := t.nextNoteIndex; i < len(t.AllNotes); i++ {
			note := t.AllNotes[i]
			if note.Target > spawnTime {
				break
			}
			notes = append(notes, note)
			if t.Name == RightTop {
				fmt.Println("Track ", t.Name, " note at ", note.Target)
			}
			t.nextNoteIndex = i + 1
		}
	}

	t.ActiveNotes = notes
}
