package state

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/cache"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/user"
)

type PlayArgs struct {
	Song       *types.Song
	Difficulty types.Difficulty
}

type Play struct {
	types.BaseGameState
	Song        *types.Song
	Difficulty  types.Difficulty
	Tracks      []*types.Track
	Score       *types.Score
	Chart       *types.Chart
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
	p.SetAction(input.ActionBack, p.pause)
	p.SetNotNavigable()
	return p
}

func (p *Play) GetTravelTime() int64 {
	return int64(travelTime / user.S().LaneSpeed)
}

func (p *Play) pause() {
	p.SetNextState(types.GameStatePause,
		&PauseArgs{
			Song:       p.Song,
			Difficulty: p.Difficulty,
			Cb: func() {
				p.countTicks = p.Song.GetCountdownTicks(true)
				p.startTime = time.Now()
				if p.elapsedTime >= user.S().AudioOffset {
					audio.SetSongPositionMS(int(p.elapsedTime - p.getGracePeriod()))
					audio.ResumeSong()
				}
			},
		})
}

func (p *Play) CurrentTime() int64 {
	return p.elapsedTime
}

func (p *Play) MaxTrackTime() int64 {
	return p.elapsedTime + p.GetTravelTime()
}

func (p *Play) getGracePeriod() int64 {
	if currentPos := audio.CurrentSongPositionMS(); currentPos <= 0 {
		return p.Song.GetBeatInterval() * 8
	}
	return p.Song.GetBeatInterval() * 4
}

func (p *Play) inGracePeriod() bool {
	// Determine current time based off grace period,
	// will be negative until the song starts
	p.elapsedTime = time.Since(p.startTime).Milliseconds() - p.getGracePeriod()

	// Play the starting ticks
	if len(p.countTicks) > 1 && p.elapsedTime >= (p.countTicks[0]+user.S().AudioOffset) {
		audio.PlaySFX(audio.SFXHat)
		p.countTicks = p.countTicks[1:]
	}

	// Start the audio when the elapsed time is equal to the offset start time
	if p.elapsedTime >= user.S().AudioOffset {
		input.K.ForceReset()
		audio.PlaySong()
		return false
	}
	return true
}

func (p *Play) Update() error {
	p.BaseGameState.Update()
	if cache.Path.IsBuilding() {
		return nil
	}

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
		activeBefore := track.Active
		track.Update(p.elapsedTime, p.GetTravelTime(), p.MaxTrackTime())
		if !activeBefore && track.Active {
			audio.PlayTrackSFX(track.Name)
		}
	}

	return nil
}

func (p *Play) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {}
