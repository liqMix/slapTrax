package state

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

type PlayArgs struct {
	Song       *types.Song
	Difficulty types.Difficulty
}

type Play struct {
	types.BaseGameState
	Song       *types.Song
	Difficulty types.Difficulty
	Tracks     []*types.Track
	Score      *types.Score
	Chart      *types.Chart

	startTime   time.Time
	elapsedTime int64
	countTicks  []int64
}

const travelTime float64 = 10000

func NewPlayState(args *PlayArgs) *Play {
	song := args.Song
	difficulty := args.Difficulty
	chart, ok := song.Charts[difficulty]
	if !ok {
		panic("No chart for difficulty")
	}

	tracks := chart.Tracks

	// Get the song audio ready
	audio.StopAll()
	audio.InitSong(song)
	for _, track := range tracks {
		track.Reset()
	}

	// set the grace period to be 4 quarter notes
	p := &Play{
		Song:        song,
		Difficulty:  difficulty,
		Tracks:      tracks,
		Chart:       chart,
		Score:       types.NewScore(song, difficulty),
		elapsedTime: 0,
		startTime:   time.Now(),
		countTicks:  song.GetCountdownTicks(false),
	}
	return p
}

func (p *Play) GetTravelTime() int64 {
	return int64(travelTime / user.S().LaneSpeed)
}

func (p *Play) restart() {
	audio.StopAll()
	p.SetNextState(types.GameStatePlay, &PlayArgs{
		Song:       p.Song,
		Difficulty: p.Difficulty,
	})
}

func (p *Play) CurrentTime() int64 {
	return p.elapsedTime
}

func (p *Play) MaxTrackTime() int64 {
	return p.elapsedTime + p.GetTravelTime()
}

func (p *Play) getGracePeriod() int64 {
	if currentPos := audio.CurrentSongPositionMS(); currentPos < 0 {
		return p.Song.GetBeatInterval() * 8
	}
	return p.Song.GetBeatInterval() * 4
}

func (p *Play) getOffsetTime() int64 {
	return user.S().AudioOffset + config.INHERENT_OFFSET
}

func (p *Play) inGracePeriod() bool {
	// Account for offets
	offsetStartTime := p.getOffsetTime()

	// Determine the grace period
	gracePeriod := offsetStartTime + p.getGracePeriod()

	// Determine current time based off grace period,
	// will be negative until the song starts
	p.elapsedTime = time.Since(p.startTime).Milliseconds() - gracePeriod

	// Play the starting ticks
	if len(p.countTicks) > 1 && p.elapsedTime >= (p.countTicks[0]+offsetStartTime) {
		audio.PlaySFX(audio.SFXHat)
		p.countTicks = p.countTicks[1:]
	}

	// Start the audio when the elapsed time is equal to the offset start time
	if p.elapsedTime >= offsetStartTime {
		audio.PlaySong()
		return false
	}
	return true
}

func (p *Play) handleAction(action PlayAction) {
	switch action {
	case RestartAction:
		p.restart()
	case PauseAction:
		audio.PauseSong()
		p.SetNextState(types.GameStatePause,
			&PauseArgs{
				song:       p.Song,
				difficulty: p.Difficulty,

				cb: func() {
					offsetTime := p.getOffsetTime()
					p.countTicks = p.Song.GetCountdownTicks(true)
					p.startTime = time.Now()
					audio.SetSongPositionMS(int(p.elapsedTime + offsetTime))
				},
			})
		return
	}
}

func (p *Play) Update() error {
	if !audio.IsSongPlaying() {
		if !p.inGracePeriod() {
			stillPlaying := false
			for _, track := range p.Tracks {
				if track.HasMoreNotes() {
					stillPlaying = true
					break
				}
			}
			if !stillPlaying {
				p.SetNextState(types.GameStateResult, &ResultStateArgs{
					Score: p.Score,
				})
			}
		}
	} else {
		p.elapsedTime = int64(audio.CurrentSongPositionMS()) + user.S().AudioOffset
	}

	// Update the tracks
	for _, track := range p.Tracks {
		p.updateTrack(track, p.elapsedTime, p.Score)
	}

	// Handle input
	for action, keys := range PlayActions {
		if input.K.AreAny(keys, input.JustPressed) {
			p.handleAction(action)
		}
	}

	return nil
}

func (p *Play) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {}
