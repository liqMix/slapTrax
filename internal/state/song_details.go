package state

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/external"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

func NewSongDetails() *SongDetails {
	d := &SongDetails{}
	center := ui.Point{
		X: detailsTop.X,
		Y: detailsTop.Y,
	}
	offset := float64(ui.TextHeight(nil)) * 1.2

	// Art
	e := ui.NewElement()
	e.SetCenter(center)
	e.SetSize(ui.Point{X: artSize, Y: artSize})
	d.art = e
	center.Y += artSize/2 + offset

	// Track Details
	scale := 1.5
	b := ui.NewButton()
	b.SetCenter(center)
	b.SetScale(scale)
	b.SetTextBold(true)
	d.title = b
	center.Y += offset

	scale = 1.2
	b = ui.NewButton()
	b.SetCenter(center)
	b.SetScale(scale)
	d.artist = b
	center.Y += offset

	b = ui.NewButton()
	b.SetCenter(center)
	b.SetScale(scale)
	d.album = b
	center.Y += offset

	scale = 1.0
	e = ui.NewElement()
	e.SetCenter(center)
	d.year = e
	center.Y += offset

	d.bpm = ui.NewElement()
	d.bpm.SetCenter(center)
	center.Y += offset * 2

	// Difficulties
	scale = 1.2
	e = ui.NewElement()
	e.SetCenter(center)
	e.SetText(l.DIFFICULTIES)
	e.SetScale(scale)
	e.SetTextBold(true)
	d.difficultyText = e
	center.Y += offset * scale

	d.difficultyCenter = &ui.Point{
		X: center.X,
		Y: center.Y,
	}
	center.Y += offset * 2

	// Chart Details
	e = ui.NewElement()
	e.SetCenter(center)
	e.SetText(l.CHART)
	e.SetScale(scale)
	e.SetTextBold(true)
	d.chartText = e
	center.Y += offset

	e = ui.NewElement()
	e.SetCenter(center)
	center.Y += offset
	d.version = e

	b = ui.NewButton()
	b.SetCenter(center)
	center.Y += offset * 2
	d.charter = b

	center.Y += offset

	height := center.Y - detailsTop.Y + offset
	d.panelSize = &ui.Point{
		X: artSize,
		Y: height,
	}
	return d
}

func (s *SongDetails) UpdateDetails(song *types.Song) {
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
		s.title.SetTrigger(func() {
			external.OpenURL(song.TitleLink)
		})

		s.artist.SetTrigger(func() {
			external.OpenURL(song.ArtistLink)
		})

		s.album.SetTrigger(func() {
			external.OpenURL(song.AlbumLink)
		})
		if links.CharterLink != "" {
			s.charter.SetTrigger(func() {
				external.OpenURL(links.CharterLink)
			})
		} else {
			s.charter.SetTrigger(nil)
		}
	} else {
		s.title.SetTrigger(nil)
		s.artist.SetTrigger(nil)
		s.album.SetTrigger(nil)
		s.charter.SetTrigger(nil)
	}

	s.art.SetImage(song.Art)
	s.bpm.SetText(fmt.Sprintf("%d BPM", song.BPM))
	s.version.SetText(song.Version)
	s.year.SetText(fmt.Sprintf("%d", song.Year))

	difficulties := song.GetDifficulties()
	s.difficulties = make([]*ui.Element, 0, len(difficulties))
	spacing := 0.025
	center := ui.Point{
		X: s.difficultyCenter.X,
		Y: s.difficultyCenter.Y,
	}

	totalWidth := 0.0
	// Calculate total width of all text
	for _, d := range difficulties {
		diffS := d.String()
		totalWidth += float64(ui.TextWidth(nil, diffS))
	}
	// Add spacing between elements
	totalWidth += spacing * float64(len(difficulties)-1)

	// Start from leftmost position
	center.X = center.X - (totalWidth / 2)

	for _, diff := range difficulties {
		d := ui.NewElement()
		diffWidth := float64(ui.TextWidth(nil, diff.String()))
		center.X += diffWidth / 2
		// Center each element at its position
		d.SetCenter(center)
		d.SetText(diff.String())
		d.SetTextColor(diff.Color())
		s.difficulties = append(s.difficulties, d)

		// Move to next position
		center.X += diffWidth + spacing
	}
}

func (s *SongDetails) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	ui.DrawNoteThemedRect(screen, &detailsCenter, s.panelSize)
	s.art.Draw(screen, opts)
	s.title.Draw(screen, opts)
	s.artist.Draw(screen, opts)
	s.album.Draw(screen, opts)
	s.year.Draw(screen, opts)
	s.bpm.Draw(screen, opts)
	s.chartText.Draw(screen, opts)
	s.version.Draw(screen, opts)
	s.charter.Draw(screen, opts)
	s.difficultyText.Draw(screen, opts)
	for _, d := range s.difficulties {
		d.Draw(screen, opts)
	}
}
func (s *SongDetails) Update() {
	s.title.Update()
	s.artist.Update()
	s.album.Update()
	s.charter.Update()
}
