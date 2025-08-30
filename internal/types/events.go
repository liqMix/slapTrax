package types

import (
	"fmt"
	"sort"
	"sync"

	"github.com/liqmix/slaptrax/internal/types/schema"
)

// Event interface defines the contract for all game events
type Event interface {
	// GetTime returns when this event should trigger (in milliseconds)
	GetTime() int64

	// Execute performs the event action
	Execute(ctx *EventContext) error

	// GetType returns the event type identifier
	GetType() string

	// GetDuration returns how long the event lasts (0 for instant events)
	GetDuration() int64

	// IsActive returns true if the event is currently active
	IsActive(currentTime int64) bool

	// Reset resets the event to its initial state
	Reset()
}

// EventContext provides access to game systems for event execution
type EventContext struct {
	CurrentTime int64
	Song        *Song
	Chart       *Chart
	AudioOffset int64 // Audio offset in milliseconds

	// System accessors for event execution
	AudioInitSong        func(*Song)  // Function to initialize song
	AudioPlaySong        func()       // Function to play song
	AudioSetSongPosition func(int)    // Function to set song position in milliseconds
	AudioPlaySFX         func(string) // Function to play sound effects
	RenderSystem         interface{}  // TODO: Replace with actual render system interface
	EffectSystem         interface{}  // TODO: Replace with actual effect system interface
}

// BaseEvent provides common functionality for all events
type BaseEvent struct {
	Time     int64 // ms from start (calculated from TimeBeat)
	Duration int64
	Type     string
	Active   bool
	executed bool
	mu       sync.RWMutex
	
	// Beat-based positioning (musical time, BPM-independent)
	TimeBeat     float64 // beat position from start of song
	DurationBeat float64 // duration in beats
}

func (e *BaseEvent) GetTime() int64 {
	return e.Time
}

func (e *BaseEvent) GetType() string {
	return e.Type
}

func (e *BaseEvent) GetDuration() int64 {
	return e.Duration
}

func (e *BaseEvent) IsActive(currentTime int64) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.executed && currentTime >= e.Time {
		return true
	}

	if e.Duration > 0 {
		return currentTime >= e.Time && currentTime <= e.Time+e.Duration
	}

	return false
}

func (e *BaseEvent) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Active = false
	e.executed = false
}

// RecalculateTime updates the millisecond timestamp from beat position
func (e *BaseEvent) RecalculateTime(bpm float64) {
	if e.TimeBeat > 0 {
		quarterNoteMs := 60000.0 / bpm
		e.Time = int64(e.TimeBeat * quarterNoteMs)
		if e.DurationBeat > 0 {
			e.Duration = int64(e.DurationBeat * quarterNoteMs)
		}
	}
}

// SetBeatPosition sets the beat position and calculates millisecond timestamp
func (e *BaseEvent) SetBeatPosition(timeBeat float64, bpm float64) {
	e.TimeBeat = timeBeat
	e.RecalculateTime(bpm)
}

func (e *BaseEvent) markExecuted() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.executed = true
}

// BPMChangeEvent handles tempo changes during gameplay
type BPMChangeEvent struct {
	BaseEvent
	NewBPM int `json:"bpm"`
}

func NewBPMChangeEvent(time int64, newBPM int) *BPMChangeEvent {
	return &BPMChangeEvent{
		BaseEvent: BaseEvent{
			Time: time,
			Type: schema.EventTypeBPMChange,
		},
		NewBPM: newBPM,
	}
}

func (e *BPMChangeEvent) Execute(ctx *EventContext) error {
	e.markExecuted()
	// TODO: Implement BPM change logic
	// This would update the audio manager's tempo
	return nil
}

// LaneEffectEvent triggers visual effects on specific tracks
type LaneEffectEvent struct {
	BaseEvent
	Track     string  `json:"track"`
	Effect    string  `json:"effect"`
	Color     string  `json:"color,omitempty"`
	Intensity float32 `json:"intensity,omitempty"`
}

func NewLaneEffectEvent(time int64, duration int64, track, effect string) *LaneEffectEvent {
	return &LaneEffectEvent{
		BaseEvent: BaseEvent{
			Time:     time,
			Duration: duration,
			Type:     schema.EventTypeLaneEffect,
		},
		Track:     track,
		Effect:    effect,
		Intensity: 1.0,
	}
}

