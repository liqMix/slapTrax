package state

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

var lastIdx = 0

type SongSelectionArgs struct {
	Song *types.Song
}
type SongSelection struct {
	types.BaseGameState

	options          []*SongOption
	details          *ui.SongDetails
	songList         *ui.SongSelector
	leaderboard      *ui.Leaderboard
	currentIdx       int
	uiIdxToOptionIdx map[int]int
}

type SongOption struct {
	song       *types.Song
	difficulty types.Difficulty
}

func NewSongSelectionState() *SongSelection {
	audio.FadeOutBGM()

	s := &SongSelection{leaderboard: ui.NewLeaderboard()}
	s.SetAction(input.ActionBack, func() {
		audio.StopAll()
		audio.PlaySFX(audio.SFXBack)
		s.SetNextState(types.GameStateTitle, nil)
	})

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

	//order by difficulty,
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
	s.options = songs

	center := ui.Point{
		X: 0.75,
		Y: 0.5,
	}
	songElements := ui.NewSongSelector()
	songElements.SetCenter(center)

	uiIdx := 0
	s.uiIdxToOptionIdx = make(map[int]int)

	size := ui.Point{X: 0.25, Y: 0.15}
	for _, d := range sortedD {
		e := ui.NewElement()
		e.SetText(l.String(d.String()))
		e.SetDisabled(true)
		e.SetSize(size)
		e.SetTextScale(2)
		e.SetTextBold(true)
		e.SetTextColor(d.Color())
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

	details := ui.NewSongDetails()
	s.details = details
	if len(s.options) > 0 {
		idx := 0
		if lastIdx > 0 {
			idx = lastIdx
			s.songList.Select(idx)
		}
		idx = s.uiIdxToOptionIdx[idx]
		song := s.options[idx].song
		diff := s.options[idx].difficulty
		s.details.UpdateDetails(song, diff)
		audio.PlaySongPreview(song)
	}
	return s
}

func (s *SongSelection) Update() error {
	s.BaseGameState.Update()

	s.songList.Update()
	idx := s.songList.GetIndex()
	if idx != s.currentIdx {
		s.currentIdx = idx
		lastIdx = s.currentIdx

		optionIdx := s.uiIdxToOptionIdx[s.currentIdx]
		song := s.options[optionIdx].song
		diff := s.options[optionIdx].difficulty
		s.details.UpdateDetails(song, diff)
		audio.PlaySongPreview(song)
		s.leaderboard.FetchScores(song.Hash, int(diff), float64(song.BPM))
	}

	s.details.Update()
	s.leaderboard.Update()
	return nil
}

func (s *SongSelection) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.songList.Draw(screen, opts)
	s.details.Draw(screen, opts)
	s.leaderboard.Draw(screen, opts)
}
