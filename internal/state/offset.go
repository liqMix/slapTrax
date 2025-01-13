package state

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/user"
)

const (
	offsetTravelTime int64   = 1500
	centerWindow     float64 = 0.01
)

type Offset struct {
	types.BaseGameState

	bgmWasPlaying   bool
	cb              func()
	startTime       time.Time // Renamed from anchorTime for clarity
	elapsedTime     int64
	backwards       bool
	travelTime      int64
	playedTick      bool
	initAudioOffset int64
	initInputOffset int64

	NoteProgress float64
	CenterWindow float64
	AudioOffset  int64
	InputOffset  int64
	HitDiff      int64
}

func NewOffsetState(args *FloatStateArgs) *Offset {
	bgmWasPlaying := false
	if audio.GetBGM() != nil && audio.GetBGM().IsPlaying() {
		audio.FadeOutBGM()
		bgmWasPlaying = true
	}

	// Gotta have some volume...
	if user.S().SFXVolume < 0.5 {
		audio.SetSFXVolume(0.5)
	}

	aOffset := int64(user.S().AudioOffset)
	iOffset := int64(user.S().InputOffset)

	now := time.Now()
	state := &Offset{
		bgmWasPlaying:   bgmWasPlaying,
		cb:              args.Cb,
		AudioOffset:     aOffset,
		InputOffset:     iOffset,
		initAudioOffset: aOffset,
		initInputOffset: iOffset,
		NoteProgress:    0.5,
		CenterWindow:    centerWindow,
		travelTime:      offsetTravelTime,
		startTime:       now,
	}

	return state
}

func resetVolume() {
	audio.SetSFXVolume(user.S().SFXVolume)
}

func (s *Offset) getCenterTime() int64 {
	return (s.travelTime / 2) + s.AudioOffset
}

func (s *Offset) handleInput() error {
	if input.JustActioned(input.ActionUp) {
		s.AudioOffset += 5
	} else if input.JustActioned(input.ActionDown) {
		s.AudioOffset -= 5
	} else if input.JustActioned(input.ActionLeft) {
		s.InputOffset -= 5
	} else if input.JustActioned(input.ActionRight) {
		s.InputOffset += 5
	} else if input.K.Is(ebiten.KeySpace, input.JustPressed) {
		hitTime := s.elapsedTime + s.InputOffset
		s.HitDiff = s.getCenterTime() - hitTime
	} else if input.K.Is(ebiten.KeyR, input.JustPressed) {
		s.resetOffsets()
	} else if input.K.Is(ebiten.Key0, input.JustPressed) {
		s.AudioOffset = 0
		s.InputOffset = 0
	} else if input.JustActioned(input.ActionSelect) {
		resetVolume()
		audio.PlaySFX(audio.SFXSelect)
		return s.saveAndExit()
	} else if input.JustActioned(input.ActionBack) {
		resetVolume()
		return s.exit()
	}

	return nil
}

func (s *Offset) resetOffsets() {
	s.AudioOffset = s.initAudioOffset
	s.InputOffset = s.initInputOffset
}

func (s *Offset) saveAndExit() error {
	user.S().AudioOffset = s.AudioOffset
	user.S().InputOffset = s.InputOffset
	if s.cb != nil {
		s.cb()
	}
	return s.exit()
}

func (s *Offset) exit() error {
	audio.StopAll()
	if s.bgmWasPlaying {
		audio.FadeInBGM()
	}
	s.SetNextState(types.GameStateBack, nil)
	return nil
}

func (s *Offset) Update() error {
	s.BaseGameState.Update()
	now := time.Now()
	elapsed := now.Sub(s.startTime).Milliseconds()
	progress := float64(elapsed) / float64(s.travelTime)
	if s.backwards {
		progress = 1 - progress
	}

	// Check bounds and handle wraparound
	if progress <= 0 || progress >= 1 {
		progress = math.Max(0, math.Min(1, progress))
		s.backwards = progress > 0
		s.startTime = now
		s.playedTick = false
	}

	if !s.playedTick {
		tickAtProgress := ((float64(s.travelTime/2) + float64(s.AudioOffset)) / float64(s.travelTime))

		if s.backwards {
			tickAtProgress = 1 - tickAtProgress
		}

		// Check if current progress is within threshold of either tick point
		diff := math.Abs(progress - tickAtProgress)

		if diff <= s.CenterWindow {
			audio.PlaySFX(audio.SFXHat)
			s.playedTick = true
		}
	}

	s.NoteProgress = progress
	s.elapsedTime = elapsed

	if err := s.handleInput(); err != nil {
		return err
	}

	// Ensure offsets stay within bounds
	s.AudioOffset = clamp(s.AudioOffset, -s.travelTime, s.travelTime)
	s.InputOffset = clamp(s.InputOffset, -s.travelTime, s.travelTime)

	return nil
}

func clamp(val, min, max int64) int64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (s *Offset) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {}
