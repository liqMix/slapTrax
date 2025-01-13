package ui

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type SongSelector struct {
	*UIGroup

	minIdx int
	maxIdx int
}

const maxDisplayedItems = 15

func NewSongSelector() *SongSelector {
	return &SongSelector{
		UIGroup: NewUIGroup(),
	}
}

func (s *SongSelector) Update() {
	s.UIGroup.Update()
	if s.GetCenter() == nil {
		return
	}
	position := s.GetCenter()
	yOffset := 0.05
	xOffset := 0.025
	opacity := 1.5
	min := s.currentIdx - maxDisplayedItems/2
	if min < 0 {
		min = 0
	}
	s.minIdx = min
	max := s.currentIdx + maxDisplayedItems/2
	if max >= len(s.items) {
		max = len(s.items) - 1
	}
	s.maxIdx = max

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
	for i := s.minIdx; i <= s.maxIdx; i++ {
		item := s.items[i]
		item.Draw(screen, opts)
	}
}
