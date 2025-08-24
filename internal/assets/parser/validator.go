package parser

import (
	"fmt"
	"math"
	"sort"

	"github.com/liqmix/slaptrax/internal/types/schema"
)

// Validator provides validation for song data
type Validator struct {
	maxBPM        int
	minBPM        int
	maxDuration   int64 // milliseconds
	maxNoteCount  int
	maxEventCount int
}

// NewValidator creates a new validator with default limits
func NewValidator() *Validator {
	return &Validator{
		maxBPM:        300,
		minBPM:        60,
		maxDuration:   600000, // 10 minutes
		maxNoteCount:  10000,  // Per chart
		maxEventCount: 1000,   // Per chart
	}
}

// ValidationConfig allows customizing validation limits
type ValidationConfig struct {
	MaxBPM        int
	MinBPM        int
	MaxDuration   int64
	MaxNoteCount  int
	MaxEventCount int
}

// NewValidatorWithConfig creates a validator with custom limits
func NewValidatorWithConfig(config ValidationConfig) *Validator {
	return &Validator{
		maxBPM:        config.MaxBPM,
		minBPM:        config.MinBPM,
		maxDuration:   config.MaxDuration,
		maxNoteCount:  config.MaxNoteCount,
		maxEventCount: config.MaxEventCount,
	}
}

// ValidateSchema performs comprehensive validation on song schema
func (v *Validator) ValidateSchema(song *schema.SongDataV2) error {
	// Basic schema validation (already done in schema.FromJSON)
	if err := song.Validate(); err != nil {
		return err
	}
	
	// Extended validation
	if err := v.validateMetadata(&song.Metadata); err != nil {
		return fmt.Errorf("metadata validation failed: %w", err)
	}
	
	if err := v.validateAudio(&song.Audio); err != nil {
		return fmt.Errorf("audio validation failed: %w", err)
	}
	
	for diffStr, chart := range song.Charts {
		if err := v.validateChart(&chart); err != nil {
			return fmt.Errorf("chart %s validation failed: %w", diffStr, err)
		}
	}
	
	return nil
}

// validateMetadata validates song metadata
func (v *Validator) validateMetadata(meta *schema.SongMetadata) error {
	// BPM validation
	if meta.BPM < v.minBPM || meta.BPM > v.maxBPM {
		return fmt.Errorf("BPM %d outside valid range [%d, %d]", meta.BPM, v.minBPM, v.maxBPM)
	}
	
	// Duration validation
	if meta.Duration > v.maxDuration {
		return fmt.Errorf("song duration %dms exceeds maximum %dms", meta.Duration, v.maxDuration)
	}
	
	// Preview validation
	if meta.PreviewStart < 0 {
		return fmt.Errorf("preview start cannot be negative")
	}
	
	if meta.Duration > 0 && meta.PreviewStart >= meta.Duration {
		return fmt.Errorf("preview start %dms is beyond song duration %dms", 
			meta.PreviewStart, meta.Duration)
	}
	
	// Year validation
	if meta.Year > 0 && (meta.Year < 1900 || meta.Year > 2100) {
		return fmt.Errorf("year %d seems unrealistic", meta.Year)
	}
	
	// Difficulty range validation
	if len(meta.DifficultyRange) == 2 {
		min, max := meta.DifficultyRange[0], meta.DifficultyRange[1]
		if min < 1 || max > 10 || min > max {
			return fmt.Errorf("invalid difficulty range [%d, %d]", min, max)
		}
	}
	
	return nil
}

// validateAudio validates audio configuration
func (v *Validator) validateAudio(audio *schema.AudioInfo) error {
	if audio.File == "" {
		return fmt.Errorf("audio file not specified")
	}
	
	// Validate audio file extension
	validExtensions := []string{".ogg", ".mp3", ".wav"}
	hasValidExt := false
	for _, ext := range validExtensions {
		if len(audio.File) >= len(ext) && 
		   audio.File[len(audio.File)-len(ext):] == ext {
			hasValidExt = true
			break
		}
	}
	
	if !hasValidExt {
		return fmt.Errorf("audio file must have extension: %v", validExtensions)
	}
	
	// Preview duration validation
	if audio.PreviewDuration < 0 {
		return fmt.Errorf("preview duration cannot be negative")
	}
	
	if audio.PreviewDuration > 60000 { // 1 minute max preview
		return fmt.Errorf("preview duration %dms exceeds maximum 60000ms", audio.PreviewDuration)
	}
	
	return nil
}

// validateChart validates chart data
func (v *Validator) validateChart(chart *schema.ChartDataV2) error {
	// Basic validation already done in schema
	
	// Note count validation
	if chart.NoteCount > v.maxNoteCount {
		return fmt.Errorf("note count %d exceeds maximum %d", chart.NoteCount, v.maxNoteCount)
	}
	
	// Event count validation
	if len(chart.Events) > v.maxEventCount {
		return fmt.Errorf("event count %d exceeds maximum %d", len(chart.Events), v.maxEventCount)
	}
	
	// Validate note timing
	if err := v.validateNoteTiming(chart); err != nil {
		return fmt.Errorf("note timing validation failed: %w", err)
	}
	
	// Validate events
	if err := v.validateEvents(chart.Events); err != nil {
		return fmt.Errorf("event validation failed: %w", err)
	}
	
	// Validate difficulty vs note density - DISABLED
	// Charts can have any difficulty regardless of note density
	// if err := v.validateDifficultyConsistency(chart); err != nil {
	//	return fmt.Errorf("difficulty consistency validation failed: %w", err)
	// }
	
	return nil
}

