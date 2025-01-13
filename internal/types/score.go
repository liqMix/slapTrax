package types

import "image/color"

type SongRating int

const MaxScore = 100000

const (
	RatingSSS SongRating = iota
	RatingSS
	RatingS
	RatingA
	RatingB
	RatingF
)

var Ratings = []SongRating{RatingSSS, RatingSS, RatingS, RatingA, RatingB, RatingF}

func (r SongRating) String() string {
	switch r {
	case RatingSSS:
		return "SSS"
	case RatingSS:
		return "SS"
	case RatingS:
		return "S"
	case RatingA:
		return "A"
	case RatingB:
		return "B"
	case RatingF:
		return "F"
	}
	return "F"
}

func (r SongRating) Color() color.RGBA {
	switch r {
	case RatingSSS:
		return Red.C()
	case RatingSS:
		return Orange.C()
	case RatingS:
		return Yellow.C()
	case RatingA:
		return LightBlue.C()
	case RatingB:
		return Green.C()
	case RatingF:
		return Gray.C()
	}
	return Gray.C()
}

func (r SongRating) Threshold() int {
	switch r {
	case RatingSSS:
		return MaxScore
	case RatingSS:
		return MaxScore * 0.9
	case RatingS:
		return MaxScore * 0.8
	case RatingA:
		return MaxScore * 0.7
	case RatingB:
		return MaxScore * 0.6
	case RatingF:
		return 0
	}
	return 0
}

func GetSongRating(score int) SongRating {
	for _, r := range Ratings {
		if score >= r.Threshold() {
			return r
		}
	}

	return RatingF
}

type Score struct {
	Song       *Song
	Difficulty Difficulty
	TotalNotes int

	TotalScore int
	Rating     SongRating
	Perfect    int
	Good       int
	Bad        int
	Miss       int

	Combo    int
	MaxCombo int

	Early      int
	Late       int
	HitRecords []*HitRecord
	hitValue   int
}

var score *Score

func NewScore(song *Song, difficulty Difficulty) *Score {
	chart := song.Charts[difficulty]
	totalNotes := chart.TotalNotes

	score = &Score{
		Song:       song,
		Difficulty: difficulty,
		TotalNotes: totalNotes,
		HitRecords: make([]*HitRecord, 0, totalNotes),
		hitValue:   MaxScore / (totalNotes + chart.TotalHoldNotes),
	}
	return score
}

func (s *Score) Reset() {
	s.Perfect = 0
	s.Good = 0
	s.Bad = 0
	s.Miss = 0

	s.Combo = 0
	s.MaxCombo = 0
}

func (s *Score) GetLastHitRecord() *HitRecord {
	if len(s.HitRecords) == 0 {
		return nil
	}
	return s.HitRecords[len(s.HitRecords)-1]
}

func AddHit(h *HitRecord) {
	s := score
	hitType := h.HitRating
	if hitType == None || hitType == Miss {
		return
	}
	s.HitRecords = append(s.HitRecords, h)
	switch hitType {
	case Perfect:
		s.Perfect++
	case Good:
		s.Good++
	case Bad:
		s.Bad++
	}
	s.Combo++
	if s.Combo > s.MaxCombo {
		s.MaxCombo = s.Combo
	}

	timing := h.HitTiming
	if timing == HitTimingEarly {
		s.Early++
	} else if timing == HitTimingLate {
		s.Late++
	}
	s.TotalScore += int(hitType.Value() * float64(s.hitValue))
}

func (s *Score) AddMiss(n *Note) {
	record := &HitRecord{
		Note:      n,
		HitRating: Miss,
	}
	s.HitRecords = append(s.HitRecords, record)
	s.Miss++
	s.Combo = 0
}

func (s *Score) GetAccuracy() float64 {
	return float64(s.Perfect) / float64(s.TotalNotes)
}
