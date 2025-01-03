package types

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type Checksum string

type Chart struct {
	Difficulty Difficulty
	TotalNotes int
	Tracks     []*Track
}
type Track struct {
	Name        TrackName
	AllNotes    []*Note
	ActiveNotes []*Note

	Active      bool
	StaleActive bool

	NextNoteIndex int
}

func NewTrack(name TrackName, notes []*Note, beatInterval int64) *Track {
	// Reset the notes
	for _, n := range notes {
		n.Reset()
	}

	// Sort the notes by target time
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Target < notes[j].Target
	})

	return &Track{
		Name:     name,
		AllNotes: notes,
	}
}

func (t *Track) Reset() {
	t.ActiveNotes = make([]*Note, 0)
	t.Active = false
	t.StaleActive = false
	t.NextNoteIndex = 0
	for _, n := range t.AllNotes {
		n.Reset()
	}
}

func (t Track) IsPressed() bool {
	return t.Active || t.StaleActive
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
	Checksum  Checksum              // calculated from folder contents
	Charts    map[Difficulty]*Chart // charts for the song

	FolderName string
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
	defaultLink := ""

	artistLink := s.ArtistLink
	hasArtistLink := artistLink != ""

	albumLink := s.AlbumLink
	hasAlbumLink := albumLink != ""

	titleLink := s.TitleLink
	hasTitleLink := titleLink != ""

	if hasArtistLink {
		defaultLink = artistLink
	} else if hasAlbumLink {
		defaultLink = albumLink
	} else if hasTitleLink {
		defaultLink = titleLink
	}

	if defaultLink == "" {
		return nil
	}

	if !hasArtistLink {
		artistLink = defaultLink
	}
	if !hasAlbumLink {
		albumLink = defaultLink
	}
	if !hasTitleLink {
		titleLink = defaultLink
	}

	return &SongLinks{
		ArtistLink:  artistLink,
		AlbumLink:   albumLink,
		TitleLink:   titleLink,
		CharterLink: s.ChartedByLink,
	}
}
