package parser

import (
	"crypto/sha256"
	"fmt"

	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/types/schema"
)

// JSONParser handles parsing of JSON song files
type JSONParser struct {
	validator *Validator
}

// NewJSONParser creates a new JSON parser instance
func NewJSONParser() *JSONParser {
	return &JSONParser{
		validator: NewValidator(),
	}
}

// ParseSongData converts JSON song data to game objects
func (p *JSONParser) ParseSongData(data []byte) (*types.Song, error) {
	logger.Debug("Parsing JSON song data")

	// Parse JSON schema
	songData, err := schema.FromJSON(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate schema
	if err := p.validator.ValidateSchema(songData); err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	// Convert to game objects
	song, err := p.convertToSong(songData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to game objects: %w", err)
	}

	logger.Info("Successfully parsed JSON song: %s by %s", song.Title, song.Artist)
	return song, nil
}

// convertToSong converts schema.SongDataV2 to types.Song
func (p *JSONParser) convertToSong(data *schema.SongDataV2) (*types.Song, error) {
	song := &types.Song{
		Title:         data.Metadata.Title,
		TitleLink:     data.Metadata.TitleLink,
		Artist:        data.Metadata.Artist,
		ArtistLink:    data.Metadata.ArtistLink,
		Album:         data.Metadata.Album,
		AlbumLink:     data.Metadata.AlbumLink,
		Year:          data.Metadata.Year,
		BPM:           data.Metadata.BPM,
		Length:        int(data.Metadata.Duration),
		PreviewStart:  data.Metadata.PreviewStart,
		ChartedBy:     data.Metadata.ChartedBy,
		ChartedByLink: data.Metadata.ChartedByLink,
		Version:       data.Metadata.Version,
		AudioPath:     data.Audio.File,
		Charts:        make(map[types.Difficulty]*types.Chart),
	}

	// Parse charts
	for diffStr, chartData := range data.Charts {
		chart, err := p.convertToChart(song, &chartData)
		if err != nil {
			return nil, fmt.Errorf("failed to convert chart %s: %w", diffStr, err)
		}

		song.Charts[types.Difficulty(chartData.Difficulty)] = chart
	}

	// Generate hash based on song metadata and chart data
	song.Hash = p.generateHash(data)

	return song, nil
}

// generateHash creates a unique hash for the song based on its data
func (p *JSONParser) generateHash(data *schema.SongDataV2) string {
	hasher := sha256.New()
	
	// Hash metadata
	hasher.Write([]byte(data.Metadata.Title))
	hasher.Write([]byte(data.Metadata.Artist))
	hasher.Write([]byte(fmt.Sprintf("%d", data.Metadata.BPM)))
	hasher.Write([]byte(data.Audio.File))
	
	// Hash chart data in deterministic order
	for _, chartData := range data.Charts {
		hasher.Write([]byte(fmt.Sprintf("diff:%d", chartData.Difficulty)))
		hasher.Write([]byte(fmt.Sprintf("notes:%d", chartData.NoteCount)))
		
		// Include some track data to ensure uniqueness
		for trackName, notes := range chartData.Tracks {
			hasher.Write([]byte(trackName))
			for _, note := range notes {
				hasher.Write([]byte(fmt.Sprintf("t:%d:d:%d", note.Time, note.Duration)))
			}
		}
	}
	
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// convertToChart converts schema.ChartDataV2 to types.Chart
func (p *JSONParser) convertToChart(song *types.Song, data *schema.ChartDataV2) (*types.Chart, error) {
	chart := &types.Chart{
		TotalNotes:     data.NoteCount,
		TotalHoldNotes: data.HoldCount,
		Tracks:         make([]*types.Track, 0),
		EventManager:   types.NewEventManager(),
	}

	// Convert notes for each track
	notes := make(map[types.TrackName][]*types.Note)
	trackMapping := map[string]types.TrackName{
		schema.TrackLeftTop:      types.TrackLeftTop,
		schema.TrackLeftBottom:   types.TrackLeftBottom,
		schema.TrackCenterTop:    types.TrackCenterTop,
		schema.TrackCenterBottom: types.TrackCenterBottom,
		schema.TrackRightTop:     types.TrackRightTop,
		schema.TrackRightBottom:  types.TrackRightBottom,
	}

	// Initialize note arrays for all tracks
	for _, trackName := range types.TrackNames() {
		notes[trackName] = []*types.Note{}
	}

	// Process notes for each track
	for trackStr, trackNotes := range data.Tracks {
		trackName, ok := trackMapping[trackStr]
		if !ok {
			return nil, fmt.Errorf("invalid track name: %s", trackStr)
		}

		for _, noteData := range trackNotes {
			gameNotes, err := p.convertNote(noteData, trackName)
			if err != nil {
				return nil, fmt.Errorf("failed to convert note in track %s: %w", trackStr, err)
			}
			notes[trackName] = append(notes[trackName], gameNotes...)
		}
	}

	// Mark multi-notes as non-solo
	p.markMultiNotes(notes)

	// Create tracks from notes
	beatInterval := song.GetBeatInterval()
	for _, trackName := range types.TrackNames() {
		track := types.NewTrack(trackName, notes[trackName], beatInterval)
		chart.Tracks = append(chart.Tracks, track)
	}

	// Convert events
	if len(data.Events) > 0 {
		for _, eventData := range data.Events {
			event, err := types.CreateEventFromData(eventData)
			if err != nil {
				logger.Warn("Failed to create event: %v", err)
				continue
			}
			chart.EventManager.AddEvent(event)
		}

		logger.Debug("Loaded %d events for chart", chart.EventManager.GetEventCount())
	}

	return chart, nil
}

// convertNote converts schema.NoteData to types.Note(s)
func (p *JSONParser) convertNote(data schema.NoteData, defaultTrack types.TrackName) ([]*types.Note, error) {
	switch data.Type {
	case schema.NoteTypeTap:
		note := types.NewNote(defaultTrack, data.Time, 0)
		return []*types.Note{note}, nil

	case schema.NoteTypeHold:
		if data.Duration <= 0 {
			return nil, fmt.Errorf("hold note missing duration")
		}
		note := types.NewNote(defaultTrack, data.Time, data.Time+data.Duration)
		return []*types.Note{note}, nil

	case schema.NoteTypeMulti:
		if len(data.Tracks) < 2 {
			return nil, fmt.Errorf("multi note must specify at least 2 tracks")
		}

		trackMapping := map[string]types.TrackName{
			schema.TrackLeftTop:      types.TrackLeftTop,
			schema.TrackLeftBottom:   types.TrackLeftBottom,
			schema.TrackCenterTop:    types.TrackCenterTop,
			schema.TrackCenterBottom: types.TrackCenterBottom,
			schema.TrackRightTop:     types.TrackRightTop,
			schema.TrackRightBottom:  types.TrackRightBottom,
		}

		var notes []*types.Note
		for _, trackStr := range data.Tracks {
			trackName, ok := trackMapping[trackStr]
			if !ok {
				return nil, fmt.Errorf("invalid track name in multi note: %s", trackStr)
			}

			var note *types.Note
			if data.Duration > 0 {
				note = types.NewNote(trackName, data.Time, data.Time+data.Duration)
			} else {
				note = types.NewNote(trackName, data.Time, 0)
			}
			note.SetSolo(false) // Multi notes are never solo
			notes = append(notes, note)
		}

		return notes, nil

	default:
		return nil, fmt.Errorf("unknown note type: %s", data.Type)
	}
}

// markMultiNotes identifies notes that occur at the same time and marks them as non-solo
func (p *JSONParser) markMultiNotes(notes map[types.TrackName][]*types.Note) {
	// Count notes at each timestamp
	noteCounts := make(map[int64]int)
	for _, trackNotes := range notes {
		for _, note := range trackNotes {
			noteCounts[note.Target]++
		}
	}

	// Mark notes that share timestamps as non-solo
	for _, trackNotes := range notes {
		for _, note := range trackNotes {
			if noteCounts[note.Target] > 1 {
				note.SetSolo(false)
			}
		}
	}
}

// ValidateGameplay performs additional validation specific to gameplay
func (p *JSONParser) ValidateGameplay(song *types.Song) error {
	for difficulty, chart := range song.Charts {
		if err := p.validateChartGameplay(chart, difficulty); err != nil {
			return fmt.Errorf("gameplay validation failed for difficulty %d: %w", difficulty, err)
		}
	}
	return nil
}

// validateChartGameplay validates chart for gameplay issues
func (p *JSONParser) validateChartGameplay(chart *types.Chart, difficulty types.Difficulty) error {
	// Check for reasonable note density
	if chart.TotalNotes == 0 {
		return fmt.Errorf("chart has no notes")
	}

	// Validate timing consistency
	for _, track := range chart.Tracks {
		if err := p.validateTrackTiming(track); err != nil {
			return fmt.Errorf("track %s timing validation failed: %w", track.Name, err)
		}
	}

	return nil
}

// validateTrackTiming checks for timing issues in a track
func (p *JSONParser) validateTrackTiming(track *types.Track) error {
	if len(track.AllNotes) == 0 {
		return nil // Empty tracks are OK
	}

	// Check notes are sorted by time
	for i := 1; i < len(track.AllNotes); i++ {
		if track.AllNotes[i].Target < track.AllNotes[i-1].Target {
			return fmt.Errorf("notes not sorted by time")
		}
	}

	// Check for overlapping hold notes
	for i := 0; i < len(track.AllNotes); i++ {
		note := track.AllNotes[i]
		if !note.IsHoldNote() {
			continue
		}

		for j := i + 1; j < len(track.AllNotes); j++ {
			nextNote := track.AllNotes[j]
			if nextNote.Target >= note.TargetRelease {
				break // No more overlapping notes
			}

			// Found an overlapping note
			logger.Warn("Overlapping notes in track %s at times %d and %d",
				track.Name, note.Target, nextNote.Target)
		}
	}

	return nil
}

// Statistics provides parsing statistics
type Statistics struct {
	TotalCharts    int
	TotalNotes     int
	TotalHoldNotes int
	TotalEvents    int
	ParsingTime    int64 // milliseconds
}

// GetStatistics returns parsing statistics for a song
func (p *JSONParser) GetStatistics(song *types.Song) Statistics {
	stats := Statistics{
		TotalCharts: len(song.Charts),
	}

	for _, chart := range song.Charts {
		stats.TotalNotes += chart.TotalNotes
		stats.TotalHoldNotes += chart.TotalHoldNotes
		// TODO: Add event count when Chart struct is updated
	}

	return stats
}