func (e *LaneEffectEvent) Execute(ctx *EventContext) error {
	e.markExecuted()
	// TODO: Implement lane effect logic
	// This would trigger visual effects in the render system
	return nil
}

// BackgroundChangeEvent changes the background during gameplay
type BackgroundChangeEvent struct {
	BaseEvent
	Image      string `json:"image"`
	Transition string `json:"transition,omitempty"`
}

func NewBackgroundChangeEvent(time int64, image, transition string) *BackgroundChangeEvent {
	return &BackgroundChangeEvent{
		BaseEvent: BaseEvent{
			Time: time,
			Type: schema.EventTypeBackgroundChange,
		},
		Image:      image,
		Transition: transition,
	}
}

func (e *BackgroundChangeEvent) Execute(ctx *EventContext) error {
	e.markExecuted()
	// TODO: Implement background change logic
	return nil
}

// SpeedChangeEvent modifies note travel speed
type SpeedChangeEvent struct {
	BaseEvent
	Multiplier float32 `json:"multiplier"`
}

func NewSpeedChangeEvent(time int64, multiplier float32) *SpeedChangeEvent {
	return &SpeedChangeEvent{
		BaseEvent: BaseEvent{
			Time: time,
			Type: schema.EventTypeSpeedChange,
		},
		Multiplier: multiplier,
	}
}

func (e *SpeedChangeEvent) Execute(ctx *EventContext) error {
	e.markExecuted()
	// TODO: Implement speed change logic
	return nil
}

// ColorChangeEvent changes track colors
type ColorChangeEvent struct {
	BaseEvent
	Track string `json:"track"`
	Color string `json:"color"`
}

func NewColorChangeEvent(time int64, track, color string) *ColorChangeEvent {
	return &ColorChangeEvent{
		BaseEvent: BaseEvent{
			Time: time,
			Type: schema.EventTypeColorChange,
		},
		Track: track,
		Color: color,
	}
}

func (e *ColorChangeEvent) Execute(ctx *EventContext) error {
	e.markExecuted()
	// TODO: Implement color change logic
	return nil
}

// ParticleEffectEvent triggers particle effects
type ParticleEffectEvent struct {
	BaseEvent
	Effect        string  `json:"effect"`
	Position      string  `json:"position,omitempty"`
	Intensity     float32 `json:"intensity,omitempty"`
	ParticleCount int     `json:"particle_count,omitempty"`
}

func NewParticleEffectEvent(time int64, duration int64, effect string) *ParticleEffectEvent {
	return &ParticleEffectEvent{
		BaseEvent: BaseEvent{
			Time:     time,
			Duration: duration,
			Type:     schema.EventTypeParticleEffect,
		},
		Effect:        effect,
		Intensity:     1.0,
		ParticleCount: 100,
	}
}

func (e *ParticleEffectEvent) Execute(ctx *EventContext) error {
	e.markExecuted()
	// TODO: Implement particle effect logic
	return nil
}

// AudioStartEvent marks where audio playback should begin
type AudioStartEvent struct {
	BaseEvent
}

func NewAudioStartEvent(time int64) *AudioStartEvent {
	return &AudioStartEvent{
		BaseEvent: BaseEvent{
			Time: time,
			Type: schema.EventTypeAudioStart,
			// TimeBeat will be calculated when BPM context is available
		},
	}
}

// NewAudioStartEventFromBeat creates an audio start event using beat position
func NewAudioStartEventFromBeat(timeBeat float64, bpm float64) *AudioStartEvent {
	quarterNoteMs := 60000.0 / bpm
	time := int64(timeBeat * quarterNoteMs)
	
	return &AudioStartEvent{
		BaseEvent: BaseEvent{
			Time:     time,
			Type:     schema.EventTypeAudioStart,
			TimeBeat: timeBeat,
		},
	}
}

