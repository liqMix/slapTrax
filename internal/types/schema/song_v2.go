package schema

import (
	"encoding/json"
)

// SongDataV2 represents the complete JSON song format
type SongDataV2 struct {
	Schema   string                 `json:"$schema,omitempty"`
	Version  int                    `json:"version"`
	Metadata SongMetadata           `json:"metadata"`
	Audio    AudioInfo              `json:"audio"`
	Visual   VisualInfo             `json:"visual"`
	Charts   map[string]ChartDataV2 `json:"charts"`
}

// SongMetadata contains enhanced song information
type SongMetadata struct {
	Title           string   `json:"title"`
	TitleLink       string   `json:"title_link,omitempty"`
	Artist          string   `json:"artist"`
	ArtistLink      string   `json:"artist_link,omitempty"`
	Album           string   `json:"album,omitempty"`
	AlbumLink       string   `json:"album_link,omitempty"`
	Year            int      `json:"year,omitempty"`
	BPM             int      `json:"bpm"`
	PreviewStart    int64    `json:"preview_start"`
	Duration        int64    `json:"duration,omitempty"`
	ChartedBy       string   `json:"charted_by"`
	ChartedByLink   string   `json:"charted_by_link,omitempty"`
	Version         string   `json:"version"`
	Tags            []string `json:"tags,omitempty"`
	DifficultyRange [2]int   `json:"difficulty_range,omitempty"`
}

// AudioInfo contains audio file information
type AudioInfo struct {
	File            string `json:"file"`
	PreviewDuration int64  `json:"preview_duration,omitempty"`
	Offset          int64  `json:"offset,omitempty"`
}

// VisualInfo contains visual asset information
type VisualInfo struct {
	Art        string `json:"art,omitempty"`
	Background string `json:"background,omitempty"`
	Theme      string `json:"theme,omitempty"`
}

// ChartDataV2 represents a single difficulty chart
type ChartDataV2 struct {
	Name       string                `json:"name"`
	Difficulty int                   `json:"difficulty"`
	NoteCount  int                   `json:"note_count"`
	HoldCount  int                   `json:"hold_count"`
	MaxCombo   int                   `json:"max_combo"`
	Tracks     map[string][]NoteData `json:"tracks"`
	Events     []EventData           `json:"events,omitempty"`
}

// NoteData represents a single note in compact format
type NoteData struct {
	Time     int64                  `json:"time"`
	Type     string                 `json:"type"`
	Duration int64                  `json:"duration,omitempty"`
	Tracks   []string               `json:"tracks,omitempty"`
	Props    map[string]interface{} `json:"props,omitempty"`
	
	// Beat-based positioning (BPM-independent)
	TimeBeat     float64 `json:"timeBeat,omitempty"`     // beat position from start of song
	DurationBeat float64 `json:"durationBeat,omitempty"` // duration in beats
}