// validateNoteTiming validates note timing consistency
func (v *Validator) validateNoteTiming(chart *schema.ChartDataV2) error {
	var allNotes []schema.NoteData
	
	// Collect all notes from all tracks
	for trackName, notes := range chart.Tracks {
		for i, note := range notes {
			if note.Time < 0 {
				return fmt.Errorf("note %d in track %s has negative time %d", i, trackName, note.Time)
			}
			
			if note.Type == schema.NoteTypeHold && note.Duration <= 0 {
				return fmt.Errorf("hold note %d in track %s has invalid duration %d", 
					i, trackName, note.Duration)
			}
			
			allNotes = append(allNotes, note)
		}
	}
	
	// Sort notes by time
	sort.Slice(allNotes, func(i, j int) bool {
		return allNotes[i].Time < allNotes[j].Time
	})
	
	// Check for reasonable spacing between notes
	const minNoteSpacing = 50 // 50ms minimum between notes
	for i := 1; i < len(allNotes); i++ {
		prev, curr := allNotes[i-1], allNotes[i]
		if curr.Time-prev.Time < minNoteSpacing && curr.Time != prev.Time {
			return fmt.Errorf("notes too close together: %dms and %dms (minimum spacing: %dms)",
				prev.Time, curr.Time, minNoteSpacing)
		}
	}
	
	return nil
}

// validateEvents validates event data
func (v *Validator) validateEvents(events []schema.EventData) error {
	// Sort events by time for validation
	sortedEvents := make([]schema.EventData, len(events))
	copy(sortedEvents, events)
	sort.Slice(sortedEvents, func(i, j int) bool {
		return sortedEvents[i].Time < sortedEvents[j].Time
	})
	
	for i, event := range sortedEvents {
		if err := v.validateEvent(event); err != nil {
			return fmt.Errorf("event %d validation failed: %w", i, err)
		}
	}
	
	// Check for event conflicts
	if err := v.checkEventConflicts(sortedEvents); err != nil {
		return err
	}
	
	return nil
}

