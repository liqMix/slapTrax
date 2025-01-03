package state

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/locale"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/resource"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

type SongDetails struct {
	// Clickable to navigate to song meta links
	title   *ui.Button
	artist  *ui.Button
	album   *ui.Button
	charter *ui.Button

	// Display only
	art     *ui.Element
	bpm     *ui.Element
	version *ui.Element
	year    *ui.Element

	chartText      *ui.Element
	difficultyText *ui.Element
	difficulties   []*ui.Element

	// Positions
	panelSize        *ui.Point
	difficultyCenter *ui.Point
}

var artSize = 0.2
var detailsTop = ui.Point{X: 0.25, Y: 0.27}
var detailsCenter = ui.Point{X: 0.25, Y: 0.50}

func NewSongDetails() *SongDetails {
	d := &SongDetails{}
	center := ui.Point{
		X: detailsTop.X,
		Y: detailsTop.Y,
	}
	offset := float64(ui.TextHeight()) * 1.2

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
	e.SetText(locale.String(types.L_DIFFICULTIES))
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
	e.SetText(locale.String(types.L_CHART))
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
	s.album.SetText(song.Album)
	s.charter.SetText(song.ChartedBy)

	links := song.GetSongLinks()
	if links != nil {
		s.title.SetTrigger(func() {
			resource.OpenURL(song.TitleLink)
		})

		s.artist.SetTrigger(func() {
			resource.OpenURL(song.ArtistLink)
		})

		s.album.SetTrigger(func() {
			resource.OpenURL(song.AlbumLink)
		})
		if links.CharterLink != "" {
			s.charter.SetTrigger(func() {
				resource.OpenURL(links.CharterLink)
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
		totalWidth += float64(ui.TextWidth(diffS))
	}
	// Add spacing between elements
	totalWidth += spacing * float64(len(difficulties)-1)

	// Start from leftmost position
	center.X = center.X - (totalWidth / 2)

	for _, diff := range difficulties {
		d := ui.NewElement()
		diffWidth := float64(ui.TextWidth(diff.String()))
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
	ui.DrawFilledRect(screen, &detailsCenter, s.panelSize, types.Gray)
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

type SongSelection struct {
	types.BaseGameState

	songs      []*types.Song
	details    *SongDetails
	currentIdx int
	songList   *ui.UIGroup
}

func (s *SongSelection) SelectSong(song *types.Song) {
	if song == nil {
		logger.Debug("Invalid song selection - song is nil")
		return
	}

	s.SetNextState(types.GameStateDifficultySelection, &DifficultySelectionArgs{
		song: song,
	})
}

func NewSongSelectionState() *SongSelection {
	songs := resource.GetAllSongs()
	s := &SongSelection{songs: songs}

	center := ui.Point{
		X: 0.75,
		Y: 0.25,
	}
	offset := float64(ui.TextHeight() * 2)

	songElements := ui.NewUIGroup()
	for _, song := range songs {
		e := ui.NewElement()
		e.SetCenter(center)
		e.SetText(song.Title)
		e.SetTrigger(func() {
			s.SelectSong(song)
		})
		songElements.Add(e)
		center.Y += offset
	}
	s.songList = songElements
	s.currentIdx = 0

	details := NewSongDetails()
	s.details = details
	if len(songs) > 0 {
		song := songs[0]
		s.details.UpdateDetails(song)
		audio.PlaySongPreview(song)
	}
	return s
}

func (s *SongSelection) Update() error {
	if input.K.Is(ebiten.KeyEscape, input.JustPressed) {
		audio.StopAll()
		s.SetNextState(types.GameStateTitle, nil)
	}

	s.songList.Update()
	currentIdx := s.songList.GetIndex()
	if currentIdx != s.currentIdx {
		s.currentIdx = currentIdx
		song := s.songs[s.currentIdx]
		s.details.UpdateDetails(song)
		audio.PlaySongPreview(song)
	}
	s.details.Update()
	return nil
}

func (s *SongSelection) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.songList.Draw(screen, opts)
	s.details.Draw(screen, opts)
}
