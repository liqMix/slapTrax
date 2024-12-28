package song

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"gopkg.in/yaml.v2"
)

var songs = make(map[Checksum]Song)

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
	songDir, err := os.ReadDir(config.SONG_DIR)
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

// // Ensure we have the files we need in the song directory
// func validateSongDir(folderName string) bool {
// 	dirName := path.Join(config.SONG_DIR, folderName)
// 	_, err := os.Stat(dirName)
// 	if os.IsNotExist(err) {
// 		fmt.Println("Song directory does not exist: " + dirName)
// 		return false
// 	}

// 	// metadata file
// 	metaPath := path.Join(dirName, config.SONG_META_NAME)
// 	_, err = os.Stat(metaPath)
// 	if os.IsNotExist(err) {
// 		fmt.Println("Song metadata file does not exist: " + metaPath)
// 		return false
// 	}

// 	// audio file (audio.wav/.ogg/.mp3)
// 	songPath := path.Join(dirName, config.SONG_AUDIO_NAME)
// 	_, err = os.Stat(songPath + ".wav")
// 	if os.IsNotExist(err) {
// 		_, err = os.Stat(songPath + ".ogg")
// 		if os.IsNotExist(err) {
// 			_, err = os.Stat(songPath + ".mp3")
// 			if os.IsNotExist(err) {
// 				fmt.Println("Song audio file does not exist: " + songPath)
// 				return false
// 			}
// 		}
// 	}

// 	// at least one chart file (.midi with integer name)
// 	chartDir, err := os.ReadDir(dirName)
// 	if err != nil {
// 		return false
// 	}
// 	for _, entry := range chartDir {
// 		name := entry.Name()
// 		if !entry.IsDir() && filepath.Ext(name) == ".midi" {
// 			// check if integer directly
// 			_, err := strconv.Atoi(name[:len(name)-5])
// 			if err != nil {
// 				fmt.Println("Invalid chart file name: " + name)
// 				return false
// 			}
// 		}
// 	}

// 	// maybe no art is ok ?
// 	// _, err = os.Stat(dirName + "/art.png")
// 	// if os.IsNotExist(err) {
// 	// 	return false
// 	// }
// 	return true
// }

func loadSong(folderName string) *Song {
	// Load the song from the directory
	songPath := path.Join(config.SONG_DIR, folderName)
	songFile, err := os.ReadFile(path.Join(songPath, config.SONG_META_NAME))
	if err != nil {
		fmt.Println("Error reading song metadata file: " + err.Error())
		return nil
	}

	var song Song
	if err := yaml.Unmarshal(songFile, &song); err != nil {
		fmt.Println("Error unmarshalling song metadata: " + err.Error())
		return nil
	}

	// Load the album art
	artPath := songPath + "/art.png"
	artImg, _, err := ebitenutil.NewImageFromFile(artPath)
	if err != nil {
		fmt.Println("Error loading album art: " + err.Error())
	} else {
		song.Art = *artImg
	}

	// Identify the audio file
	audioPath := songPath + "/audio.wav"
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		audioPath = songPath + "/audio.ogg"
		if _, err := os.Stat(audioPath); os.IsNotExist(err) {
			audioPath = songPath + "/audio.mp3"
			if _, err := os.Stat(audioPath); os.IsNotExist(err) {
				fmt.Println("No audio file found for song " + song.Title)
				return nil
			}
		}
	}
	song.AudioPath = audioPath

	// Load the charts
	song.Charts = make(map[Difficulty]Chart)

	// Find all  files in the song directory
	// and load them as charts with the chart name as the difficulty.
	// If not an integer, ignore it.
	chartDir, err := os.ReadDir(songPath)
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

			difficulty := Difficulty(d)
			chartPath := path.Join(songPath, name)
			chartFile, err := os.ReadFile(chartPath)
			if err != nil {
				fmt.Println("Error reading chart file: " + err.Error())
				continue
			}

			chart, err := ParseChart(&song, chartFile)
			if err != nil {
				fmt.Println("Error parsing chart file: " + err.Error())
				continue
			}
			if chart != nil {
				chart.Difficulty = difficulty
				song.Charts[difficulty] = *chart
			}
		}
	}

	if len(song.Charts) == 0 {
		fmt.Println("No valid charts found for song " + song.Title)
		return nil
	}

	// Calculate the checksum
	song.Checksum = calculateChecksum(songPath)
	song.folderName = folderName
	return &song
}

func calculateChecksum(songPath string) Checksum {
	// Create hash
	hasher := sha256.New()

	// Get all files
	var files []string
	err := filepath.Walk(songPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return ""
	}

	// Sort for consistent ordering
	sort.Strings(files)

	// Hash each file's contents
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return ""
		}
		hasher.Write(data)
	}

	// Return hex string of hash
	return Checksum(fmt.Sprintf("%x", hasher.Sum(nil)))
}

// Returns a song by checksum
func Get(checksum Checksum) *Song {
	song, ok := songs[checksum]
	if !ok {
		return nil
	}
	return &song
}

// Returns all songs
func GetAll() []*Song {
	var songList []*Song
	for _, song := range songs {
		songList = append(songList, &song)
	}
	return songList
}

func GetTestSong() *Song {
	for _, song := range songs {
		fmt.Println(song.Title)
		if song.Title == "another" {
			return &song
		}
	}
	panic("No test song found")
}
