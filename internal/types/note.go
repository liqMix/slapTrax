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
	Target        int64 // ms from start of song it should be played
	TargetRelease int64 // ms the note should be held until

	HitTime     int64 // ms the note was hit
	ReleaseTime int64 // ms the note was released

	Progress        float64    // the note's progress towards the target down the track
	ReleaseProgress float64    // the note releases's progress
	MarkerType      MarkerType // Allows for special markers to be in the track ?

	HitRating HitRating // The rating of the hit

	Solo bool // If the note is paired with other notes
	
	// Hold note interval tracking
	HoldIntervals       []int64 // Timestamps for each interval check
	HoldIntervalsHit    []bool  // Track which intervals were hit
	LastCheckedInterval int     // Progression tracker
	MissedInitial       bool    // For dimmer rendering when missed
	IsInactive          bool    // For released/missed hold notes that can be reactivated
	IsActive            bool    // Whether the hold note is currently being held
}

var noteId int = 0

func NewNote(trackName TrackName, target, targetRelease int64) *Note {
	noteId++
	return &Note{
		Id:            noteId,
		TrackName:     trackName,
		Target:        target,
		TargetRelease: targetRelease,
		Solo:          true,
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
	n.LastCheckedInterval = 0
	n.MissedInitial = false
	n.IsInactive = false
	n.IsActive = false
	n.HoldIntervals = nil
	n.HoldIntervalsHit = nil
}

func (n *Note) SetSolo(solo bool) {
	n.Solo = solo
}

func (n *Note) IsHoldNote() bool {
	return !user.S().DisableHoldNotes && n.TargetRelease > 0
}

func (n *Note) WasHit() bool {
	return n.HitRating != None && n.HitRating != Slop
}

func (n *Note) WasReleased() bool {
	return n.ReleaseTime > 0
}

func (n *Note) CanReactivate(currentTime int64) bool {
	if !n.IsHoldNote() {
		return false
	}
	
	// Can reactivate if:
	// 1. Note is in inactive state OR missed initial hit
	// 2. Current time is still within the hold note's duration window  
	// 3. Haven't processed all intervals yet
	
	hasIntervals := len(n.HoldIntervals) > 0
	withinDurationWindow := currentTime <= n.TargetRelease + LatestWindow
	hasRemainingIntervals := !hasIntervals || n.LastCheckedInterval < len(n.HoldIntervals)
	
	// For notes that missed the initial hit, allow reactivation during the hold duration
	if n.MissedInitial && !n.IsActive {
		return withinDurationWindow && hasRemainingIntervals && currentTime >= n.Target
	}
	
	// For released notes, allow reactivation
	return n.IsInactive && withinDurationWindow && hasRemainingIntervals
}

func (n *Note) Reactivate() {
	if !n.IsHoldNote() {
		return
	}
	
	// Reset to active state
	n.IsActive = true
	n.IsInactive = false
	
	// Clear release time to allow continued holding
	n.ReleaseTime = 0
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
	// Only process miss once per note
	if n.HitRating != None {
		return
	}
	
	n.HitRating = Slop
	if n.IsHoldNote() {
		n.MissedInitial = true
	}
	score.AddMiss(n)
}

func (n *Note) Release(releaseTime int64) {
	if !n.IsHoldNote() {
		return
	}

	// Only set release time if not already released (allow multiple releases for reactivation)
	if n.ReleaseTime == 0 {
		n.ReleaseTime = releaseTime + user.S().InputOffset
	}
	
	// Process any intervals that were being held up to release time
	if len(n.HoldIntervals) > 0 {
		// Check any intervals that should have been processed by release time
		currentTime := releaseTime + user.S().InputOffset
		for i := n.LastCheckedInterval; i < len(n.HoldIntervals); i++ {
			if n.HoldIntervals[i] <= currentTime {
				// This interval was held during the active period
				n.HoldIntervalsHit[i] = true
				AddHoldInterval(true)
				n.LastCheckedInterval = i + 1
			} else {
				// Haven't reached this interval yet - stop processing
				break
			}
		}
	}
	
	// Note: Don't mark remaining intervals as missed here since reactivation is possible
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

// CalculateIntervals divides hold note into 1/16 intervals (rounded down)
func (n *Note) CalculateIntervals() {
	if !n.IsHoldNote() {
		return
	}
	
	duration := n.TargetRelease - n.Target
	
	// Calculate 1/16 note intervals - assume 120 BPM for simplicity
	// 1/16 note at 120 BPM = (60000 / 120) / 4 = 125ms
	sixteenthNoteMs := int64(125) // This could be made dynamic based on song BPM
	
	intervalCount := int(duration / sixteenthNoteMs)
	if intervalCount < 1 {
		intervalCount = 1 // Minimum 1 interval
	}
	
	n.HoldIntervals = make([]int64, intervalCount)
	n.HoldIntervalsHit = make([]bool, intervalCount)
	
	intervalDuration := duration / int64(intervalCount)
	
	for i := 0; i < intervalCount; i++ {
		n.HoldIntervals[i] = n.Target + (int64(i+1) * intervalDuration)
	}
	
	// Update score to know about intervals and calculate interval value
	if score != nil {
		score.holdIntervalValue = score.hitValue / intervalCount
	}
}

// CheckHoldProgress evaluates intervals during hold period
func (n *Note) CheckHoldProgress(currentTime int64, trackActive bool) {
	if !n.IsHoldNote() || len(n.HoldIntervals) == 0 {
		return
	}
	
	// Only check if note was hit or missed but can be reactivated
	if !n.WasHit() && !n.MissedInitial {
		return
	}
	
	// Check intervals that have passed since last check
	for i := n.LastCheckedInterval; i < len(n.HoldIntervals); i++ {
		intervalTime := n.HoldIntervals[i]
		
		if currentTime >= intervalTime {
			// This interval has passed - check if note is currently active
			// For reactivated notes, trackActive && n.IsActive determines the hit
			hit := trackActive && n.IsActive && (n.WasHit() || n.MissedInitial)
			n.HoldIntervalsHit[i] = hit
			
			// Break combo immediately for missed intervals
			// This preserves the original behavior where missed intervals break combo
			AddHoldInterval(hit)
			
			n.LastCheckedInterval = i + 1
		} else {
			break // Haven't reached this interval yet
		}
	}
}

// GetHoldOpacity returns rendering opacity based on active status
func (n *Note) GetHoldOpacity() float32 {
	if !n.IsHoldNote() {
		return 1.0
	}
	
	// Active notes have full opacity for wobble effect
	if n.IsActive {
		return 1.0
	}
	
	// Inactive notes (released or missed) are dimmed but still visible
	if n.IsInactive || n.MissedInitial {
		return 0.3
	}
	
	// Default opacity for unhit notes
	return 1.0
}
