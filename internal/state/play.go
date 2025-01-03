package state

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
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
	Song        *types.Song
	Difficulty  types.Difficulty
	Tracks      []*types.Track
	Score       *types.Score
	Chart       *types.Chart
	startTime   time.Time
	elapsedTime int64
	TravelTime  float64
}

func NewPlayState(args *PlayArgs) *Play {
	playSong := args.Song
	difficulty := args.Difficulty
	chart, ok := playSong.Charts[difficulty]
	if !ok {
		panic("No chart for difficulty")
	}

	tracks := chart.Tracks

	// Get the song audio ready
	assets.InitSong(playSong)
	for _, track := range tracks {
		track.Reset()
	}
	return &Play{
		Song:        playSong,
		Difficulty:  difficulty,
		Tracks:      tracks,
		Chart:       chart,
		Score:       types.NewScore(chart.TotalNotes),
		elapsedTime: 0,
		startTime:   time.Now(),
		TravelTime:  float64(config.TRAVEL_TIME) / user.S().Gameplay.NoteSpeed,
	}
}

// func (p *Play) IsTrackPressed(trackName song.TrackName) bool {
// 	keys := TrackNameToKeys[trackName]
// 	return input.K.AreAny(keys, input.Held)
// }

func (p *Play) inGracePeriod() bool {
	if assets.IsSongPlaying() {
		return false
	}

	// Update elapsed time
	p.elapsedTime = time.Since(p.startTime).Milliseconds() - int64(p.TravelTime)/2

	// Start the audio when the elapsed time is equal to the audio offset
	if p.elapsedTime >= -200+user.S().Gameplay.AudioOffset {
		assets.PlaySong()
		return false
	}
	return true
}

func (p *Play) handleAction(action PlayAction) {
	switch action {
	case RestartAction:
		// Stop the song
		assets.StopSong()
		assets.InitSong(p.Song)

		// Reset the tracks
		for _, track := range p.Tracks {
			track.Reset()
		}

		// Reset the score
		p.Score.Reset()
		p.elapsedTime = 0
		p.startTime = time.Now()
		return
	case PauseAction:
		assets.PauseSong()
		p.SetNextState(types.GameStatePause,
			&PauseArgs{
				song:       p.Song,
				difficulty: p.Difficulty,
			})
		return
	}
}

func (p *Play) Update() error {
	if !p.inGracePeriod() {
		p.elapsedTime = int64(assets.CurrentSongPositionMS()) + user.S().Gameplay.AudioOffset
	}

	// Update the tracks
	for _, track := range p.Tracks {
		p.updateTrack(track, p.elapsedTime, p.TravelTime, p.Score)
	}

	// Handle input
	for action, keys := range PlayActions {
		if input.K.AreAny(keys, input.JustPressed) {
			p.handleAction(action)
		}
	}

	// for _, track := range p.Tracks {
	// 	keys := TrackNameToKeys[track.Name]
	// 	if input.K.AreAny(keys, input.JustPressed) {
	// 		track.Activate(p.elapsedTime)
	// 	} else if input.K.AreAll(keys, input.None) {
	// 		track.Release(p.elapsedTime)
	// 	}
	// }

	return nil
}

func (p *Play) CurrentTime() int64 {
	return p.elapsedTime
}

func (p *Play) updateTrackInput(t *types.Track) {
	if !t.Active && !t.StaleActive {
		if input.K.AreAny(types.TrackNameToKeys[t.Name], input.JustPressed) {
			t.Active = true
			return
		}
	}

	if t.StaleActive || t.Active {
		if !input.K.AreAny(types.TrackNameToKeys[t.Name], input.Held) {
			t.Active = false
			t.StaleActive = false
		}
	}
}

func (p *Play) updateTrack(t *types.Track, currentTime int64, travelTime float64, score *types.Score) {
	p.updateTrackInput(t)
	// Reset the new hits
	notes := make([]*types.Note, 0, len(t.ActiveNotes))

	// Only update notes that are currently visible
	for _, n := range t.ActiveNotes {
		n.Update(currentTime, travelTime)

		// Active Track
		if t.Active {
			if n.WasHit() && !n.IsHoldNote() {
				continue
			}

			if !t.StaleActive {
				if n.Hit(currentTime+user.S().Gameplay.InputOffset, score) {
					t.StaleActive = true
					if n.IsHoldNote() {
						notes = append(notes, n)
					}
					continue
				}
			}

		} else {
			// not active track
			if n.IsHoldNote() && n.WasHit() {
				if n.WasReleased() {
					continue
				}

				n.Release(currentTime + user.S().Gameplay.InputOffset)
				continue
			}
		}

		// Note not yet reached the out of bounds window
		windowEnd := n.Target + types.LatestWindow
		releaseWindowEnd := n.TargetRelease + types.LatestWindow
		if currentTime < windowEnd || (n.IsHoldNote() && currentTime < releaseWindowEnd) {
			notes = append(notes, n)
			continue
		}

		// Drop expired notes
		n.Miss(score)

		// // Handle hold notes
		// if n.IsHoldNote() && n.WasHit() {
		// 	// If player released or the hold window expired while holding
		// 	// if t.releaseTs > minTs {
		// 	// 	n.Release(t.releaseTs)
		// 	// } else if currentTime > releaseWindowEnd {
		// 	// 	n.Release(currentTime)
		// 	// } else {
		// 	// 	notes = append(notes, n) // Keep held or missed-hit hold notes
		// 	// }
		// 	// if hitRating != None {
		// 	// TODO: determine handling release scores
		// 	// t.newHits = append(t.newHits, hitRating)
		// 	// }
		// 	continue
		// }
	}

	// Add new approaching notes
	spawnTime := currentTime + int64(travelTime*2)
	if t.NextNoteIndex < len(t.AllNotes) {
		for i := t.NextNoteIndex; i < len(t.AllNotes); i++ {
			note := t.AllNotes[i]
			if note.Target > spawnTime {
				break
			}
			notes = append(notes, note)
			t.NextNoteIndex = i + 1
		}
	}

	t.ActiveNotes = notes

	// dont stay active if no notes in window
	if t.Active && !t.StaleActive {
		for _, n := range t.ActiveNotes {
			if n.InWindow(currentTime-types.EarliestWindow, currentTime+types.LatestWindow) {
				return
			}
		}
		t.StaleActive = true
	}
}

func (p *Play) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {}
