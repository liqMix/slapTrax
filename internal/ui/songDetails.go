package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

var (
	detailsTop    = Point{X: 0.5, Y: 0.27}
	detailsCenter = Point{X: 0.5, Y: 0.50}
	artSize       = 0.2
)

type SongDetails struct {
	// Clickable to navigate to song meta links
	title   *HotLink
	artist  *HotLink
	album   *HotLink
	charter *HotLink

	// Display only
	art     *Element
	bpm     *Element
	version *Element
	year    *Element

	chartText      *Element
	difficultyText *Element
	difficulties   []*Element

	// Positions
	panelSize        *Point
	difficultyCenter *Point
}

func NewSongDetails() *SongDetails {
	d := &SongDetails{}
	offset := float64(TextHeight(nil)) * 1.3
	size := Point{X: 0.15, Y: 0.1}
	center := Point{
		X: detailsTop.X,
		Y: detailsTop.Y + offset/2,
	}

	// Art
	e := NewElement()
	e.SetCenter(center)
	e.SetSize(size)
	d.art = e
	center.Y += size.Y + offset*2

	// Track Details
	size = Point{X: 0.1, Y: 0.025}
	b := NewHotLink()
	b.SetCenter(center)
	b.SetSize(size)
	b.SetTextBold(true)
	b.SetTextScale(1.5)
	d.title = b
	center.Y += offset * 1.2

	b = NewHotLink()
	b.SetCenter(center)
	b.SetSize(size.Scale(1.3))
	d.artist = b
	center.Y += offset

	b = NewHotLink()
	b.SetCenter(center)
	b.SetSize(size.Scale(1.3))
	d.album = b
	center.Y += offset * 2

	size = Point{X: 0.15, Y: 0.1}
	e = NewElement()
	e.SetSize(size)
	e.SetCenter(center)
	d.year = e
	center.Y += offset

	d.bpm = NewElement()
	d.bpm.SetCenter(center)
	center.Y += offset * 2

	// Difficulties
	// e = NewElement()
	// e.SetCenter(center)
	// e.SetText(l.String(l.DIFFICULTIES))
	// e.SetSize(size.Scale(1.2))
	// e.SetTextBold(true)
	// d.difficultyText = e
	// center.Y += offset * e.GetSize().Y

	d.difficultyCenter = &Point{
		X: center.X,
		Y: center.Y,
	}
	center.Y += offset * 2

	// Chart Details
	e = NewElement()
	e.SetCenter(center)
	e.SetText(l.String(l.CHART))
	e.SetSize(size)
	e.SetTextBold(true)
	d.chartText = e
	// center.Y += offset * 2

	e = NewElement()
	e.SetCenter(center)
	e.SetSize(size)
	center.Y += offset
	d.version = e

	b = NewHotLink()
	b.SetCenter(center)
	b.SetSize(size)
	// center.Y += offset * 2
	d.charter = b

	center.Y += offset * 2

	height := center.Y - detailsTop.Y + offset
	d.panelSize = &Point{
		X: artSize,
		Y: height,
	}
	return d
}

func (s *SongDetails) UpdateDetails(song *types.Song, difficulty types.Difficulty) {
	if song == nil {
		logger.Debug("Invalid song details update - song is nil")
		return
	}
	s.title.SetText(song.Title)
	s.artist.SetText(song.Artist)
	album := song.Album
	if album == "" {
		album = " "
	}
	s.album.SetText(album)
	s.charter.SetText(song.ChartedBy)

	links := song.GetSongLinks()
	if links != nil {
		s.title.SetURL(links.TitleLink)
		s.artist.SetURL(links.ArtistLink)
		s.album.SetURL(links.AlbumLink)

		if links.CharterLink != "" {
			s.charter.SetURL(links.CharterLink)
		} else {
			s.charter.SetURL("")
		}
	} else {
		s.title.SetURL("")
		s.artist.SetURL("")
		s.album.SetURL("")
		s.charter.SetURL("")
	}

	s.art.SetImage(song.Art)
	s.bpm.SetText(fmt.Sprintf("%d BPM", song.BPM))
	s.version.SetText(song.Version)
	s.year.SetText(fmt.Sprintf("%d", song.Year))

	difficulties := song.GetDifficulties()
	s.difficulties = make([]*Element, 0, len(difficulties))
	spacing := 0.025
	center := Point{
		X: s.difficultyCenter.X,
		Y: s.difficultyCenter.Y,
	}

	totalWidth := 0.0

	// Calculate total width of all text
	for _, d := range difficulties {
		diffS := d.String()
		totalWidth += float64(TextWidth(nil, diffS) * 1.25)
	}
	// Add spacing between elements
	totalWidth += spacing * float64(len(difficulties)-1)

	// Start from leftmost position
	center.X = center.X - (totalWidth / 1.95)

	for _, diff := range difficulties {
		d := NewElement()
		diffWidth := float64(TextWidth(nil, diff.String()))
		center.X += diffWidth / 2

		// Center each element at its position
		d.SetCenter(center)
		if diff == difficulty {
			d.SetTextScale(1.25)
			d.SetTextBold(true)
			d.SetTextColor(diff.Color())
		} else {
			d.SetTextScale(0.8)
			d.SetTextColor(types.Gray.C())
		}
		d.SetText(diff.String())
		s.difficulties = append(s.difficulties, d)

		// Move to next position
		center.X += diffWidth + spacing
	}
}

func (s *SongDetails) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	DrawNoteThemedRect(screen, &detailsCenter, s.panelSize)
	s.art.Draw(screen, opts)
	s.title.Draw(screen, opts)
	s.artist.Draw(screen, opts)
	s.album.Draw(screen, opts)
	s.year.Draw(screen, opts)
	s.bpm.Draw(screen, opts)
	// s.chartText.Draw(screen, opts)
	s.version.Draw(screen, opts)
	// s.charter.Draw(screen, opts)
	// s.difficultyText.Draw(screen, opts)
	for _, d := range s.difficulties {
		d.Draw(screen, opts)
	}
}

func (s *SongDetails) Update() {
	s.title.Update()
	s.artist.Update()
	s.album.Update()
	// s.charter.Update()
}
