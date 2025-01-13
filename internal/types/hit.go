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
	if diff > int64(Perfect) {
		return HitTimingEarly
	} else if diff < int64(-Perfect) {
		return HitTimingLate
	}
	return HitTimingNone
}

type HitRating int

const (
	Perfect HitRating = iota
	Good
	Bad
	Miss
	None
)

func (r HitRating) String() string {
	switch r {
	case Perfect:
		return "PERFECT"
	case Good:
		return "GOOD"
	case Bad:
		return "OK"
	case Miss:
		return "MISS"
	}
	return ""
}

func (r HitRating) Value() float64 {
	switch r {
	case Perfect:
		return 1
	case Good:
		return 0.5
	case Bad:
		return 0.25
	}
	return 0
}

func (r HitRating) Color() GameColor {
	switch r {
	case Perfect:
		return Green
	case Good:
		return Yellow
	case Bad:
		return Orange
	case Miss:
		return Red
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
	case Perfect:
		return 60
	case Good:
		return 70 * scale
	case Bad:
		return 120 * scale
	}
	return 0
}

func GetHitRating(diff int64) HitRating {
	d := math.Abs(float64(diff))

	early := GetHitTiming(diff) == HitTimingEarly
	if d < Perfect.Window(early) {
		return Perfect
	} else if d < Good.Window(early) {
		return Good
	} else if d < Bad.Window(early) {
		return Bad
	}

	return None
}

var EarliestWindow = int64(Bad.Window(true))
var LatestWindow = int64(Bad.Window(false))
