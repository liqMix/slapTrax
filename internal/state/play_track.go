package state

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

func (p *Play) updateTrack(t *types.Track, currentTime int64, score *types.Score) {
	t.CheckInput(currentTime)

	// Reset the new hits
	notes := make([]*types.Note, 0, len(t.ActiveNotes))

	// Only update notes that are currently visible
	for _, n := range t.ActiveNotes {
		n.Update(currentTime, p.GetTravelTime())

		if n.IsHoldNote() {
			if n.WasReleased() {
				continue
			}

			if !n.WasHit() && t.Active && !t.StaleActive {
				if n.Hit(currentTime, score) {
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
					t.StaleActive = true
					continue
				}
			}
		}

		// Note not yet reached the out of bounds window
		windowEnd := n.Target + types.LatestWindow
		if n.IsHoldNote() {
			windowEnd = n.TargetRelease + types.LatestWindow
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
			if note.Target > p.MaxTrackTime() {
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
			if n.InWindow(currentTime-types.EarliestWindow, currentTime+types.LatestWindow) {
				return
			}
		}
		t.StaleActive = true
	}
}
