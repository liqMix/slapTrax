package assets

import (
	"crypto/sha256"
	"embed"
	"encoding/binary"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/liqmix/slaptrax/internal/logger"
)

//go:embed songs/**/*.yaml songs/**/*.png songs/**/*.mid
var songFS embed.FS

const songDir = "songs"
const songAudioPrefix = "audio."
const songMetaFilename = "meta.yaml"
const songArtPath = "art.png"
const defaultArtPath = "default_art.png"

var songs map[string]*SongData = make(map[string]*SongData)

func InitSongs() {
	songDirs := readSongDir()
	for _, songDir := range songDirs {
		// Load the song from the directory
		song := loadSong(songDir)
		if song != nil {
			songs[song.GetHash()] = song
		}
	}
}

func GetAllSongData() []*SongData {
	songData := make([]*SongData, 0, len(songs))
	for _, song := range songs {
		songData = append(songData, song)
	}
	return songData
}

func GetSongData(hash string) *SongData {
	song, ok := songs[hash]
	if !ok {
		return nil
	}
	return song
}

type SongData struct {
	hash       string
	FolderName string
	Meta       []byte
	Art        *ebiten.Image
	AudioPath  string
	Charts     map[int][]byte
}

func (sd *SongData) GetHash() string {
	if sd.hash != "" {
		return sd.hash
	}
	hasher := sha256.New()

	// Hash the metadata
	hasher.Write(sd.Meta)

	// Hash chart data in a deterministic order
	difficulties := make([]int, 0, len(sd.Charts))
	for diff := range sd.Charts {
		difficulties = append(difficulties, diff)
	}
	sort.Ints(difficulties)

	for _, diff := range difficulties {
		// Include difficulty level in hash to ensure unique hashes even if chart data is identical
		diffBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(diffBytes, uint32(diff))
		hasher.Write(diffBytes)

		hasher.Write(sd.Charts[diff])
	}

	// Store the computed hash
	sd.hash = fmt.Sprintf("%x", hasher.Sum(nil))
	return sd.hash
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

func findAudioFile(folderName string) (string, error) {
	dirPath := path.Join(songDir, folderName)
	entries, err := audioFS.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			logger.Debug("\tSkipping directory %s", entry.Name())
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, songAudioPrefix) {
			logger.Debug("\tSkipping file %s", name)
			continue
		}

		for _, audioExt := range []AudioExt{Wav, Ogg, Mp3} {
			if audioExt.Is(name) {
				return path.Join(dirPath, name), nil
			}
		}
	}

	return "", fmt.Errorf("no audio file found in %s", folderName)
}

func loadSong(folderName string) *SongData {
	logger.Info("Loading song %s", folderName)
	song := SongData{}

	// Load the song from the directory
	songPath := path.Join(songDir, folderName)
	metaFile, err := songFS.ReadFile(path.Join(songPath, songMetaFilename))
	if err != nil {
		logger.Error("\tError reading song metadata file: %s", err)
		return nil
	}
	song.Meta = metaFile

	// Load the album art
	artPath := path.Join(songPath, songArtPath)
	artImg, _, err := ebitenutil.NewImageFromFileSystem(songFS, artPath)
	if err != nil {
		logger.Warn("\tUnable to find art for %s, falling back to default", folderName)
		song.Art = GetImage(defaultArtPath)
	} else {
		song.Art = artImg
	}

	// Identify the audio file
	audioPath, err := findAudioFile(folderName)
	if err != nil {
		logger.Error("\tError finding audio file: %s", err)
		return nil
	}
	song.AudioPath = audioPath

	// Load the charts
	charts := make(map[int][]byte)

	// Find all  files in the song directory
	// and load them as charts with the chart name as the difficulty.
	// If filename isn't an integer, ignore it.
	chartDir, err := songFS.ReadDir(songPath)
	if err != nil {
		logger.Error("\tUnable to read song directory: %w", err)
		return nil
	}

	for _, entry := range chartDir {
		name := entry.Name()
		if !entry.IsDir() && filepath.Ext(name) == ".mid" {
			d, err := strconv.Atoi(name[:len(name)-4])
			if err != nil {
				logger.Warn("\tInvalid chart file name: %s", name)
				continue
			}

			chartPath := path.Join(songPath, name)
			chartFile, err := songFS.ReadFile(chartPath)
			if err != nil {
				logger.Error("Error reading chart file: %s", err)
				continue
			}
			charts[d] = chartFile

		}
	}

	if len(charts) == 0 {
		logger.Error("\tNo charts found in song directory: %s", folderName)
		return nil
	}

	song.Charts = charts
	return &song
}
