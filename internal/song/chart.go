package song

import "github.com/liqmix/ebiten-holiday-2024/internal/l"

// hmm.. maybe this should be a number from 1-10 or something for more representative difficulty levels
type Difficulty int

const (
	Easy Difficulty = iota
	Hard
)

func (sd Difficulty) String() string {
	switch sd {
	case Easy:
		return l.String(l.DIFFICULTY_EASY)
	case Hard:
		return l.String(l.DIFFICULTY_HARD)
	}
	return l.String(l.UNKNOWN)
}

type Chart struct {
	Difficulty Difficulty
	Tracks     []Track
}

// Parse the chart file into a set of track and associated notes
func ParseChart(difficulty Difficulty, chartData []byte) *Chart {
	return nil
}
