package types

import (
	"math"

	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/user"
)

type MarkerType int

const (
	MarkerTypeNone MarkerType = iota
	// MarkerTypeBeat
	// MarkerTypeMeasure
)

type Note struct {
	Id            int
	TrackName     TrackName
	Target        int64 // ms from start of song it should be played (calculated from TargetBeat)
	TargetRelease int64 // ms the note should be held until (calculated from TargetReleaseBeat)
	
	// Beat-based positioning (musical time, BPM-independent)
	TargetBeat        float64 // beat position from start of song
	TargetReleaseBeat float64 // beat position where note should be released

	HitTime     int64 // ms the note was hit
	ReleaseTime int64 // ms the note was released

	Progress        float64    // the note's progress towards the target down the track
	ReleaseProgress float64    // the note releases's progress
	MarkerType      MarkerType // Allows for special markers to be in the track ?

	HitRating HitRating // The rating of the hit

	Solo bool // If the note is paired with other notes
}

var noteId int = 0

func NewNote(trackName TrackName, target, targetRelease int64) *Note {
	noteId++
	return &Note{
		Id:            noteId,
		TrackName:     trackName,
		Target:        target,
		TargetRelease: targetRelease,
		// Beat positions will be calculated when BPM context is available
		TargetBeat:        0, // Will be set by editor when BPM is known
		TargetReleaseBeat: 0, // Will be set by editor when BPM is known
		Solo:              true,
	}
}

// NewNoteFromBeats creates a note using beat positions and calculates millisecond positions
func NewNoteFromBeats(trackName TrackName, targetBeat, targetReleaseBeat float64, bpm float64) *Note {
	noteId++
	
	// Convert beats to milliseconds with rounding for integer storage
	quarterNoteMs := 60000.0 / bpm
	target := int64(math.Round(targetBeat * quarterNoteMs))
	targetRelease := int64(0)
	if targetReleaseBeat > 0 {
		targetRelease = int64(math.Round(targetReleaseBeat * quarterNoteMs))
	}
	
	return &Note{
		Id:                noteId,
		TrackName:         trackName,
		Target:            target,
		TargetRelease:     targetRelease,
		TargetBeat:        targetBeat,
		TargetReleaseBeat: targetReleaseBeat,
		Solo:              true,
	}
}

func NewMarker(target int64, markerType MarkerType) *Note {
	return &Note{
		Target:     target,
		MarkerType: markerType,
	}
}

func (n *Note) Reset() {
	n.HitTime = 0
	n.ReleaseTime = 0
	n.Progress = 0
	n.HitRating = None
}

func (n *Note) SetSolo(solo bool) {
	n.Solo = solo
}

func (n *Note) IsHoldNote() bool {
	return !user.S().DisableHoldNotes && n.TargetRelease > 0
}

func (n *Note) WasHit() bool {
	return n.HitRating != None
}

func (n *Note) WasReleased() bool {
	return n.ReleaseTime > 0
}

func (n *Note) Hit(hitTime int64, score *Score) bool {
	if n.WasHit() {
		return false
	}

	diff := n.Target - hitTime + user.S().InputOffset
	timing := GetHitTiming(diff)
	rating := GetHitRating(diff)
	if rating == None {
		return false
	}

	n.HitTime = hitTime
	n.HitRating = rating
	logger.Debug("Id: %d | Hit: %s | Diff: %d | Target: %d | HitTime: %d", n.Id, n.HitRating, int(diff), n.Target, n.HitTime)
	AddHit(&HitRecord{
		Note:      n,
		HitDiff:   n.Target - hitTime,
		HitRating: n.HitRating,
		HitTiming: timing,
	})

	return true
}

func (n *Note) Miss(score *Score) {
	n.HitRating = Miss
	score.AddMiss(n)
}

func (n *Note) Release(releaseTime int64) {
	if !n.IsHoldNote() || !n.WasHit() || n.WasReleased() {
		return
	}

	if n.WasHit() {
		n.ReleaseTime = releaseTime + user.S().InputOffset
		if n.IsHoldNote() {
			// force rating to miss if the note was released early
			// *use more generous window for release
			diff := n.TargetRelease - n.ReleaseTime

			hit := &HitRecord{
				Note:        n,
				HitDiff:     n.Target - n.HitTime,
				ReleaseDiff: diff,
			}

			if diff > int64(Bad.Window(false)*3) {
				n.HitRating = Miss
				hit.HitRating = Miss
				hit.HitTiming = HitTimingEarly

			} else {
				n.HitRating = GetHitRating(hit.HitDiff)
				hit.HitRating = n.HitRating
				hit.HitTiming = GetHitTiming(hit.HitDiff)
			}

			AddHit(hit)
		}
	}
}

func (n *Note) InWindow(start, end int64) bool {
	if n.Target >= start && n.Target <= end {
		return true
	}
	if n.TargetRelease >= start && n.TargetRelease <= end {
		return true
	}
	return false
}

// Updates note's progress towards the target
// 0 = not started, 1 = at target
func (n *Note) Update(currentTime int64, travelTime int64) {
	n.Progress = GetTrackProgress(n.Target, currentTime, travelTime)
	if n.IsHoldNote() {
		n.ReleaseProgress = GetTrackProgress(n.TargetRelease, currentTime, travelTime)
	}
}

func GetTrackProgress(targetTime, currentTime, travelTime int64) float64 {
	return math.Max(0, 1-float64(targetTime-currentTime)/float64(travelTime))
}
