package types

import (
	"sort"

	"github.com/liqmix/slaptrax/internal/input"
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
		}
	}

	if t.StaleActive || t.Active {
		if input.NotActioned(t.Name.Action()) {
			t.Active = false
			t.StaleActive = false
		}
	}

	hit := None

	// Reset active notes
	notes := make([]*Note, 0, len(t.ActiveNotes))

	// Only update notes that are currently visible
	for _, n := range t.ActiveNotes {
		n.Update(currentTime, travelTime)
		
		// Initialize hold intervals if not done yet
		if n.IsHoldNote() && len(n.HoldIntervals) == 0 {
			n.CalculateIntervals()
		}
		
		// Check hold progress for active hold notes
		if n.IsHoldNote() && (n.WasHit() || n.MissedInitial) {
			n.CheckHoldProgress(currentTime, n.IsActive)
		}

		if n.IsHoldNote() {
			// Check if initial hit window has passed for hold notes FIRST
			if !n.WasHit() && !n.MissedInitial {
				initialWindowEnd := n.Target + LatestWindow
				if currentTime >= initialWindowEnd {
					// Initial hit window missed, mark as such but keep for reactivation
					n.MissedInitial = true
					n.IsInactive = true
					n.Miss(score) // Score the missed initial hit
				}
			}
			
			// Handle hold note logic with reactivation support
			if t.Active {
				// Check for reactivation of inactive notes (including missed initial)
				if n.CanReactivate(currentTime) {
					n.Reactivate()
				}
				
				// Hit the note if it hasn't been hit yet and is in timing window
				if !n.WasHit() && n.Target - LatestWindow <= currentTime && currentTime <= n.Target + LatestWindow {
					if n.Hit(currentTime, score) {
						hit = n.HitRating
					}
				}
				
				// Set note as active if track is active and note is not currently inactive
				// OR if we can reactivate it (for missed initial notes)
				if !n.IsInactive || n.CanReactivate(currentTime) {
					n.IsActive = true
					if n.MissedInitial {
						n.IsInactive = false // Clear inactive state for reactivated missed notes
					}
				}
			} else {
				// Track not active - make hold note inactive
				if n.IsActive {
					n.Release(currentTime)
				}
				n.IsActive = false
				n.IsInactive = true
			}
			
			notes = append(notes, n)
			continue
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
		if n.IsHoldNote() {
			// Remove hold notes that are past their release window
			pastReleaseWindow := currentTime >= n.TargetRelease + LatestWindow
			
			if pastReleaseWindow {
				// Mark any remaining intervals as missed and remove the note
				if len(n.HoldIntervals) > 0 {
					for i := n.LastCheckedInterval; i < len(n.HoldIntervals); i++ {
						n.HoldIntervalsHit[i] = false
						AddHoldIntervalWithCombo(false, true) // Definitive miss - break combo
						n.LastCheckedInterval = i + 1
					}
				}
				continue // Remove the note completely
			} else {
				// Still within valid time - keep for potential reactivation
				notes = append(notes, n)
				continue
			}
		}
		
		// Regular note missed
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
