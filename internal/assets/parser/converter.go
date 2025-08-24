package parser

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/types/schema"
	"gopkg.in/yaml.v2"
)

// Converter handles conversion between different song formats
type Converter struct {
	parser *JSONParser
}

// NewConverter creates a new format converter
func NewConverter() *Converter {
	return &Converter{
		parser: NewJSONParser(),
	}
}

// MIDIToJSON converts a MIDI+YAML song to JSON format
func (c *Converter) MIDIToJSON(songData *types.SongData) (*schema.SongDataV2, error) {
	logger.Info("Converting MIDI song %s to JSON format", songData.FolderName)

	// Parse existing metadata
	var oldMeta struct {
		Title         string `yaml:"title"`
		TitleLink     string `yaml:"title_link"`
		Artist        string `yaml:"artist"`
		ArtistLink    string `yaml:"artist_link"`
		Album         string `yaml:"album"`
		AlbumLink     string `yaml:"album_link"`
		Year          int    `yaml:"year"`
		BPM           int    `yaml:"bpm"`
		PreviewStart  int64  `yaml:"preview_start"`
		ChartedBy     string `yaml:"charted_by"`
		ChartedByLink string `yaml:"charted_by_link"`
		Version       string `yaml:"version"`
	}

	if err := yaml.Unmarshal(songData.Meta, &oldMeta); err != nil {
		return nil, fmt.Errorf("failed to parse existing metadata: %w", err)
	}

	// Create new JSON structure
	jsonSong := &schema.SongDataV2{
		Schema:  "https://slaptrax.dev/schema/song/v2.json",
		Version: 2,
		Metadata: schema.SongMetadata{
			Title:         oldMeta.Title,
			TitleLink:     oldMeta.TitleLink,
			Artist:        oldMeta.Artist,
			ArtistLink:    oldMeta.ArtistLink,
			Album:         oldMeta.Album,
			AlbumLink:     oldMeta.AlbumLink,
			Year:          oldMeta.Year,
			BPM:           oldMeta.BPM,
			PreviewStart:  oldMeta.PreviewStart,
			ChartedBy:     oldMeta.ChartedBy,
			ChartedByLink: oldMeta.ChartedByLink,
			Version:       oldMeta.Version,
		},
		Audio: schema.AudioInfo{
			File:            filepath.Base(songData.AudioPath),
			PreviewDuration: 15000, // Default 15 seconds
		},
		Visual: schema.VisualInfo{
			Art:   "art.png",
			Theme: "default",
		},
		Charts: make(map[string]schema.ChartDataV2),
	}

	// Convert each chart
	for difficulty, chartData := range songData.Charts {
		chartV2, err := c.convertMIDIChart(chartData, difficulty, oldMeta.BPM)
		if err != nil {
			return nil, fmt.Errorf("failed to convert chart difficulty %d: %w", difficulty, err)
		}

		jsonSong.Charts[strconv.Itoa(difficulty)] = *chartV2
	}

	logger.Info("Successfully converted MIDI song to JSON with %d charts", len(jsonSong.Charts))
	return jsonSong, nil
}

