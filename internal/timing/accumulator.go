package timing

import "time"

// FixedStep handles fixed timestep updates with interpolation
type FixedStep[T any] struct {
	UpdateRate float64   // Fixed update rate in seconds (e.g. 1/240)
	MaxSteps   int       // Maximum steps per frame to prevent spiral of death
	LastUpdate time.Time // Last update time

	accumulator float64
	prev, curr  T // Previous and current states for interpolation

	// UpdateFn is called at a fixed timestep to update the state
	UpdateFn func(T) T
}

// NewFixedStep creates a new fixed timestep accumulator
func NewFixedStep[T any](updateRate float64, maxSteps int, initial T, updateFn func(T) T) *FixedStep[T] {
	return &FixedStep[T]{
		UpdateRate: updateRate,
		MaxSteps:   maxSteps,
		LastUpdate: time.Now().Round(0),
		prev:       initial,
		curr:       initial,
		UpdateFn:   updateFn,
	}
}

// Update runs the fixed timestep update loop and returns interpolated state
// The alpha value can be used for custom interpolation if needed
func (f *FixedStep[T]) Update() (state T, alpha float64) {
	now := time.Now().Round(0)
	frameTime := now.Sub(f.LastUpdate).Seconds()
	f.LastUpdate = now

	f.accumulator += frameTime
	steps := 0

	// Update in fixed timesteps
	for f.accumulator >= f.UpdateRate && steps < f.MaxSteps {
		f.prev = f.curr
		f.curr = f.UpdateFn(f.curr)

		f.accumulator -= f.UpdateRate
		steps++
	}

	// Calculate interpolation alpha
	alpha = f.accumulator / f.UpdateRate

	return f.curr, alpha
}

// GetPrevious returns the previous state
func (f *FixedStep[T]) GetPrevious() T {
	return f.prev
}

// GetCurrent returns the current state
func (f *FixedStep[T]) GetCurrent() T {
	return f.curr
}

// Reset resets the accumulator with a new state
func (f *FixedStep[T]) Reset(state T) {
	f.prev = state
	f.curr = state
	f.accumulator = 0
	f.LastUpdate = time.Now().Round(0)
}

// Common helper functions for interpolation

// LerpFloat interpolates between two float64 values
func LerpFloat(start, end, alpha float64) float64 {
	return start + (end-start)*alpha
}

// LerpInt interpolates between two int values
func LerpInt(start, end int, alpha float64) int {
	return start + int(float64(end-start)*alpha)
}
