package types

import (
	"math"
)

type HitRecord struct {
	Note *Note

	HitDiff     int64
	ReleaseDiff int64

	HitRating HitRating
	HitTiming HitTiming
}

type HitTiming int

const (
	HitTimingEarly HitTiming = iota
	HitTimingLate
	HitTimingNone
)

func (r HitTiming) String() string {
	switch r {
	case HitTimingEarly:
		return "Early"
	case HitTimingLate:
		return "Late"
	}
	return "None"
}

func GetHitTiming(diff int64) HitTiming {
	if diff > int64(Slap) {
		return HitTimingEarly
	} else if diff < int64(-Slap) {
		return HitTimingLate
	}
	return HitTimingNone
}

type HitRating int

const (
	Slap HitRating = iota
	Slip
	Slop
	None
)

func (r HitRating) String() string {
	switch r {
	case Slap:
		return "SLAP"
	case Slip:
		return "SLIP"
	case Slop:
		return "SLOP"
	}
	return ""
}

func (r HitRating) Value() float64 {
	switch r {
	case Slap:
		return 1
	case Slip:
		return 0.5
	}
	return 0
}

func (r HitRating) Color() GameColor {
	switch r {
	case Slap:
		return Green
	case Slip:
		return Yellow
	case Slop:
		return Gray
	}
	return White
}

// Loosen the window for early hits
var earlyScale = 0.4

func (r HitRating) Window(early bool) float64 {
	scale := 1.0
	if early {
		scale = earlyScale
	}
	switch r {
	case Slap:
		return 60
	case Slip:
		return 120 * scale
	}
	return 0
}

func GetHitRating(diff int64) HitRating {
	d := math.Abs(float64(diff))

	early := GetHitTiming(diff) == HitTimingEarly
	if d < Slap.Window(early) {
		return Slap
	} else if d < Slip.Window(early) {
		return Slip
	}

	return None
}

var EarliestWindow = int64(Slip.Window(true))
var LatestWindow = int64(Slip.Window(false))
