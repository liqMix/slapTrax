package assets

import (
	"bytes"
	"embed"
	"fmt"
	"path"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/liqmix/slaptrax/internal/assets/parser"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types"
)

// JSONLoader handles JSON format songs
type JSONLoader struct {
	parser *parser.JSONParser
	fs     embed.FS
}

// NewJSONLoader creates a new JSON song loader
func NewJSONLoader(fs embed.FS) *JSONLoader {
	return &JSONLoader{
		parser: parser.NewJSONParser(),
		fs:     fs,
	}
}

// LoadSong loads a song from JSON format
func (jl *JSONLoader) LoadSong(folderName string) (*types.Song, error) {
	logger.Info("Loading JSON song %s", folderName)
	songPath := path.Join(songDir, folderName)

	// Look for JSON file
	jsonPath := path.Join(songPath, "song.json")
	jsonData, err := jl.fs.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	// Parse JSON to internal format
	song, err := jl.parser.ParseSongData(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON song: %w", err)
	}

	// Load audio file path
	audioPath := path.Join(songPath, song.AudioPath)
	song.AudioPath = audioPath

	// Load artwork
	artPath := path.Join(songPath, "art.png")
	artData, err := jl.fs.ReadFile(artPath)
	if err != nil {
		logger.Debug("No artwork found for %s: %v", folderName, err)
	} else {
		artImage, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(artData))
		if err != nil {
			logger.Debug("Failed to load artwork for %s: %v", folderName, err)
		} else {
			song.Art = artImage
		}
	}

	// Validate the song has at least one chart
	if len(song.Charts) == 0 {
		return nil, fmt.Errorf("song %s has no charts", folderName)
	}

	logger.Debug("Successfully loaded song %s with %d charts", song.Title, len(song.Charts))
	return song, nil
}

// CanLoad checks if this loader can handle the directory
func (jl *JSONLoader) CanLoad(folderName string) bool {
	songPath := path.Join(songDir, folderName)
	jsonPath := path.Join(songPath, "song.json")
	
	// Check if song.json exists
	_, err := jl.fs.ReadFile(jsonPath)
	return err == nil
}

// Global loader instance
var globalLoader *JSONLoader

// InitLoaders initializes the song loading system
func InitLoaders() {
	globalLoader = NewJSONLoader(songFS)
	logger.Debug("Initialized JSON song loader")
}

// LoadSongJSON loads a song using the JSON loader
func LoadSongJSON(folderName string) (*types.Song, error) {
	if globalLoader == nil {
		InitLoaders()
	}

	return globalLoader.LoadSong(folderName)
}

// CanLoadSong checks if a song can be loaded
func CanLoadSong(folderName string) bool {
	if globalLoader == nil {
		InitLoaders()
	}

	return globalLoader.CanLoad(folderName)
}