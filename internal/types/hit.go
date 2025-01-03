package types

import (
	"math"
)

type HitRecord struct {
	Note      *Note
	Diff      int64
	HitType   HitType
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
	if diff > 0 {
		return HitTimingEarly
	} else if diff < 0 {
		return HitTimingLate
	}
	return HitTimingNone
}

type HitType int

const (
	Perfect HitType = iota
	Good
	Bad
	Miss
	None
)

func (r HitType) String() string {
	switch r {
	case Perfect:
		return "Perfect"
	case Good:
		return "Good"
	case Bad:
		return "Bad"
	case Miss:
		return "Miss"
	}
	return "None"
}

func (r HitType) Value() int {
	switch r {
	case Perfect:
		return 10
	case Good:
		return 5
	case Bad:
		return 0
	}
	return 0
}

// Loosen the window for early hits
var earlyScale = 0.4

func (r HitType) Window(early bool) float64 {
	scale := 1.0
	if early {
		scale = earlyScale
	}
	switch r {
	case Perfect:
		return 60
	case Good:
		return 100 * scale
	case Bad:
		return 150 * scale
	}
	return 0
}

func GetHitRating(diff int64) HitType {
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
