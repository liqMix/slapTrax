package ui

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type SongSelector struct {
	*UIGroup
}

const maxDisplayedItems = 10

func NewSongSelector() *SongSelector {
	return &SongSelector{
		UIGroup: NewUIGroup(),
	}
}

// We'll draw the selected item centered at position.
// Previous and next items will be drawn above and below the selected item with opacity on their distance
func (s *SongSelector) Update() {
	s.UIGroup.Update()
	if s.GetCenter() == nil {
		return
	}
	position := s.GetCenter()
	yOffset := 0.05
	xOffset := 0.025
	opacity := 0.8
	min := s.currentIdx - maxDisplayedItems/2
	if min < 0 {
		min = 0
	}
	max := s.currentIdx + maxDisplayedItems/2
	if max >= len(s.items) {
		max = len(s.items) - 1
	}
	for i := min; i <= max; i++ {
		item := s.items[i]
		if i == s.currentIdx {
			item.SetCenter(Point{
				X: position.X,
				Y: position.Y,
			})
			item.SetOpacity(1.0)
		} else {
			distance := float64(i - s.currentIdx)
			absDist := math.Abs(distance)
			item.SetCenter(Point{
				X: position.X + xOffset*absDist,
				Y: position.Y + yOffset*float64(distance),
			})
			item.SetOpacity(opacity * absDist)
		}
	}
}

func (s *SongSelector) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	for _, item := range s.items {
		item.Draw(screen, opts)
	}
}
