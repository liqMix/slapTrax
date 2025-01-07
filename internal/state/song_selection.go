package state

import (
	"sort"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
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
var detailsTop = ui.Point{X: 0.5, Y: 0.27}
var detailsCenter = ui.Point{X: 0.5, Y: 0.50}

type SongSelection struct {
	types.BaseGameState

	options  []*SongOption
	details  *SongDetails
	songList *ui.SongSelector

	currentIdx       int
	uiIdxToOptionIdx map[int]int
}

type SongOption struct {
	song       *types.Song
	difficulty types.Difficulty
}

func NewSongSelectionState() *SongSelection {
	allSongs := types.GetAllSongs()
	songs := make([]*SongOption, 0)
	for _, song := range allSongs {
		for _, diff := range song.GetDifficulties() {
			songs = append(songs, &SongOption{
				song:       song,
				difficulty: diff,
			})
		}
	}

	//orde by difficulty,
	//then by song title
	difficulties := map[types.Difficulty]int{}
	for _, o := range songs {
		difficulties[o.difficulty]++
	}
	sortedD := make([]types.Difficulty, 0, len(difficulties))
	for d := range difficulties {
		sortedD = append(sortedD, d)
	}
	sort.Slice(sortedD, func(i, j int) bool {
		return sortedD[i] < sortedD[j]
	})
	sort.Slice(songs, func(i, j int) bool {
		if songs[i].difficulty != songs[j].difficulty {
			return songs[i].difficulty < songs[j].difficulty
		}
		return songs[i].song.Title < songs[j].song.Title
	})

	s := &SongSelection{options: songs}

	center := ui.Point{
		X: 0.75,
		Y: 0.5,
	}

	songElements := ui.NewSongSelector()
	songElements.SetCenter(center)

	uiIdx := 0
	s.uiIdxToOptionIdx = make(map[int]int)
	for _, d := range sortedD {
		count := difficulties[d]
		e := ui.NewElement()
		e.SetText(l.String(d.String()) + " (" + (strconv.Itoa(count) + ")"))
		e.SetDisabled(true)
		e.SetScale(1.2)
		songElements.Add(e)
		uiIdx++

		for i, o := range s.options {
			if o.difficulty != d {
				continue
			}
			e := ui.NewElement()
			e.SetText(o.song.Title)
			e.SetTrigger(func() {
				s.SetNextState(types.GameStatePlay, &PlayArgs{
					Song:       o.song,
					Difficulty: o.difficulty,
				})
			})
			songElements.Add(e)
			s.uiIdxToOptionIdx[uiIdx] = i
			uiIdx++
		}
	}
	s.songList = songElements

	details := NewSongDetails()
	s.details = details
	if len(s.options) > 0 {
		song := s.options[0].song
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
	idx := s.songList.GetIndex()
	if idx != s.currentIdx {
		s.currentIdx = idx
		optionIdx := s.uiIdxToOptionIdx[s.currentIdx]
		song := s.options[optionIdx].song
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
