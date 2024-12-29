package score

import "github.com/liqmix/ebiten-holiday-2024/internal/score/hit"

type ScoreRating string

type scoreRating struct {
	SSS ScoreRating
	SS  ScoreRating
	S   ScoreRating
	A   ScoreRating
	B   ScoreRating
	F   ScoreRating
}

var Rating = scoreRating{
	SSS: "SSS",
	SS:  "SS",
	S:   "S",
	A:   "A",
	B:   "B",
	F:   "F",
}

// scoreThresholds is the minimum percentage of perfect/good hits required to achieve a rating
var scoreThresholds map[ScoreRating]int = map[ScoreRating]int{
	Rating.SSS: 100,
	Rating.SS:  90,
	Rating.S:   80,
	Rating.A:   70,
	Rating.B:   60,
	Rating.F:   0,
}

type Score struct {
	TotalNotes int

	Rating  ScoreRating
	Perfect int
	Good    int
	Bad     int
	Miss    int

	Combo    int
	MaxCombo int
}

func NewScore(totalNotes int) *Score {
	return &Score{
		TotalNotes: totalNotes,
	}
}

func (s *Score) AddHit(hitType hit.HitRating) {
	switch hitType {
	case hit.Rating.Perfect:
		s.Perfect++
	case hit.Rating.Good:
		s.Good++
	case hit.Rating.Bad:
		s.Bad++
	case hit.Rating.Miss:
		s.Miss++
	}
	if hitType != hit.Rating.Miss {
		s.Combo++
		if s.Combo > s.MaxCombo {
			s.MaxCombo = s.Combo
		}
	} else {
		s.Combo = 0
	}
}

func (s *Score) GetScore() int {
	return s.Perfect*hit.Value.Perfect +
		s.Good*hit.Value.Good +
		s.Bad*hit.Value.Bad
}

func (s *Score) GetRating() ScoreRating {
	percentage := (s.Perfect + s.Good) / s.TotalNotes
	for rating, threshold := range scoreThresholds {
		if percentage >= threshold {
			return rating
		}
	}
	return Rating.F
}
