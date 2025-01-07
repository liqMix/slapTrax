package types

import (
	"fmt"
	"image/color"

	"github.com/liqmix/ebiten-holiday-2024/internal/l"
)

type Difficulty int

func (d Difficulty) String() string {
	return fmt.Sprintf("%d", d)
}

func (d Difficulty) Level() string {
	if d < 5 {
		return l.DIFFICULTY_EASY
	}
	if d < 8 {
		return l.DIFFICULTY_MEDIUM
	}
	if d <= 10 {
		return l.DIFFICULTY_HARD
	}
	return l.UNKNOWN
}

func (d Difficulty) Color() color.RGBA {
	color := White
	if d < 5 {
		color = Green
	} else if d < 8 {
		color = Yellow
	} else if d <= 10 {
		color = Red
	}
	return color.C()
}
