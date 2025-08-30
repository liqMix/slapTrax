package types

import (
	"fmt"
	"math"
)

// MusicPosition represents a position in musical time using measures, beats, and divisions
type MusicPosition struct {
	Measure  int     // 1-based measure number (1 = first measure)
	Beat     int     // 1-based beat number within the measure (1 = first beat)
	Division int     // Division within the beat (0-based, where max = timeDivision-1)
}

// NewMusicPosition creates a new MusicPosition
func NewMusicPosition(measure, beat, division int) MusicPosition {
	return MusicPosition{
		Measure:  measure,
		Beat:     beat,
		Division: division,
	}
}

// NewMusicPositionFromBeats converts from total beat count to measure-based position
func NewMusicPositionFromBeats(totalBeats float64, beatsPerMeasure, timeDivision int) MusicPosition {
	// Calculate measure (1-based)
	measureIndex := int(totalBeats) / beatsPerMeasure
	measure := measureIndex + 1
	
	// Calculate beat within measure (1-based)
	beatWithinMeasure := totalBeats - float64(measureIndex*beatsPerMeasure)
	beat := int(beatWithinMeasure) + 1
	
	// Calculate division within beat (0-based)
	beatFraction := beatWithinMeasure - float64(int(beatWithinMeasure))
	divisionsPerBeat := float64(timeDivision) / 4.0
	division := int(math.Round(beatFraction * divisionsPerBeat))
	
	return MusicPosition{
		Measure:  measure,
		Beat:     beat,
		Division: division,
	}
}

// ToBeats converts to total beat count
func (mp MusicPosition) ToBeats(beatsPerMeasure, timeDivision int) float64 {
	// Convert to 0-based for calculation
	measureIndex := float64(mp.Measure - 1)
	beatIndex := float64(mp.Beat - 1)
	
	// Calculate total beats up to this measure
	totalBeats := measureIndex * float64(beatsPerMeasure)
	
	// Add beats within current measure
	totalBeats += beatIndex
	
	// Add fractional beat from division
	divisionsPerBeat := float64(timeDivision) / 4.0
	if divisionsPerBeat > 0 {
		totalBeats += float64(mp.Division) / divisionsPerBeat
	}
	
	return totalBeats
}

// ToMilliseconds converts to millisecond timestamp
func (mp MusicPosition) ToMilliseconds(bpm float64, beatsPerMeasure, timeDivision int) int64 {
	totalBeats := mp.ToBeats(beatsPerMeasure, timeDivision)
	quarterNoteMs := 60000.0 / bpm
	return int64(math.Round(totalBeats * quarterNoteMs))
}

// Add adds time to this position
func (mp MusicPosition) Add(measures, beats, divisions int, beatsPerMeasure, timeDivision int) MusicPosition {
	result := mp
	
	// Add divisions
	result.Division += divisions
	divisionsPerBeat := timeDivision / 4
	if divisionsPerBeat > 0 && result.Division >= divisionsPerBeat {
		extraBeats := result.Division / divisionsPerBeat
		result.Beat += extraBeats
		result.Division %= divisionsPerBeat
	}
	
	// Add beats
	result.Beat += beats
	if result.Beat > beatsPerMeasure {
		extraMeasures := (result.Beat - 1) / beatsPerMeasure
		result.Measure += extraMeasures
		result.Beat = ((result.Beat - 1) % beatsPerMeasure) + 1
	}
	
	// Add measures
	result.Measure += measures
	
	return result
}

// IsZero returns true if this is the zero position
func (mp MusicPosition) IsZero() bool {
	return mp.Measure == 0 && mp.Beat == 0 && mp.Division == 0
}

// String returns a string representation
func (mp MusicPosition) String() string {
	return fmt.Sprintf("%d:%d:%d", mp.Measure, mp.Beat, mp.Division)
}

// Compare compares two positions (-1: less, 0: equal, 1: greater)
func (mp MusicPosition) Compare(other MusicPosition) int {
	if mp.Measure < other.Measure {
		return -1
	}
	if mp.Measure > other.Measure {
		return 1
	}
	
	if mp.Beat < other.Beat {
		return -1
	}
	if mp.Beat > other.Beat {
		return 1
	}
	
	if mp.Division < other.Division {
		return -1
	}
	if mp.Division > other.Division {
		return 1
	}
	
	return 0
}

// Before returns true if this position is before the other
func (mp MusicPosition) Before(other MusicPosition) bool {
	return mp.Compare(other) < 0
}

// After returns true if this position is after the other
func (mp MusicPosition) After(other MusicPosition) bool {
	return mp.Compare(other) > 0
}

// Equal returns true if positions are equal
func (mp MusicPosition) Equal(other MusicPosition) bool {
	return mp.Compare(other) == 0
}