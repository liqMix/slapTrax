package types

import (
	"image/color"
	"strconv"
)

type Difficulty int

func (d Difficulty) String() string {
	return strconv.Itoa(int(d))
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
