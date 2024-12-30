package play

import (
	"fmt"
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
	Song        *song.Song
	Difficulty  song.Difficulty
	Tracks      []*song.Track
	Score       score.Score
	Chart       *song.Chart
	startTime   time.Time
	elapsedTime int64
}

func New(arg interface{}) *State {
	args, ok := arg.(PlayArgs)
	if !ok {
		panic("Invalid play state args")
	}
	playSong := args.Song
	difficulty := args.Difficulty
	chart, ok := playSong.Charts[difficulty]
	if !ok {
		panic("No chart for difficulty")
	}

	tracks := chart.Tracks

	// If we don't have edge tracks, remove them
	if config.NO_EDGE_TRACKS || !chart.HasEdgeTracks() {
		fmt.Println("Removing edge tracks")
		newTracks := []*song.Track{}
		for _, track := range chart.Tracks {
			switch track.Name {
			case song.EdgeTop, song.EdgeTap1, song.EdgeTap2, song.EdgeTap3:
				continue
			default:
				newTracks = append(newTracks, track)
			}
		}
		tracks = newTracks
	}

	// Get the song audio ready
	audio.InitSong(playSong)
	return &State{
		Song:        playSong,
		Difficulty:  difficulty,
		Tracks:      tracks,
		Chart:       chart,
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
