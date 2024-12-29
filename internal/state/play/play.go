package play

import (
	"time"

	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
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
	Tracks     []*song.Track
	Score      score.Score

	startTime   time.Time
	elapsedTime int64
}

func New(arg interface{}) *State {
	args, ok := arg.(PlayArgs)
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
		elapsedTime: 0,
		startTime:   time.Now(),
	}
}

func (p *State) IsTrackPressed(trackName song.TrackName) bool {
	keys := TrackNameToKeys[trackName]
	return input.AnyKeysPressed(keys)
}

func (p *State) inGracePeriod() bool {
	if audio.IsSongPlaying() {
		return false
	}

	// Update elapsed time
	p.elapsedTime = time.Since(p.startTime).Milliseconds() - config.GracePeriod

	// Start the audio when the elapsed time is equal to the audio offset
	if p.elapsedTime >= config.AUDIO_OFFSET {
		audio.PlaySong()
		return false
	}
	return true
}

func (p *State) Update() (*types.GameState, interface{}, error) {
	if !p.inGracePeriod() {
		p.elapsedTime = int64(audio.CurrentSongPositionMS()) + config.AUDIO_OFFSET
	}

	// Handle input
	for _, track := range p.Tracks {
		keys := TrackNameToKeys[track.Name]
		if input.AnyKeysJustPressed(keys) {
			track.Activate(p.elapsedTime)
		} else if input.AllKeysReleased(keys) {
			track.Release(p.elapsedTime)
		}
	}
	// Update the tracks
	for _, track := range p.Tracks {
		track.Update(p.elapsedTime)
		for _, hit := range track.NewHits() {
			p.Score.AddHit(hit)
		}
	}
	return nil, nil, nil
}

func (p *State) CurrentTime() int64 {
	return p.elapsedTime
}
