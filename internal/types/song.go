package types

import (
	"fmt"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/logger"
	"gopkg.in/yaml.v2"
)

func GetAllSongs(allSongData []*SongData) []*Song {
	songs := make([]*Song, 0, len(allSongData))
	for _, songData := range allSongData {
		song, err := NewSong(songData)
		if err != nil {
			continue
		}
		songs = append(songs, song)
	}
	return songs
}

func GetAllCharts(allSongData []*SongData) []*Chart {
	charts := make([]*Chart, 0)
	for _, songData := range allSongData {
		song, err := NewSong(songData)
		if err != nil {
			logger.Error("Failed to parse song data", err)
			continue
		}
		for _, chart := range song.Charts {
			charts = append(charts, chart)
		}
	}
	return charts
}

type SongLinks struct {
	ArtistLink  string
	AlbumLink   string
	TitleLink   string
	CharterLink string
}

type Song struct {
	// metadata read from file
	Title     string `yaml:"title"`
	TitleLink string `yaml:"title_link"` // clickable external link to the song

	Artist     string `yaml:"artist"`
	ArtistLink string `yaml:"artist_link"` // clickable external link to the artist

	Album     string `yaml:"album"`      // album the song is from
	AlbumLink string `yaml:"album_link"` // clickable external link to the album

	Year         int   `yaml:"year"`          // year the song was released
	BPM          int   `yaml:"bpm"`           // beats per minute of the song
	Length       int   `yaml:"length"`        // length of the song in milliseconds (maybe can be derived from the audio)
	PreviewStart int64 `yaml:"preview_start"` // start of the preview in milliseconds

	ChartedBy     string `yaml:"charted_by"`      // name of the person who made the chart
	ChartedByLink string `yaml:"charted_by_link"` // clickable external link to the charter
	Version       string `yaml:"version"`
	////

	Art       *ebiten.Image         // album art
	AudioPath string                // path to the audio file
	Charts    map[Difficulty]*Chart // charts for the song

	FolderName string
	Hash       string
}

var parsedSongs map[string]*Song = make(map[string]*Song)

func NewSong(songData *SongData) (*Song, error) {
	if songData == nil {
		return nil, fmt.Errorf("songData is nil")
	}
	hash := songData.GetHash()
	if song, ok := parsedSongs[hash]; ok {
		return song, nil
	}

	meta := songData.Meta
	song := &Song{
		Hash: songData.GetHash(),
	}
	err := yaml.Unmarshal(meta, song)
	if err != nil {
		return nil, err
	}

	song.Art = songData.Art
	song.AudioPath = songData.AudioPath
	song.Charts = make(map[Difficulty]*Chart)
	for difficulty, chartData := range songData.Charts {
		chart, err := NewChart(song, chartData)
		if err != nil {
			return nil, err
		}

		song.Charts[Difficulty(difficulty)] = chart
	}

	parsedSongs[hash] = song
	return song, nil
}

func (s *Song) GetChart(difficulty Difficulty) *Chart {
	chart, ok := s.Charts[difficulty]
	if !ok {
		panic("No chart for difficulty")
	}
	return chart
}

func (s *Song) GetDifficulties() []Difficulty {
	difficulties := make([]Difficulty, 0, len(s.Charts))
	for difficulty := range s.Charts {
		difficulties = append(difficulties, difficulty)
	}
	sort.Slice(difficulties, func(i, j int) bool {
		return difficulties[i] < difficulties[j]
	})
	return difficulties
}

func (s *Song) GetSongLinks() *SongLinks {
	if s == nil {
		return nil
	}
	artistLink := s.ArtistLink
	albumLink := s.AlbumLink
	titleLink := s.TitleLink
	return &SongLinks{
		ArtistLink:  artistLink,
		AlbumLink:   albumLink,
		TitleLink:   titleLink,
		CharterLink: s.ChartedByLink,
	}
}

func (s *Song) GetBeatInterval() int64 {
	return int64(60000 / s.BPM)
}

func (s *Song) GetCountdownTicks(restart bool) []int64 {
	b := s.GetBeatInterval()

	if restart {
		return []int64{
			-b * 2,
			-b * 1,
			0,
		}
	}
	return []int64{
		-b * 4,
		-b * 3,
		-b * 2,
		-b * 1,
		0,
	}
}