// validateEvent validates a single event
func (v *Validator) validateEvent(event schema.EventData) error {
	if event.Time < 0 {
		return fmt.Errorf("event time cannot be negative")
	}
	
	if event.Duration < 0 {
		return fmt.Errorf("event duration cannot be negative")
	}
	
	switch event.Type {
	case schema.EventTypeBPMChange:
		return v.validateBPMChangeEvent(event)
	case schema.EventTypeLaneEffect:
		return v.validateLaneEffectEvent(event)
	case schema.EventTypeBackgroundChange:
		return v.validateBackgroundChangeEvent(event)
	case schema.EventTypeSpeedChange:
		return v.validateSpeedChangeEvent(event)
	case schema.EventTypeColorChange:
		return v.validateColorChangeEvent(event)
	case schema.EventTypeParticleEffect:
		return v.validateParticleEffectEvent(event)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

// validateBPMChangeEvent validates BPM change events
func (v *Validator) validateBPMChangeEvent(event schema.EventData) error {
	bpmVal, exists := event.Properties["bpm"]
	if !exists {
		return fmt.Errorf("BPM change event missing 'bpm' property")
	}
	
	bpm, ok := bpmVal.(float64)
	if !ok {
		return fmt.Errorf("BPM must be a number")
	}
	
	if bpm < float64(v.minBPM) || bpm > float64(v.maxBPM) {
		return fmt.Errorf("BPM %g outside valid range [%d, %d]", bpm, v.minBPM, v.maxBPM)
	}
	
	return nil
}

// validateLaneEffectEvent validates lane effect events
func (v *Validator) validateLaneEffectEvent(event schema.EventData) error {
	if event.Target == "" {
		return fmt.Errorf("lane effect event missing target track")
	}
	
	validTracks := map[string]bool{
		"all":                      true,
		schema.TrackLeftTop:        true,
		schema.TrackLeftBottom:     true,
		schema.TrackCenterTop:      true,
		schema.TrackCenterBottom:   true,
		schema.TrackRightTop:       true,
		schema.TrackRightBottom:    true,
	}
	
	if !validTracks[event.Target] {
		return fmt.Errorf("invalid target track for lane effect: %s", event.Target)
	}
	
	if event.Duration <= 0 {
		return fmt.Errorf("lane effect event must have positive duration")
	}
	
	return nil
}

// validateBackgroundChangeEvent validates background change events
func (v *Validator) validateBackgroundChangeEvent(event schema.EventData) error {
	imageVal, exists := event.Properties["image"]
	if !exists {
		return fmt.Errorf("background change event missing 'image' property")
	}
	
	image, ok := imageVal.(string)
	if !ok || image == "" {
		return fmt.Errorf("background image must be a non-empty string")
	}
	
	return nil
}

// validateSpeedChangeEvent validates speed change events
func (v *Validator) validateSpeedChangeEvent(event schema.EventData) error {
	multiplierVal, exists := event.Properties["multiplier"]
	if !exists {
		return fmt.Errorf("speed change event missing 'multiplier' property")
	}
	
	multiplier, ok := multiplierVal.(float64)
	if !ok {
		return fmt.Errorf("speed multiplier must be a number")
	}
	
	if multiplier <= 0 || multiplier > 5.0 {
		return fmt.Errorf("speed multiplier %g outside valid range (0, 5.0]", multiplier)
	}
	
	return nil
}

// validateColorChangeEvent validates color change events
func (v *Validator) validateColorChangeEvent(event schema.EventData) error {
	if event.Target == "" {
		return fmt.Errorf("color change event missing target track")
	}
	
	colorVal, exists := event.Properties["color"]
	if !exists {
		return fmt.Errorf("color change event missing 'color' property")
	}
	
	color, ok := colorVal.(string)
	if !ok || color == "" {
		return fmt.Errorf("color must be a non-empty string")
	}
	
	// Basic hex color validation
	if len(color) == 7 && color[0] == '#' {
		for i := 1; i < 7; i++ {
			c := color[i]
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return fmt.Errorf("invalid hex color: %s", color)
			}
		}
	} else if len(color) != 7 || color[0] != '#' {
		return fmt.Errorf("color must be in hex format (#RRGGBB)")
	}
	
	return nil
}

// validateParticleEffectEvent validates particle effect events
func (v *Validator) validateParticleEffectEvent(event schema.EventData) error {
	if event.Effect == "" {
		return fmt.Errorf("particle effect event missing effect type")
	}
	
	if event.Duration <= 0 {
		return fmt.Errorf("particle effect event must have positive duration")
	}
	
	// Validate optional properties
	if intensityVal, exists := event.Properties["intensity"]; exists {
		intensity, ok := intensityVal.(float64)
		if !ok || intensity < 0 || intensity > 2.0 {
			return fmt.Errorf("particle intensity must be between 0 and 2.0")
		}
	}
	
	if countVal, exists := event.Properties["particle_count"]; exists {
		count, ok := countVal.(float64)
		if !ok || count < 1 || count > 1000 {
			return fmt.Errorf("particle count must be between 1 and 1000")
		}
	}
	
	return nil
}

// checkEventConflicts checks for conflicting events
func (v *Validator) checkEventConflicts(events []schema.EventData) error {
	// Check for multiple BPM changes at the same time
	bpmTimes := make(map[int64]bool)
	
	for _, event := range events {
		if event.Type == schema.EventTypeBPMChange {
			if bpmTimes[event.Time] {
				return fmt.Errorf("multiple BPM changes at time %d", event.Time)
			}
			bpmTimes[event.Time] = true
		}
	}
	
	return nil
}

// validateDifficultyConsistency checks if note density matches difficulty
func (v *Validator) validateDifficultyConsistency(chart *schema.ChartDataV2) error {
	// Calculate notes per second
	if chart.NoteCount == 0 {
		return nil // Empty chart is OK
	}
	
	// Find the time span of the chart
	var minTime, maxTime int64 = math.MaxInt64, 0
	
	for _, notes := range chart.Tracks {
		for _, note := range notes {
			if note.Time < minTime {
				minTime = note.Time
			}
			noteEnd := note.Time
			if note.Duration > 0 {
				noteEnd = note.Time + note.Duration
			}
			if noteEnd > maxTime {
				maxTime = noteEnd
			}
		}
	}
	
	if minTime >= maxTime {
		return nil // No valid time span
	}
	
	duration := float64(maxTime-minTime) / 1000.0 // Convert to seconds
	notesPerSecond := float64(chart.NoteCount) / duration
	
	// Expected ranges for difficulty levels (notes per second)
	expectedRanges := map[int][2]float64{
		1:  {0.5, 2.0},   // Easy
		2:  {1.0, 3.0},   // 
		3:  {2.0, 4.0},   // 
		4:  {3.0, 5.0},   // 
		5:  {4.0, 6.0},   // Medium
		6:  {5.0, 7.0},   // 
		7:  {6.0, 8.0},   // 
		8:  {7.0, 9.0},   // 
		9:  {8.0, 10.0},  // 
		10: {9.0, 12.0},  // Hard
	}
	
	if expectedRange, exists := expectedRanges[chart.Difficulty]; exists {
		min, max := expectedRange[0], expectedRange[1]
		if notesPerSecond < min || notesPerSecond > max {
			return fmt.Errorf("difficulty %d expects %.1f-%.1f notes/sec, got %.1f", 
				chart.Difficulty, min, max, notesPerSecond)
		}
	}
	
	return nil
}