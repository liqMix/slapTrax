package play

import (
	"time"

	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/score"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type PlayArgs struct {
	Song       *song.Song
	Difficulty song.Difficulty
}

type State struct {
	Song       *song.Song
	Difficulty song.Difficulty
	Tracks     []song.Track
	Score      score.Score

	startTime   time.Time
	elapsedTime int64
}

func New(arg interface{}) *State {
	args, ok := arg.(*PlayArgs)
	if !ok {
		panic("Invalid play state args")
	}
	song := args.Song
	difficulty := args.Difficulty
	chart, ok := song.Charts[difficulty]
	if !ok {
		panic("No chart for difficulty")
	}

	audio.InitSong(song)
	return &State{
		Song:        song,
		Difficulty:  difficulty,
		Tracks:      chart.Tracks,
		elapsedTime: -config.GRACE_PERIOD,
		startTime:   time.Now(),
	}
}

func (p *State) inGracePeriod() bool {
	if audio.IsSongPlaying() {
		return false
	}

	// Update elapsed time
	p.elapsedTime = time.Since(p.startTime).Milliseconds() - config.GRACE_PERIOD

	// Start the audio when the elapsed time is equal to the audio offset
	if p.elapsedTime >= config.AUDIO_OFFSET {
		audio.PlaySong()
		return false
	}
	return true
}

func (p *State) Update() (*types.GameState, interface{}, error) {
	if p.inGracePeriod() {
		return nil, nil, nil
	}

	p.elapsedTime = int64(audio.CurrentSongPositionMS()) + config.AUDIO_OFFSET
	for track := range p.Tracks {
		p.Tracks[track].Update(p.elapsedTime)
	}
	return nil, nil, nil
}
