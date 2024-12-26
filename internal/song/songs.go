package song

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"

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

// Ensure we have the files we need in the song directory
func validateSongDir(folderName string) bool {
	dirName := path.Join(config.SONG_DIR, folderName)
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return false
	}

	// metadata file (song.yaml)
	metaPath := path.Join(dirName, config.SONG_META_NAME)
	_, err = os.Stat(metaPath)
	if os.IsNotExist(err) {
		return false
	}

	// audio file (audio.wav/.ogg/.mp3)
	songPath := path.Join(dirName, config.SONG_AUDIO_NAME)
	_, err = os.Stat(songPath + ".wav")
	if os.IsNotExist(err) {
		_, err = os.Stat(songPath + ".ogg")
		if os.IsNotExist(err) {
			_, err = os.Stat(songPath + ".mp3")
			if os.IsNotExist(err) {
				return false
			}
		}
	}

	// at least one chart file (easy.yaml/hard.yaml)
	_, err = os.Stat(dirName + "/easy.yaml")
	if os.IsNotExist(err) {
		_, err = os.Stat(dirName + "/hard.yaml")
		if os.IsNotExist(err) {
			return false
		}
	}

	// maybe no art is ok ?
	// _, err = os.Stat(dirName + "/art.png")
	// if os.IsNotExist(err) {
	// 	return false
	// }
	return true
}

func readSongDir() []string {
	songDir, err := os.ReadDir(config.SONG_DIR)
	if err != nil {
		panic(err)
	}

	songs := make([]string, 0)
	for _, entry := range songDir {
		if entry.IsDir() {
			if validateSongDir(entry.Name()) {
				songs = append(songs, entry.Name())
			}
		}
	}
	return songs
}

func loadSong(folderName string) *Song {
	// Load the song from the directory
	songPath := config.SONG_DIR + "/" + folderName
	songFile, err := os.ReadFile(songPath + "/song.yaml")
	if err != nil {
		return nil
	}

	var song Song
	if err := yaml.Unmarshal(songFile, &song); err != nil {
		return nil
	}

	// Load the album art
	artPath := songPath + "/art.png"
	artImg, _, err := ebitenutil.NewImageFromFile(artPath)
	if err != nil {
		return nil
	}
	song.Art = *artImg

	// Identify the audio file
	audioPath := songPath + "/audio.wav"
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		audioPath = songPath + "/audio.ogg"
		if _, err := os.Stat(audioPath); os.IsNotExist(err) {
			audioPath = songPath + "/audio.mp3"
			if _, err := os.Stat(audioPath); os.IsNotExist(err) {
				panic("No audio file found for song " + song.Title)
			}
		}
	}
	song.AudioPath = audioPath

	// Load the charts
	song.Charts = make(map[Difficulty]Chart)
	easyChartPath := songPath + "/easy.yaml"
	easyChartFile, err := os.ReadFile(easyChartPath)
	if err == nil {
		easyChart := ParseChart(Easy, easyChartFile)
		song.Charts[Easy] = *easyChart
	} else {
		fmt.Println("No easy chart found for song " + song.Title)
	}

	hardChartPath := songPath + "/hard.yaml"
	hardChartFile, err := os.ReadFile(hardChartPath)
	if err == nil {
		hardChart := ParseChart(Hard, hardChartFile)
		song.Charts[Hard] = *hardChart
		return nil
	} else {
		fmt.Println("No hard chart found for song " + song.Title)
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
