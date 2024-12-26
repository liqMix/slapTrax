package song

import (
	"github.com/hajimehoshi/ebiten/v2"
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

	Year         int   `yaml:"year"`          // year the song was released
	BPM          int   `yaml:"bpm"`           // beats per minute of the song
	Length       int   `yaml:"length"`        // length of the song in milliseconds (maybe can be derived from the audio)
	PreviewStart int64 `yaml:"preview_start"` // start of the preview in milliseconds

	ChartedBy     string `yaml:"charted_by"`      // name of the person who made the chart
	ChartedByLink string `yaml:"charted_by_link"` // clickable external link to the charter
	Version       string `yaml:"version"`
	////

	Art       ebiten.Image         // album art
	AudioPath string               // path to the audio file
	Checksum  Checksum             // calculated from folder contents
	Charts    map[Difficulty]Chart // charts for the song

	folderName string
}

func (s *Song) String() string {
	return s.Title + " - " + s.Artist
}

func (s *Song) GetChart(difficulty Difficulty) Chart {
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
	return difficulties
}
