package assets

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"gopkg.in/yaml.v2"
)

//go:embed songs/**/*.yaml songs/**/*.png songs/**/*.mid  songs/**/*.mp3
var songFS embed.FS

var songs = make(map[types.Checksum]types.Song)

func InitSongs() {
	// Read the song directory for valid songs
	songDirs := readSongDir()
	for _, songDir := range songDirs {
		// Load the song from the directory
		song := loadSong(songDir)
		if song != nil {
			songs[song.Checksum] = *song
		}
	}
}

func readSongDir() []string {
	songDir, err := songFS.ReadDir(config.SONG_DIR)
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

func findAudioFile(folderName string) (string, error) {
	dirPath := path.Join(config.SONG_DIR, folderName)
	entries, err := songFS.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	validExts := map[string]bool{
		".wav": true,
		".ogg": true,
		".mp3": true,
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, config.SONG_AUDIO_NAME+".") {
			continue
		}

		ext := filepath.Ext(name)
		if validExts[ext] {
			return path.Join(dirPath, name), nil
		}
	}

	return "", fmt.Errorf("no valid audio file found")
}

func loadSong(folderName string) *types.Song {
	// Load the song from the directory
	songPath := path.Join(config.SONG_DIR, folderName)
	songFile, err := songFS.ReadFile(path.Join(songPath, config.SONG_META_NAME))
	if err != nil {
		fmt.Println("Error reading song metadata file: " + err.Error())
		return nil
	}

	var song types.Song
	if err := yaml.Unmarshal(songFile, &song); err != nil {
		fmt.Println("Error unmarshalling song metadata: " + err.Error())
		return nil
	}

	// Load the album art
	artPath := path.Join(songPath, "art.png")
	artImg, _, err := ebitenutil.NewImageFromFileSystem(songFS, artPath)
	if err != nil {
		fmt.Println("Error loading album art: " + err.Error())
	} else {
		song.Art = artImg
	}

	// Identify the audio file
	audioPath, err := findAudioFile(folderName)
	if err != nil {
		logger.Error("finding audio file for %s: %w", song.Title, err)
		return nil
	}
	song.AudioPath = audioPath

	// Load the charts
	song.Charts = make(map[types.Difficulty]*types.Chart)

	// Find all  files in the song directory
	// and load them as charts with the chart name as the difficulty.
	// If not an integer, ignore it.
	chartDir, err := songFS.ReadDir(songPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, entry := range chartDir {
		name := entry.Name()
		if !entry.IsDir() && filepath.Ext(name) == ".mid" {
			// check if integer directly
			fmt.Println(name)
			d, err := strconv.Atoi(name[:len(name)-4])
			if err != nil {
				fmt.Println("Invalid chart file name: " + name)
				continue
			}

			difficulty := types.Difficulty(d)
			chartPath := path.Join(songPath, name)
			chartFile, err := songFS.ReadFile(chartPath)
			if err != nil {
				fmt.Println("Error reading chart file: " + err.Error())
				continue
			}

			chart, err := ParseChart(&song, chartFile)
			if err != nil || chart == nil {
				fmt.Println("Error parsing chart file: " + err.Error())
				continue
			}
			for _, track := range chart.Tracks {
				fmt.Println(track.Name, " has ", len(track.AllNotes), " notes")
			}

			chart.Difficulty = difficulty
			song.Charts[difficulty] = chart
		}
	}

	if len(song.Charts) == 0 {
		fmt.Println("No valid charts found for song " + song.Title)
		return nil
	}

	// Calculate the checksum
	song.Checksum, err = calculateChecksum(songPath)
	if err != nil {
		fmt.Println("Error calculating checksum: " + err.Error())
		return nil
	}
	song.FolderName = folderName
	return &song
}

func calculateChecksum(songPath string) (types.Checksum, error) {
	hasher := sha256.New()

	// Get all files recursively from embed.FS
	var files []string
	err := fs.WalkDir(songFS, songPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("walking directory: %w", err)
	}

	// Sort for consistent ordering
	sort.Strings(files)

	// Hash each file's contents
	for _, file := range files {
		data, err := songFS.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("reading file %s: %w", file, err)
		}

		// Write both filename and content to ensure unique hashes
		// even if files are renamed
		hasher.Write([]byte(path.Base(file)))
		hasher.Write(data)
	}

	return types.Checksum(fmt.Sprintf("%x", hasher.Sum(nil))), nil
}

// Returns a song by checksum
func Get(checksum types.Checksum) *types.Song {
	song, ok := songs[checksum]
	if !ok {
		return nil
	}
	return &song
}

// Returns all songs
func GetAllSongs() []*types.Song {
	var songList []*types.Song
	for _, song := range songs {
		songList = append(songList, &song)
	}

	// Sort by title
	sort.Slice(songList, func(i, j int) bool {
		return songList[i].Title < songList[j].Title
	})
	return songList
}

func GetSongByTitle(title string) *types.Song {
	for _, song := range songs {
		if song.Title == title {
			return &song
		}
	}
	return nil
}