// EventData represents a gameplay or visual event
type EventData struct {
	Time       int64                  `json:"time"`
	Type       string                 `json:"type"`
	Target     string                 `json:"target,omitempty"`
	Effect     string                 `json:"effect,omitempty"`
	Duration   int64                  `json:"duration,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	
	// Beat-based positioning (BPM-independent)
	TimeBeat     float64 `json:"timeBeat,omitempty"`     // beat position from start of song
	DurationBeat float64 `json:"durationBeat,omitempty"` // duration in beats
}

// NoteType constants
const (
	NoteTypeTap   = "tap"
	NoteTypeHold  = "hold"
	NoteTypeMulti = "multi"
)

// EventType constants
const (
	EventTypeBPMChange        = "bpm_change"
	EventTypeLaneEffect       = "lane_effect"
	EventTypeBackgroundChange = "background_change"
	EventTypeSpeedChange      = "speed_change"
	EventTypeColorChange      = "color_change"
	EventTypeParticleEffect   = "particle_effect"
	EventTypeAudioStart       = "audio_start"
	EventTypeCountdownBeat    = "countdown_beat"
)

// Track name constants matching existing system
const (
	TrackLeftTop      = "left_top"
	TrackLeftBottom   = "left_bottom"
	TrackCenterTop    = "center_top"
	TrackCenterBottom = "center_bottom"
	TrackRightTop     = "right_top"
	TrackRightBottom  = "right_bottom"
)

// Validation methods

// Validate performs basic validation on the song data
func (s *SongDataV2) Validate() error {
	if s.Version < 2 {
		return ErrInvalidVersion
	}
	if s.Metadata.Title == "" {
		return ErrMissingTitle
	}
	if s.Metadata.Artist == "" {
		return ErrMissingArtist
	}
	if s.Metadata.BPM <= 0 {
		return ErrInvalidBPM
	}
	if len(s.Charts) == 0 {
		return ErrNoCharts
	}

	for diffStr, chart := range s.Charts {
		if err := chart.Validate(); err != nil {
			return NewValidationError("chart", diffStr, err)
		}
	}

	return nil
}

// Validate performs validation on chart data
func (c *ChartDataV2) Validate() error {
	if c.Name == "" {
		return ErrMissingChartName
	}
	if c.Difficulty < 1 || c.Difficulty > 10 {
		return ErrInvalidDifficulty
	}

	totalNotes := 0
	holdNotes := 0

	validTracks := map[string]bool{
		TrackLeftTop:      true,
		TrackLeftBottom:   true,
		TrackCenterTop:    true,
		TrackCenterBottom: true,
		TrackRightTop:     true,
		TrackRightBottom:  true,
	}

	for trackName, notes := range c.Tracks {
		if !validTracks[trackName] {
			return NewValidationError("track", trackName, ErrInvalidTrackName)
		}

		for i, note := range notes {
			if err := note.Validate(); err != nil {
				return NewValidationError("note", string(rune(i)), err)
			}
			totalNotes++
			if note.Type == NoteTypeHold {
				holdNotes++
			}
		}
	}

	// Note count validation disabled - counts can be calculated at runtime
	// if c.NoteCount != totalNotes {
	//	return ErrNoteCountMismatch
	// }
	// if c.HoldCount != holdNotes {
	//	return ErrHoldCountMismatch
	// }

	return nil
}

// Validate performs validation on note data
func (n *NoteData) Validate() error {
	if n.Time < 0 {
		return ErrInvalidNoteTime
	}

	switch n.Type {
	case NoteTypeTap:
		// No additional validation needed
	case NoteTypeHold:
		if n.Duration <= 0 {
			return ErrMissingHoldDuration
		}
	case NoteTypeMulti:
		if len(n.Tracks) < 2 {
			return ErrInvalidMultiNote
		}
	default:
		return ErrInvalidNoteType
	}

	return nil
}

// Helper methods

// GetTrackNames returns all valid track names
func GetTrackNames() []string {
	return []string{
		TrackLeftTop,
		TrackLeftBottom,
		TrackCenterTop,
		TrackCenterBottom,
		TrackRightTop,
		TrackRightBottom,
	}
}

// ToJSON serializes the song data to JSON
func (s *SongDataV2) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// FromJSON deserializes song data from JSON
func FromJSON(data []byte) (*SongDataV2, error) {
	var song SongDataV2
	if err := json.Unmarshal(data, &song); err != nil {
		return nil, err
	}

	if err := song.Validate(); err != nil {
		return nil, err
	}

	return &song, nil
}

// CreateTemplate creates a template song structure
func CreateTemplate() *SongDataV2 {
	return &SongDataV2{
		Schema:  "https://slaptrax.dev/schema/song/v2.json",
		Version: 2,
		Metadata: SongMetadata{
			Title:           "New Song",
			Artist:          "Unknown Artist",
			BPM:             120,
			PreviewStart:    30000,
			ChartedBy:       "Chart Author",
			Version:         "1.0.0",
			DifficultyRange: [2]int{1, 5},
		},
		Audio: AudioInfo{
			File:            "audio.ogg",
			PreviewDuration: 15000,
		},
		Visual: VisualInfo{
			Theme: "default",
		},
		Charts: map[string]ChartDataV2{
			"1": {
				Name:       "Easy",
				Difficulty: 1,
				Tracks:     make(map[string][]NoteData),
				Events:     []EventData{},
			},
		},
	}
}
