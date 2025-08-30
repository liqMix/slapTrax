package state

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/assets"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/l"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
)

var lastEditorIdx = 0

type EditorSongSelection struct {
	types.BaseGameState

	options          []*SongOption
	details          *ui.SongDetails
	songList         *ui.SongSelector
	currentIdx       int
	uiIdxToOptionIdx map[int]int
}

func NewEditorSongSelectionState() *EditorSongSelection {
	audio.FadeOutBGM()

	s := &EditorSongSelection{}
	s.SetAction(input.ActionBack, func() {
		audio.StopAll()
		audio.PlaySFX(audio.SFXBack)
		s.SetNextState(types.GameStateChartSelection, nil)
	})

	// Load bundled songs
	allSongs := assets.GetAllLoadedSongs()
	logger.Debug("Found %d loaded songs for editor", len(allSongs))
	songs := make([]*SongOption, 0)
	for _, song := range allSongs {
		difficulties := song.GetDifficulties()
		logger.Debug("Song '%s' has %d difficulties: %v", song.Title, len(difficulties), difficulties)
		for _, diff := range difficulties {
			songs = append(songs, &SongOption{
				song:       song,
				difficulty: diff,
			})
		}
	}

	// Load user charts dynamically from file system
	userCharts, err := assets.GetUserCharts()
	if err != nil {
		logger.Warn("Failed to load user charts: %v", err)
	} else {
		logger.Debug("Found %d user charts for editor", len(userCharts))
		for _, userChart := range userCharts {
			// Convert to full song object with chart data
			userSong, err := assets.ConvertUserChartToSong(userChart)
			if err != nil {
				logger.Warn("Failed to convert user chart %s: %v", userChart.FilePath, err)
				continue
			}

			// Add all difficulties from the user chart
			for difficulty, chart := range userSong.Charts {
				if chart != nil {
					songs = append(songs, &SongOption{
						song:       userSong,
						difficulty: difficulty,
					})
				}
			}
		}
	}
	logger.Debug("Created %d song options total for editor (including user charts)", len(songs))

	// Order by difficulty, then by song title
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
				// Load selected song into editor
				if o.song.Hash != "" && len(o.song.Hash) > 4 && o.song.Hash[len(o.song.Hash)-5:] == ".json" {
					// This is a user chart - pass the file path
					s.SetNextState(types.GameStateEditor, &EditorArgs{
						Song:      nil,
						ChartPath: o.song.Hash, // Hash contains the file path for user charts
						AudioPath: "",
					})
				} else {
					// This is a bundled song - pass the song object
					s.SetNextState(types.GameStateEditor, &EditorArgs{
						Song:      o.song,
						ChartPath: "",
						AudioPath: "",
					})
				}
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
		if lastEditorIdx > 0 {
			idx = lastEditorIdx
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

func (s *EditorSongSelection) Update() error {
	s.BaseGameState.Update()

	s.songList.Update()
	idx := s.songList.GetIndex()
	if idx != s.currentIdx {
		s.currentIdx = idx
		lastEditorIdx = s.currentIdx

		optionIdx := s.uiIdxToOptionIdx[s.currentIdx]
		song := s.options[optionIdx].song
		diff := s.options[optionIdx].difficulty
		s.details.UpdateDetails(song, diff)
		audio.PlaySongPreview(song)
	}

	s.details.Update()
	return nil
}

func (s *EditorSongSelection) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.songList.Draw(screen, opts)
	s.details.Draw(screen, opts)
}