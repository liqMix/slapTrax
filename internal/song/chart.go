package song

import "github.com/liqmix/ebiten-holiday-2024/internal/l"

type Difficulty int

const (
	Easy Difficulty = iota
	Hard
)

func (sd Difficulty) String() string {
	switch sd {
	case Easy:
		return l.String(l.EASY)
	case Hard:
		return l.String(l.HARD)
	}
	return l.String(l.UNKNOWN)
}

type Chart struct {
	Difficulty Difficulty
	Tracks     []Track
}
