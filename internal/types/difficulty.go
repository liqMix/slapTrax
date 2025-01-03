package types

import (
	"fmt"
	"image/color"
)

type Difficulty int

func (d Difficulty) String() string {
	return fmt.Sprintf("%d", d)
}

func (d Difficulty) Level() string {
	if d < 5 {
		return L_DIFFICULTY_EASY
	}
	if d < 8 {
		return L_DIFFICULTY_MEDIUM
	}
	if d <= 10 {
		return L_DIFFICULTY_HARD
	}
	return L_UNKNOWN
}

func (d Difficulty) Color() color.RGBA {
	if d < 5 {
		return Green
	}
	if d < 8 {
		return Yellow
	}
	if d <= 10 {
		return Red
	}
	return White
}