func (e *AudioStartEvent) Execute(ctx *EventContext) error {
	e.markExecuted()
	
	// Debug: Print when audio start is executed
	fmt.Printf("🎶 Audio start executed at %dms (offset: %dms)\n", ctx.CurrentTime, ctx.AudioOffset)
	
	// Start the backing track audio if available
	if ctx.Song != nil && ctx.AudioInitSong != nil && ctx.AudioPlaySong != nil {
		ctx.AudioInitSong(ctx.Song)
		
		// Calculate how far into the audio file we should start playing
		// Audio offset trims silence: negative offset skips forward in the audio file
		timeFromMarker := ctx.CurrentTime - e.Time  // How far past the audio start marker we are
		audioFilePosition := timeFromMarker - ctx.AudioOffset  // Apply offset to trim silence
		
		// Set the position in the audio file
		if ctx.AudioSetSongPosition != nil {
			// Ensure we don't set a negative position (clamp to 0)
			finalPosition := int(max(0, audioFilePosition))
			ctx.AudioSetSongPosition(finalPosition)
			fmt.Printf("🎶 Setting audio file position to %dms (time from marker: %dms, offset: %dms)\n", 
				finalPosition, timeFromMarker, ctx.AudioOffset)
		}
		
		ctx.AudioPlaySong()
	}
	
	return nil
}

// CountdownBeatEvent marks countdown beats before the song starts
type CountdownBeatEvent struct {
	BaseEvent
}

func NewCountdownBeatEvent(time int64) *CountdownBeatEvent {
	return &CountdownBeatEvent{
		BaseEvent: BaseEvent{
			Time: time,
			Type: schema.EventTypeCountdownBeat,
			// TimeBeat will be calculated when BPM context is available
		},
	}
}

// NewCountdownBeatEventFromBeat creates a countdown beat event using beat position
func NewCountdownBeatEventFromBeat(timeBeat float64, bpm float64) *CountdownBeatEvent {
	quarterNoteMs := 60000.0 / bpm
	time := int64(timeBeat * quarterNoteMs)
	
	return &CountdownBeatEvent{
		BaseEvent: BaseEvent{
			Time:     time,
			Type:     schema.EventTypeCountdownBeat,
			TimeBeat: timeBeat,
		},
	}
}

func (e *CountdownBeatEvent) Execute(ctx *EventContext) error {
	e.markExecuted()
	
	// Debug: Print when countdown beat is executed
	fmt.Printf("🎵 Countdown beat executed at %dms\n", ctx.CurrentTime)
	
	// Play metronome/countdown sound
	if ctx.AudioPlaySFX != nil {
		ctx.AudioPlaySFX("select") // Using select sound which should be more audible
	}
	
	return nil
}

// EventManager handles event scheduling and execution
type EventManager struct {
	events       []Event
	activeEvents []Event
	mu           sync.RWMutex
	currentTime  int64
	nextEventIdx int // Index of the next event to execute
}

func NewEventManager() *EventManager {
	return &EventManager{
		events:       make([]Event, 0),
		activeEvents: make([]Event, 0),
		nextEventIdx: 0,
	}
}

// AddEvent adds an event to the manager
func (em *EventManager) AddEvent(event Event) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.events = append(em.events, event)

	// Keep events sorted by time
	sort.Slice(em.events, func(i, j int) bool {
		return em.events[i].GetTime() < em.events[j].GetTime()
	})
	
	// Reset queue index since events were resorted
	em.nextEventIdx = 0
}

// Update processes events for the current time
func (em *EventManager) Update(currentTime int64, ctx *EventContext) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.currentTime = currentTime
	ctx.CurrentTime = currentTime

	// Only check if we have events and haven't processed them all
	if em.nextEventIdx >= len(em.events) {
		return nil // All events processed
	}

	// Process all events whose time has come (sequential processing)
	for em.nextEventIdx < len(em.events) {
		nextEvent := em.events[em.nextEventIdx]
		
		// If this event's time hasn't come yet, we're done (events are sorted)
		if nextEvent.GetTime() > currentTime {
			break
		}
		
		// Execute the event
		fmt.Printf("⚡ Executing %s event at %dms (current: %dms)\n", 
			nextEvent.GetType(), nextEvent.GetTime(), currentTime)
		
		if err := nextEvent.Execute(ctx); err != nil {
			return fmt.Errorf("failed to execute event %s at %dms: %w",
				nextEvent.GetType(), nextEvent.GetTime(), err)
		}
		
		// Move to next event
		em.nextEventIdx++
	}

	// Update active events list (for visual/other systems that need this)
	em.activeEvents = em.activeEvents[:0] // Clear without reallocation
	for _, event := range em.events {
		if event.IsActive(currentTime) {
			em.activeEvents = append(em.activeEvents, event)
		}
	}

	return nil
}