// convertMIDIChart converts MIDI chart data to JSON format
func (c *Converter) convertMIDIChart(midiData []byte, difficulty int, bpm int) (*schema.ChartDataV2, error) {
	// Parse the MIDI data using existing logic
	song := &types.Song{BPM: bpm} // Temporary song for parsing
	chart, err := types.NewChart(song, midiData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MIDI chart: %w", err)
	}

	// Convert to JSON format
	chartV2 := &schema.ChartDataV2{
		Name:       getDifficultyName(difficulty),
		Difficulty: difficulty,
		NoteCount:  chart.TotalNotes,
		HoldCount:  chart.TotalHoldNotes,
		MaxCombo:   chart.TotalNotes + chart.TotalHoldNotes, // Estimate
		Tracks:     make(map[string][]schema.NoteData),
		Events:     []schema.EventData{}, // MIDI doesn't have events
	}

	// Convert track mapping
	trackMapping := map[types.TrackName]string{
		types.TrackLeftTop:      schema.TrackLeftTop,
		types.TrackLeftBottom:   schema.TrackLeftBottom,
		types.TrackCenterTop:    schema.TrackCenterTop,
		types.TrackCenterBottom: schema.TrackCenterBottom,
		types.TrackRightTop:     schema.TrackRightTop,
		types.TrackRightBottom:  schema.TrackRightBottom,
	}

	// Initialize all tracks
	for _, trackName := range schema.GetTrackNames() {
		chartV2.Tracks[trackName] = []schema.NoteData{}
	}

	// Convert notes
	for _, track := range chart.Tracks {
		trackName, ok := trackMapping[track.Name]
		if !ok {
			logger.Warn("Unknown track name in MIDI: %s", track.Name)
			continue
		}

		var notes []schema.NoteData
		for _, note := range track.AllNotes {
			noteData := schema.NoteData{
				Time: note.Target,
				Type: schema.NoteTypeTap,
			}

			if note.IsHoldNote() {
				noteData.Type = schema.NoteTypeHold
				noteData.Duration = note.TargetRelease - note.Target
			}

			notes = append(notes, noteData)
		}

		chartV2.Tracks[trackName] = notes
	}

	return chartV2, nil
}

// getDifficultyName returns a human-readable name for difficulty level
func getDifficultyName(difficulty int) string {
	names := map[int]string{
		1:  "Beginner",
		2:  "Easy",
		3:  "Normal",
		4:  "Hard",
		5:  "Expert",
		6:  "Master",
		7:  "Insane",
		8:  "Nightmare",
		9:  "Impossible",
		10: "Godlike",
	}

	if name, ok := names[difficulty]; ok {
		return name
	}

	return fmt.Sprintf("Level %d", difficulty)
}

// JSONToMIDI converts JSON format back to MIDI (for compatibility)
func (c *Converter) JSONToMIDI(jsonSong *schema.SongDataV2) (*types.SongData, error) {
	logger.Warn("JSON to MIDI conversion not fully implemented")
	// This would be a complex conversion back to MIDI format
	// For now, we'll focus on the JSON->Game objects path
	return nil, fmt.Errorf("JSON to MIDI conversion not implemented")
}

// EstimateFileSize estimates the JSON file size for a song
func (c *Converter) EstimateFileSize(song *schema.SongDataV2) int {
	// Rough estimation based on content
	baseSize := 1000 // Base metadata size

	for _, chart := range song.Charts {
		chartSize := 200 // Chart metadata

		for _, notes := range chart.Tracks {
			chartSize += len(notes) * 50 // ~50 bytes per note
		}

		chartSize += len(chart.Events) * 100 // ~100 bytes per event

		baseSize += chartSize
	}

	return baseSize
}

// CompressionRatio estimates compression ratio of JSON vs MIDI
func (c *Converter) CompressionRatio(midiSize, jsonSize int) float64 {
	if midiSize == 0 {
		return 0
	}
	return float64(jsonSize) / float64(midiSize)
}

