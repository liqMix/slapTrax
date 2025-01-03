package types

type SongRating int

const (
	RatingSSS SongRating = iota
	RatingSS
	RatingS
	RatingA
	RatingB
	RatingF
)

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

func (r SongRating) Threshold() int {
	switch r {
	case RatingSSS:
		return 100
	case RatingSS:
		return 90
	case RatingS:
		return 80
	case RatingA:
		return 70
	case RatingB:
		return 60
	case RatingF:
		return 0
	}
	return 0
}

// scoreThresholds is the minimum percentage of perfect/good hits required to achieve a rating
var scoreThresholds map[SongRating]int = map[SongRating]int{
	RatingSSS: 100,
	RatingSS:  90,
	RatingS:   80,
	RatingA:   70,
	RatingB:   60,
	RatingF:   0,
}

type Score struct {
	TotalNotes int

	Rating  SongRating
	Perfect int
	Good    int
	Bad     int
	Miss    int

	Combo    int
	MaxCombo int

	HitRecords []*HitRecord
}

func NewScore(totalNotes int) *Score {
	return &Score{
		TotalNotes: totalNotes,
	}
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

func (s *Score) AddHit(h *HitRecord) {
	hitType := h.HitType
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
}

func (s *Score) AddMiss(n *Note) {
	s.HitRecords = append(s.HitRecords, &HitRecord{
		Note:    n,
		Diff:    0,
		HitType: Miss,
	})
	s.Miss++
	s.Combo = 0
}

func (s *Score) GetScore() int {
	return s.Perfect*Perfect.Value() +
		s.Good*Good.Value() +
		s.Bad*Bad.Value()
}

func (s *Score) GetRating() SongRating {
	percentage := (s.Perfect + s.Good) / s.TotalNotes
	for rating, threshold := range scoreThresholds {
		if percentage >= threshold {
			return rating
		}
	}
	return RatingF
}
