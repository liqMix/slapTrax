package state

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/timing"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

const (
	offsetCenter     = 2000 // The ms center tick of the offset test sound
	inputDelayFrames = 100
	updateRate       = 1.0 / 240.0 // 240Hz fixed update rate
	maxSteps         = 3
	maxOffset        = 1000
	minOffset        = -1000
	maxHitDiffs      = 15   // Increased for better averaging
	maxAutoAdj       = 5    // Reduced for smoother adjustments
	hitWindowMs      = 200  // Valid hit window in milliseconds
	centerThreshold  = 0.01 // Threshold for center tick playback

	// Weights for different timing windows
	perfectWeight = 1.0
	greatWeight   = 0.8
	goodWeight    = 0.5
)

type Offset struct {
	types.BaseGameState

	startTime        time.Time // Renamed from anchorTime for clarity
	elapsedTime      int64
	backwards        bool
	travelTime       int64
	playedCenterTick bool
	initAudioOffset  int64
	initInputOffset  int64

	AutoAdjustInput bool
	AutoAdjustAudio bool

	NoteProgress float64
	AudioOffset  int64
	InputOffset  int64
	HitDiff      int64
	hitDiffs     []weightedHit

	timeStep       *timing.FixedStep[float64]
	lastUpdateTime time.Time
}

// weightedHit stores hit timing data with its weight for adjustment
type weightedHit struct {
	diff   int64
	weight float64
	time   time.Time
}

func NewOffsetState() *Offset {
	s := user.S()

	input.K.WatchKeys([]ebiten.Key{
		ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight,
	})

	aOffset := int64(s.Gameplay.AudioOffset)
	iOffset := int64(s.Gameplay.InputOffset)

	now := time.Now()
	state := &Offset{
		AudioOffset:     aOffset,
		InputOffset:     iOffset,
		initAudioOffset: aOffset,
		initInputOffset: iOffset,
		travelTime:      2000,
		NoteProgress:    0.5,
		startTime:       now,
		lastUpdateTime:  now,
		hitDiffs:        make([]weightedHit, 0, maxHitDiffs),
	}

	// Create fixed timestep handler with improved progress calculation
	state.timeStep = timing.NewFixedStep(updateRate, maxSteps, state.NoteProgress,
		func(progress float64) float64 {
			now := time.Now()
			elapsedMs := float64(now.Sub(state.startTime).Nanoseconds()) / float64(time.Millisecond)

			if state.backwards {
				return 1.0 - (elapsedMs / float64(state.travelTime))
			}
			return elapsedMs / float64(state.travelTime)
		})

	return state
}

// getHitWeight returns a weight based on timing accuracy
func getHitWeight(diff int64) float64 {
	absDiff := math.Abs(float64(diff))
	switch {
	case absDiff < 30:
		return perfectWeight
	case absDiff < 60:
		return greatWeight
	case absDiff < hitWindowMs/2:
		return goodWeight
	default:
		return 0.0 // Hit too far off
	}
}

// GetHitDiff calculates the timing difference with improved accuracy
func (s *Offset) GetHitDiff(hitTime int64) int64 {
	// Calculate the expected hit times for top, center, and bottom positions
	// Remove input offset from hit time since it's already added
	actualHitTime := hitTime - s.InputOffset

	// Adjust travel time for audio offset
	adjustedTravelTime := s.travelTime - s.AudioOffset

	targetTimes := []int64{
		adjustedTravelTime,     // Top
		adjustedTravelTime / 2, // Center
		0,                      // Bottom
	}

	// Find the closest target time
	var closestDiff int64
	minAbsDiff := int64(math.MaxInt64)

	for _, target := range targetTimes {
		diff := actualHitTime - target
		absDiff := abs(diff)

		if absDiff < minAbsDiff {
			minAbsDiff = absDiff
			closestDiff = diff
		}
	}

	// Only return valid hits within the hit window
	if minAbsDiff <= hitWindowMs {
		return closestDiff
	}
	return 0 // Invalid hit
}

// autoAdjustOffset uses a weighted moving average for smoother adjustments
func (s *Offset) autoAdjustOffset() int64 {
	if len(s.hitDiffs) == 0 {
		return 0
	}

	// Remove old hits
	now := time.Now()
	validHits := make([]weightedHit, 0, len(s.hitDiffs))
	for _, hit := range s.hitDiffs {
		if now.Sub(hit.time) < 30*time.Second {
			validHits = append(validHits, hit)
		}
	}
	s.hitDiffs = validHits

	if len(s.hitDiffs) == 0 {
		return 0
	}

	// Calculate weighted average
	var weightedSum float64
	var totalWeight float64

	for _, hit := range s.hitDiffs {
		if s.AutoAdjustAudio {
			// For audio offset:
			// If hit is late (positive), we need negative adjustment
			// If hit is early (negative), we need positive adjustment
			weightedSum -= float64(hit.diff) * hit.weight
		} else {
			// For input offset:
			// If hit is late (positive), we need negative adjustment to compensate
			// If hit is early (negative), we need positive adjustment to compensate
			weightedSum -= float64(hit.diff) * hit.weight // Same as audio!
		}
		totalWeight += hit.weight
	}

	if totalWeight == 0 {
		return 0
	}

	// Calculate adjustment with dynamic scaling
	adjustment := int64(weightedSum / totalWeight)

	// Scale adjustment based on magnitude of error
	scale := math.Min(1.0, math.Max(0.2, math.Abs(float64(adjustment))/100.0))
	adjustment = int64(float64(adjustment) * scale)

	// Apply maximum adjustment limit
	if adjustment > maxAutoAdj {
		adjustment = maxAutoAdj
	} else if adjustment < -maxAutoAdj {
		adjustment = -maxAutoAdj
	}

	return adjustment
}

