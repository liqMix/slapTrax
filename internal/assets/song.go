package assets

import (
	"embed"

	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types"
)

//go:embed songs/**/*.json songs/**/*.png songs/**/*.ogg
var songFS embed.FS

const songDir = "songs"
const songArtPath = "art.png"
const defaultArtPath = "default_art.png"
const songAudioFile = "audio.ogg"

var songs map[string]*types.SongData = make(map[string]*types.SongData)
var loadedSongs map[string]*types.Song = make(map[string]*types.Song)

func InitSongs() {
	// Initialize loaders
	InitLoaders()

	songDirs := readSongDir()
	for _, songDir := range songDirs {
		// Load JSON song
		song, err := LoadSongJSON(songDir)
		if err != nil {
			logger.Warn("Failed to load song %s: %v", songDir, err)
			continue
		}

		// Store loaded song
		loadedSongs[song.Hash] = song
		logger.Debug("Successfully loaded song %s (JSON format) with hash %s", songDir, song.Hash)
	}
}

func GetAllSongData() []*types.SongData {
	// Since we've moved to JSON loading, create SongData from loaded songs
	songData := make([]*types.SongData, 0, len(loadedSongs))
	for _, song := range loadedSongs {
		// Convert loaded Song back to SongData format for compatibility
		data := &types.SongData{
			FolderName: song.FolderName,
			Meta:       []byte{}, // JSON is parsed already
			Art:        song.Art,
			AudioPath:  song.AudioPath,
			Charts:     make(map[int][]byte),
		}
		
		// Create chart data for each difficulty
		for diff := range song.Charts {
			// Since charts are already parsed, we don't need the raw bytes
			// But for compatibility, we'll create empty byte slices
			data.Charts[int(diff)] = []byte{}
		}
		
		// Hash will be computed automatically
		songData = append(songData, data)
	}
	return songData
}

func GetSongData(hash string) *types.SongData {
	song, ok := songs[hash]
	if !ok {
		return nil
	}
	return song
}

// GetAllLoadedSongs returns all songs loaded with the new system
func GetAllLoadedSongs() []*types.Song {
	logger.Debug("GetAllLoadedSongs called, map has %d songs", len(loadedSongs))
	songs := make([]*types.Song, 0, len(loadedSongs))
	for hash, song := range loadedSongs {
		logger.Debug("Found song in map: '%s' with hash %s", song.Title, hash)
		songs = append(songs, song)
	}
	return songs
}

// GetLoadedSong returns a song loaded with the new system
func GetLoadedSong(hash string) *types.Song {
	song, ok := loadedSongs[hash]
	if !ok {
		return nil
	}
	return song
}

// GetSongByFolder returns a song by its folder name
func GetSongByFolder(folderName string) *types.Song {
	for _, song := range loadedSongs {
		if song.FolderName == folderName {
			return song
		}
	}
	return nil
}


func readSongDir() []string {
	songDir, err := songFS.ReadDir(songDir)
	if err != nil {
		panic(err)
	}

	songs := make([]string, 0)
	for _, entry := range songDir {
		if entry.IsDir() {
			songs = append(songs, entry.Name())
		}
	}
	return songs
}

