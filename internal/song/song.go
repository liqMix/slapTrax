package song

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"gopkg.in/yaml.v2"
)

type Checksum string

type Song struct {
	// metadata read from file
	Title     string `yaml:"title"`
	TitleLink string `yaml:"title_link"` // clickable external link to the song

	Artist     string `yaml:"artist"`
	ArtistLink string `yaml:"artist_link"` // clickable external link to the artist

	Album     string `yaml:"album"`      // album the song is from
	AlbumLink string `yaml:"album_link"` // clickable external link to the album

	Year         int `yaml:"year"`          // year the song was released
	BPM          int `yaml:"bpm"`           // beats per minute of the song
	Length       int `yaml:"length"`        // length of the song in milliseconds (maybe can be derived from the audio)
	PreviewStart int `yaml:"preview_start"` // start of the preview in milliseconds

	ChartedBy     string `yaml:"charted_by"`      // name of the person who made the chart
	ChartedByLink string `yaml:"charted_by_link"` // clickable external link to the charter
	Version       string `yaml:"version"`
	////

	Art       ebiten.Image // album art
	AudioPath string       // path to the audio file
	Checksum  Checksum     // calculated from folder contents
	Charts    []Chart      // charts for the song

	folderName string
}

var songs = make(map[Checksum]Song)

// Ensure we have the files we need in the song directory
func validateSongDir(folderName string) bool {
	dirName := config.SONG_DIR + "/" + folderName
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return false
	}

	// metadata file (song.yaml)
	_, err = os.Stat(dirName + "/song.yaml")
	if os.IsNotExist(err) {
		return false
	}

	// audio file (audio.wav/.ogg/.mp3)
	_, err = os.Stat(dirName + "/audio.wav")
	if os.IsNotExist(err) {
		_, err = os.Stat(dirName + "/audio.ogg")
		if os.IsNotExist(err) {
			_, err = os.Stat(dirName + "/audio.mp3")
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

func loadSong(folderName string) *Song {
	// Load the song from the directory
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
	song.Charts = make([]Chart, 0)
	// this'll be fun

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