func (s *Offset) Update() error {
	now := time.Now()
	// frameTime := now.Sub(s.lastUpdateTime)
	s.lastUpdateTime = now

	// Get interpolated progress with improved timing
	progress, _ := s.timeStep.Update()

	// Check bounds and handle wraparound
	if progress <= 0 || progress >= 1 {
		progress = math.Max(0, math.Min(1, progress))
		s.backwards = progress > 0
		s.startTime = now
		s.playedCenterTick = false
		s.timeStep.Reset(progress)
	} else if math.Abs(progress-0.5) < centerThreshold && !s.playedCenterTick {
		audio.StopAll()
		audio.PlaySFXWithOffset(audio.SFXOffset, offsetCenter+s.AudioOffset)
		s.playedCenterTick = true
	}

	s.NoteProgress = progress
	s.elapsedTime = now.Sub(s.startTime).Milliseconds()

	// Handle input with improved key checking
	if err := s.handleInput(); err != nil {
		return err
	}

	// Ensure offsets stay within bounds
	s.AudioOffset = clamp(s.AudioOffset, minOffset, maxOffset)
	s.InputOffset = clamp(s.InputOffset, minOffset, maxOffset)

	return nil
}

func (s *Offset) handleInput() error {
	if checkPress(ebiten.KeyArrowUp) {
		s.AudioOffset += 5
	} else if checkPress(ebiten.KeyArrowDown) {
		s.AudioOffset -= 5
	} else if checkPress(ebiten.KeyArrowLeft) {
		s.InputOffset -= 5
	} else if checkPress(ebiten.KeyArrowRight) {
		s.InputOffset += 5
	} else if input.K.Is(ebiten.KeySpace, input.JustPressed) {
		hitTime := s.elapsedTime
		s.HitDiff = s.GetHitDiff(hitTime)

		if s.HitDiff != 0 { // Only process valid hits
			weight := getHitWeight(s.HitDiff)
			hit := weightedHit{
				diff:   s.HitDiff,
				weight: weight,
				time:   time.Now(),
			}

			s.hitDiffs = append(s.hitDiffs, hit)
			if len(s.hitDiffs) > maxHitDiffs {
				s.hitDiffs = s.hitDiffs[1:]
			}

			if s.AutoAdjustAudio || s.AutoAdjustInput {
				adjustment := s.autoAdjustOffset()
				if s.AutoAdjustAudio {
					s.AudioOffset += adjustment
				} else {
					s.InputOffset += adjustment
				}
			}
		}
	} else if input.K.Is(ebiten.KeyR, input.JustPressed) {
		s.resetOffsets()
	} else if input.K.Is(ebiten.Key0, input.JustPressed) {
		s.AudioOffset = 0
		s.InputOffset = 0
	} else if input.K.Is(ebiten.KeyEnter, input.JustPressed) {
		return s.saveAndExit()
	} else if input.K.Is(ebiten.KeyEscape, input.JustPressed) {
		return s.exit()
	} else if input.K.Is(ebiten.KeyA, input.JustPressed) {
		s.toggleAudioAdjust()
	} else if input.K.Is(ebiten.KeyI, input.JustPressed) {
		s.toggleInputAdjust()
	}

	return nil
}

func (s *Offset) resetOffsets() {
	s.AudioOffset = s.initAudioOffset
	s.InputOffset = s.initInputOffset
	s.hitDiffs = make([]weightedHit, 0, maxHitDiffs)
}

func (s *Offset) saveAndExit() error {
	user.S().Gameplay.AudioOffset = s.AudioOffset
	user.S().Gameplay.InputOffset = s.InputOffset
	return s.exit()
}

func (s *Offset) exit() error {
	input.K.ClearWatchedKeys()
	audio.StopAll()
	s.SetNextState(types.GameStateBack, nil)
	return nil
}

func (s *Offset) toggleAudioAdjust() {
	s.AutoAdjustAudio = !s.AutoAdjustAudio
	if s.AutoAdjustAudio {
		s.AutoAdjustInput = false
	}
	s.hitDiffs = make([]weightedHit, 0, maxHitDiffs)
}

func (s *Offset) toggleInputAdjust() {
	s.AutoAdjustInput = !s.AutoAdjustInput
	if s.AutoAdjustInput {
		s.AutoAdjustAudio = false
	}
	s.hitDiffs = make([]weightedHit, 0, maxHitDiffs)
}

// Helper functions
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func clamp(value, min, max int64) int64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func checkPress(key ebiten.Key) bool {
	return input.K.Is(key, input.JustPressed) || input.K.IsKeyHeldFor(key, inputDelayFrames)
}

func (s *Offset) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {}