// GetActiveEvents returns currently active events
func (em *EventManager) GetActiveEvents() []Event {
	em.mu.RLock()
	defer em.mu.RUnlock()

	// Return a copy to avoid race conditions
	active := make([]Event, len(em.activeEvents))
	copy(active, em.activeEvents)
	return active
}

// Reset resets all events to their initial state
func (em *EventManager) Reset() {
	em.mu.Lock()
	defer em.mu.Unlock()

	for _, event := range em.events {
		event.Reset()
	}
	em.activeEvents = em.activeEvents[:0]
	em.currentTime = 0
	em.nextEventIdx = 0 // Reset to start of event queue
}

// Clear removes all events
func (em *EventManager) Clear() {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.events = em.events[:0]
	em.activeEvents = em.activeEvents[:0]
	em.currentTime = 0
	em.nextEventIdx = 0 // Reset queue index
}

// GetEventCount returns the total number of events
func (em *EventManager) GetEventCount() int {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return len(em.events)
}

// CreateEventFromData creates an Event from schema.EventData
func CreateEventFromData(data schema.EventData) (Event, error) {
	return CreateEventFromDataWithBPM(data, 120.0) // Default BPM fallback
}

// CreateEventFromDataWithBPM creates an Event from schema.EventData with BPM context
func CreateEventFromDataWithBPM(data schema.EventData, bpm float64) (Event, error) {
	switch data.Type {
	case schema.EventTypeBPMChange:
		bpm, ok := data.Properties["bpm"].(float64)
		if !ok {
			return nil, fmt.Errorf("BPM change event missing bpm property")
		}
		return NewBPMChangeEvent(data.Time, int(bpm)), nil

	case schema.EventTypeLaneEffect:
		track := data.Target
		if track == "" {
			track = "all"
		}
		effect := data.Effect
		if effect == "" {
			effect = "pulse"
		}
		return NewLaneEffectEvent(data.Time, data.Duration, track, effect), nil

	case schema.EventTypeBackgroundChange:
		image, ok := data.Properties["image"].(string)
		if !ok {
			return nil, fmt.Errorf("background change event missing image property")
		}
		transition, _ := data.Properties["transition"].(string)
		return NewBackgroundChangeEvent(data.Time, image, transition), nil

	case schema.EventTypeSpeedChange:
		multiplier, ok := data.Properties["multiplier"].(float64)
		if !ok {
			return nil, fmt.Errorf("speed change event missing multiplier property")
		}
		return NewSpeedChangeEvent(data.Time, float32(multiplier)), nil

	case schema.EventTypeColorChange:
		track := data.Target
		color, ok := data.Properties["color"].(string)
		if !ok {
			return nil, fmt.Errorf("color change event missing color property")
		}
		return NewColorChangeEvent(data.Time, track, color), nil

	case schema.EventTypeParticleEffect:
		effect := data.Effect
		if effect == "" {
			effect = "explosion"
		}
		return NewParticleEffectEvent(data.Time, data.Duration, effect), nil

	case schema.EventTypeAudioStart:
		if data.TimeBeat > 0 {
			// Use beat-based positioning
			return NewAudioStartEventFromBeat(data.TimeBeat, bpm), nil
		}
		// Fall back to millisecond-based positioning
		return NewAudioStartEvent(data.Time), nil

	case schema.EventTypeCountdownBeat:
		if data.TimeBeat > 0 {
			// Use beat-based positioning
			return NewCountdownBeatEventFromBeat(data.TimeBeat, bpm), nil
		}
		// Fall back to millisecond-based positioning
		return NewCountdownBeatEvent(data.Time), nil

	default:
		return nil, fmt.Errorf("unknown event type: %s", data.Type)
	}
}

// Helper function for max of two int64 values
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
