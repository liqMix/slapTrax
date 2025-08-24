package types

import (
	"encoding/json"
	"errors"

	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types/schema"
)

type Chart struct {
	TotalNotes     int
	TotalHoldNotes int
	Tracks         []*Track
	EventManager   *EventManager // Event system for visual/gameplay effects
}

func NewChart(song *Song, data []byte) (*Chart, error) {
	logger.Debug("Parsing JSON chart for %s", song.Title)
	
	// Parse the JSON chart data
	var chartData schema.ChartDataV2
	if err := json.Unmarshal(data, &chartData); err != nil {
		return nil, errors.New("failed to parse JSON chart data: " + err.Error())
	}
	
	chart := &Chart{
		EventManager: NewEventManager(),
		TotalNotes:   chartData.NoteCount,
		TotalHoldNotes: chartData.HoldCount,
	}
	
	chart.Tracks = make([]*Track, 0)
	notes := make(map[TrackName][]*Note)
	
	// Initialize all tracks
	for _, name := range TrackNames() {
		notes[name] = []*Note{}
	}
	
	// Convert JSON note data to internal Note format
	for trackNameStr, jsonNotes := range chartData.Tracks {
		trackName := stringToTrackName(trackNameStr)
		
		// Skip unknown track names
		if trackName == TrackUnknown {
			logger.Debug("Unknown track name: %s", trackNameStr)
			continue
		}
		
		for _, jsonNote := range jsonNotes {
			var note *Note
			if jsonNote.Duration > 0 {
				// Hold note
				note = NewNote(trackName, jsonNote.Time, jsonNote.Time+jsonNote.Duration)
			} else {
				// Tap note
				note = NewNote(trackName, jsonNote.Time, 0)
			}
			notes[trackName] = append(notes[trackName], note)
		}
	}
	
	// Identify notes that have the same start time and mark them as non-solo
	noteCounts := make(map[int64]int)
	for _, trackNotes := range notes {
		for _, note := range trackNotes {
			noteCounts[note.Target]++
		}
	}
	
	for _, trackNotes := range notes {
		for _, note := range trackNotes {
			if noteCounts[note.Target] > 1 {
				note.SetSolo(false)
			}
		}
	}
	
	// Create tracks from notes
	beatInterval := int64((60000.0 / float64(song.BPM)) / 4) // Quarter note interval in ms
	for _, name := range TrackNames() {
		track := NewTrack(name, notes[name], beatInterval)
		chart.Tracks = append(chart.Tracks, track)
	}
	
	// Load events into the event manager
	for _, eventData := range chartData.Events {
		event, err := createEventFromData(eventData)
		if err != nil {
			logger.Debug("Failed to create event: %v", err)
			continue
		}
		chart.EventManager.AddEvent(event)
	}
	
	if chart.TotalNotes == 0 {
		return nil, errors.New("no notes found in chart")
	}
	
	logger.Debug("Loaded chart with %d notes (%d holds) and %d events", 
		chart.TotalNotes, chart.TotalHoldNotes, len(chartData.Events))
	
	return chart, nil
}

// Helper function to convert string to TrackName
func stringToTrackName(trackNameStr string) TrackName {
	switch trackNameStr {
	case "left_bottom":
		return TrackLeftBottom
	case "left_top":
		return TrackLeftTop
	case "center_bottom":
		return TrackCenterBottom
	case "center_top":
		return TrackCenterTop
	case "right_bottom":
		return TrackRightBottom
	case "right_top":
		return TrackRightTop
	default:
		return TrackUnknown
	}
}

// Helper function to create events from JSON data
func createEventFromData(eventData schema.EventData) (Event, error) {
	switch eventData.Type {
	case schema.EventTypeBPMChange:
		if bpm, ok := eventData.Properties["bpm"].(float64); ok {
			return NewBPMChangeEvent(eventData.Time, int(bpm)), nil
		}
		return nil, errors.New("missing bpm property for BPM change event")
		
	case schema.EventTypeLaneEffect:
		track := eventData.Target
		effect := eventData.Effect
		duration := eventData.Duration
		return NewLaneEffectEvent(eventData.Time, duration, track, effect), nil
		
	case schema.EventTypeBackgroundChange:
		image := eventData.Target
		transition := "fade" // Default transition
		if trans, ok := eventData.Properties["transition"].(string); ok {
			transition = trans
		}
		return NewBackgroundChangeEvent(eventData.Time, image, transition), nil
		
	case schema.EventTypeSpeedChange:
		if multiplier, ok := eventData.Properties["multiplier"].(float64); ok {
			return NewSpeedChangeEvent(eventData.Time, float32(multiplier)), nil
		}
		return nil, errors.New("missing multiplier property for speed change event")
		
	case schema.EventTypeColorChange:
		track := eventData.Target
		color := "default"
		if colorVal, ok := eventData.Properties["color"].(string); ok {
			color = colorVal
		}
		return NewColorChangeEvent(eventData.Time, track, color), nil
		
	case schema.EventTypeParticleEffect:
		effect := eventData.Effect
		duration := eventData.Duration
		return NewParticleEffectEvent(eventData.Time, duration, effect), nil
		
	default:
		return nil, errors.New("unknown event type: " + eventData.Type)
	}
}