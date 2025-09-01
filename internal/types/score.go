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
	Slap       int
	Slip       int
	Slop       int

	Combo    int
	MaxCombo int

	Early      int
	Late       int
	HitRecords []*HitRecord
	hitValue   int

	// Hold note interval tracking
	HoldIntervals    int // Total intervals across all holds
	HoldIntervalsHit int // Successfully held intervals
	holdIntervalValue int // Points per interval
}

var score *Score

func NewScore(song *Song, difficulty Difficulty) *Score {
	chart := song.Charts[difficulty]
	totalNotes := chart.TotalNotes

	// Hold notes are worth 2× regular notes: initial hit + intervals
	totalScoreUnits := totalNotes + (chart.TotalHoldNotes * 2)
	
	score = &Score{
		Song:       song,
		Difficulty: difficulty,
		TotalNotes: totalNotes,
		HitRecords: make([]*HitRecord, 0, totalNotes),
		hitValue:   MaxScore / totalScoreUnits,
	}
	return score
}

func (s *Score) Reset() {
	s.Slap = 0
	s.Slip = 0
	s.Slop = 0

	s.Combo = 0
	s.MaxCombo = 0
	
	s.HoldIntervals = 0
	s.HoldIntervalsHit = 0
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
	if hitType == None || hitType == Slop {
		return
	}
	s.HitRecords = append(s.HitRecords, h)
	switch hitType {
	case Slap:
		s.Slap++
	case Slip:
		s.Slip++
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
		HitRating: Slop,
	}
	s.HitRecords = append(s.HitRecords, record)
	s.Slop++
	s.Combo = 0
}

func (s *Score) GetAccuracy() float64 {
	return float64(s.Slap) / float64(s.TotalNotes)
}

// AddHoldInterval adds scoring for hold note intervals
func AddHoldInterval(hit bool) {
	AddHoldIntervalWithCombo(hit, true)
}

// AddHoldIntervalWithCombo adds scoring for hold note intervals with optional combo breaking
func AddHoldIntervalWithCombo(hit bool, breakComboOnMiss bool) {
	s := score
	s.HoldIntervals++
	
	if hit {
		s.HoldIntervalsHit++
		// Award interval points - each interval is worth hitValue/intervalCount
		s.TotalScore += s.holdIntervalValue
	} else {
		// Only break combo if this is a definitive miss (not temporary release)
		if breakComboOnMiss {
			s.Combo = 0
		}
	}
}
