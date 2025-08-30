package assets

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/types/schema"
)

const (
	userChartsDir = "storage/charts"
)

// EnsureUserChartsDirectory creates the user charts directory if it doesn't exist
func EnsureUserChartsDirectory() error {
	if err := os.MkdirAll(userChartsDir, 0755); err != nil {
		return err
	}
	return nil
}

// GetUserChartsPath returns the path to the user charts directory
func GetUserChartsPath() string {
	return userChartsDir
}

// GetUserCharts scans the user charts directory for song directories
// This function is called dynamically each time to detect new charts added at runtime
func GetUserCharts() ([]*UserChart, error) {
	if err := EnsureUserChartsDirectory(); err != nil {
		return nil, err
	}

	var userCharts []*UserChart

	// Scan for song directories (each directory should contain song.json and audio.ogg)
	entries, err := os.ReadDir(userChartsDir)
	if err != nil {
		logger.Error("Failed to read user charts directory: %v", err)
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip files, we only want directories
		}

		songDir := filepath.Join(userChartsDir, entry.Name())
		songJsonPath := filepath.Join(songDir, "song.json")
		audioPath := filepath.Join(songDir, "audio.ogg")

		// Check if this is a valid song directory (has song.json and audio.ogg)
		if _, err := os.Stat(songJsonPath); err != nil {
			logger.Debug("Skipping %s: no song.json found", songDir)
			continue
		}
		if _, err := os.Stat(audioPath); err != nil {
			logger.Debug("Skipping %s: no audio.ogg found", songDir)
			continue
		}

		logger.Debug("Found valid user chart directory: %s", songDir)

		// Try to load the chart
		chart, err := LoadUserChart(songJsonPath)
		if err != nil {
			logger.Warn("Failed to load user chart %s: %v", songJsonPath, err)
			continue
		}

		// Set the correct audio path
		chart.AudioPath = audioPath
		chart.DirectoryPath = songDir

		logger.Debug("Loaded user chart: %s by %s", chart.Title, chart.Artist)
		userCharts = append(userCharts, chart)
	}

	logger.Debug("Successfully loaded %d user charts", len(userCharts))
	return userCharts, nil
}

// UserChart represents a user-created chart
type UserChart struct {
	FilePath      string // Path to song.json
	DirectoryPath string // Path to song directory
	Title         string
	Artist        string
	BPM           int
	ChartName     string
	Difficulty    int
	AudioPath     string // Path to audio.ogg
	ArtPath       string // Path to art.png (optional)
}

// LoadUserChart loads a user chart from a JSON file
func LoadUserChart(filePath string) (*UserChart, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	songData, err := schema.FromJSON(data)
	if err != nil {
		return nil, err
	}

	// Find the first chart for basic info
	var chartName string
	var difficulty int
	for _, chart := range songData.Charts {
		chartName = chart.Name
		difficulty = chart.Difficulty
		break
	}

	// Audio file should always be audio.ogg in the same directory
	chartDir := filepath.Dir(filePath)
	audioPath := filepath.Join(chartDir, "audio.ogg")
	
	// Check for optional album art
	artPath := filepath.Join(chartDir, "art.png")
	if _, err := os.Stat(artPath); err != nil {
		artPath = "" // No art file found
	}

	userChart := &UserChart{
		FilePath:      filePath,
		DirectoryPath: chartDir,
		Title:         songData.Metadata.Title,
		Artist:        songData.Metadata.Artist,
		BPM:           songData.Metadata.BPM,
		ChartName:     chartName,
		Difficulty:    difficulty,
		AudioPath:     audioPath,
		ArtPath:       artPath,
	}

	return userChart, nil
}

// ConvertUserChartToSong converts a UserChart to a full Song object with chart data
func ConvertUserChartToSong(userChart *UserChart) (*types.Song, error) {
	// Load the full chart data
	data, err := os.ReadFile(userChart.FilePath)
	if err != nil {
		return nil, err
	}

	songData, err := schema.FromJSON(data)
	if err != nil {
		return nil, err
	}

	// Create the song object
	song := &types.Song{
		Title:     songData.Metadata.Title + " (Custom)",
		Artist:    songData.Metadata.Artist,
		BPM:       songData.Metadata.BPM,
		Hash:      userChart.FilePath, // Use file path as unique identifier
		AudioPath: userChart.AudioPath,
		Charts:    make(map[types.Difficulty]*types.Chart),
	}

	// Convert all charts
	for _, chartData := range songData.Charts {
		chart := &types.Chart{
			Tracks:         make([]*types.Track, 0),
			TotalNotes:     chartData.NoteCount,
			TotalHoldNotes: chartData.HoldCount,
		}

		// Convert tracks from schema format
		for trackNameStr, noteDataList := range chartData.Tracks {
			// Convert string to TrackName
			var trackName types.TrackName
			switch trackNameStr {
			case "left_bottom":
				trackName = types.TrackLeftBottom
			case "left_top":
				trackName = types.TrackLeftTop
			case "center_bottom":
				trackName = types.TrackCenterBottom
			case "center_top":
				trackName = types.TrackCenterTop
			case "right_bottom":
				trackName = types.TrackRightBottom
			case "right_top":
				trackName = types.TrackRightTop
			default:
				trackName = types.TrackUnknown
			}

			if trackName == types.TrackUnknown {
				continue
			}

			track := &types.Track{
				Name:     trackName,
				AllNotes: make([]*types.Note, len(noteDataList)),
			}

			// Convert notes
			for i, noteData := range noteDataList {
				note := &types.Note{
					Id:        i,
					TrackName: trackName,
					Target:    noteData.Time,
				}
				
				// Handle hold notes
				if noteData.Duration > 0 {
					note.TargetRelease = noteData.Time + noteData.Duration
				} else {
					note.TargetRelease = noteData.Time
				}
				
				track.AllNotes[i] = note
			}

			chart.Tracks = append(chart.Tracks, track)
		}

		song.Charts[types.Difficulty(chartData.Difficulty)] = chart
	}

	return song, nil
}

// SaveUserChart saves a chart to the user charts directory in song format
func SaveUserChart(songData *schema.SongDataV2, audioSourcePath, title string) (string, error) {
	if err := EnsureUserChartsDirectory(); err != nil {
		return "", err
	}

	// Create song directory
	songDir, err := CreateUserChartDirectory(title)
	if err != nil {
		return "", err
	}

	// Save song.json
	songJsonPath := filepath.Join(songDir, "song.json")
	data, err := songData.ToJSON()
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(songJsonPath, data, 0644); err != nil {
		return "", err
	}

	// Copy audio file to audio.ogg
	if audioSourcePath != "" {
		audioDestPath := filepath.Join(songDir, "audio.ogg")
		if err := copyFile(audioSourcePath, audioDestPath); err != nil {
			return "", err
		}
	}

	return songDir, nil
}

// CreateUserChartDirectory creates a directory for a new user chart
func CreateUserChartDirectory(title string) (string, error) {
	if err := EnsureUserChartsDirectory(); err != nil {
		return "", err
	}

	// Sanitize title for directory name
	sanitized := strings.ReplaceAll(title, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, "\\", "_")
	sanitized = strings.ToLower(sanitized)
	
	chartDir := filepath.Join(userChartsDir, sanitized)
	if err := os.MkdirAll(chartDir, 0755); err != nil {
		return "", err
	}

	return chartDir, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceData, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, sourceData, 0644)
}