// ValidateConversion ensures conversion maintains data integrity
func (c *Converter) ValidateConversion(original *types.SongData, converted *schema.SongDataV2) error {
	// Parse both versions
	originalSong, err := types.NewSong(original)
	if err != nil {
		return fmt.Errorf("failed to parse original song: %w", err)
	}

	convertedSong, err := c.parser.ParseSongData(mustMarshalJSON(converted))
	if err != nil {
		return fmt.Errorf("failed to parse converted song: %w", err)
	}

	// Compare metadata
	if originalSong.Title != convertedSong.Title {
		return fmt.Errorf("title mismatch: %s != %s", originalSong.Title, convertedSong.Title)
	}

	if originalSong.Artist != convertedSong.Artist {
		return fmt.Errorf("artist mismatch: %s != %s", originalSong.Artist, convertedSong.Artist)
	}

	if originalSong.BPM != convertedSong.BPM {
		return fmt.Errorf("BPM mismatch: %d != %d", originalSong.BPM, convertedSong.BPM)
	}

	// Compare charts
	if len(originalSong.Charts) != len(convertedSong.Charts) {
		return fmt.Errorf("chart count mismatch: %d != %d",
			len(originalSong.Charts), len(convertedSong.Charts))
	}

	for difficulty := range originalSong.Charts {
		originalChart := originalSong.Charts[difficulty]
		convertedChart, ok := convertedSong.Charts[difficulty]
		if !ok {
			return fmt.Errorf("missing chart for difficulty %d", difficulty)
		}

		if originalChart.TotalNotes != convertedChart.TotalNotes {
			return fmt.Errorf("note count mismatch for difficulty %d: %d != %d",
				difficulty, originalChart.TotalNotes, convertedChart.TotalNotes)
		}

		if originalChart.TotalHoldNotes != convertedChart.TotalHoldNotes {
			return fmt.Errorf("hold note count mismatch for difficulty %d: %d != %d",
				difficulty, originalChart.TotalHoldNotes, convertedChart.TotalHoldNotes)
		}
	}

	logger.Info("Conversion validation passed")
	return nil
}

// BatchConvert converts multiple songs from MIDI to JSON
type BatchConverter struct {
	converter *Converter
	stats     BatchStats
}

// BatchStats tracks conversion statistics
type BatchStats struct {
	TotalSongs    int
	SuccessCount  int
	FailureCount  int
	TotalNotes    int
	TotalCharts   int
	OriginalSize  int64
	ConvertedSize int64
	Errors        []string
}

// NewBatchConverter creates a batch converter
func NewBatchConverter() *BatchConverter {
	return &BatchConverter{
		converter: NewConverter(),
		stats:     BatchStats{},
	}
}

// ConvertAll converts all songs in the provided map
func (bc *BatchConverter) ConvertAll(songs map[string]*types.SongData) map[string]*schema.SongDataV2 {
	bc.stats = BatchStats{} // Reset stats
	bc.stats.TotalSongs = len(songs)

	converted := make(map[string]*schema.SongDataV2)

	for hash, songData := range songs {
		jsonSong, err := bc.converter.MIDIToJSON(songData)
		if err != nil {
			bc.stats.FailureCount++
			bc.stats.Errors = append(bc.stats.Errors,
				fmt.Sprintf("Failed to convert %s: %v", songData.FolderName, err))
			logger.Error("Failed to convert song %s: %v", songData.FolderName, err)
			continue
		}

		// Validate conversion
		if err := bc.converter.ValidateConversion(songData, jsonSong); err != nil {
			bc.stats.FailureCount++
			bc.stats.Errors = append(bc.stats.Errors,
				fmt.Sprintf("Validation failed for %s: %v", songData.FolderName, err))
			logger.Error("Validation failed for song %s: %v", songData.FolderName, err)
			continue
		}

		converted[hash] = jsonSong
		bc.stats.SuccessCount++
		bc.stats.TotalCharts += len(jsonSong.Charts)

		// Count notes
		for _, chart := range jsonSong.Charts {
			bc.stats.TotalNotes += chart.NoteCount
		}

		// Estimate sizes
		bc.stats.ConvertedSize += int64(bc.converter.EstimateFileSize(jsonSong))
		// Original size would need to be calculated from MIDI + metadata
	}

	logger.Info("Batch conversion completed: %d/%d songs successful",
		bc.stats.SuccessCount, bc.stats.TotalSongs)

	return converted
}

// GetStats returns conversion statistics
func (bc *BatchConverter) GetStats() BatchStats {
	return bc.stats
}

// Helper function to marshal JSON (panics on error for internal use)
func mustMarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return data
}
