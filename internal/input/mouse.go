package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
)

type mouse struct {
	sX, sY int
	cX, cY float64

	// Buttons
	left, right mousebutton
}
type mousebutton struct {
	b    ebiten.MouseButton
	s    Status
	held int
}

func newMousebutton(b ebiten.MouseButton) mousebutton {
	return mousebutton{
		b: b,
		s: None,
	}
}

func newMouse() mouse {
	return mouse{
		left:  newMousebutton(ebiten.MouseButtonLeft),
		right: newMousebutton(ebiten.MouseButtonRight),
	}
}

func (mb *mousebutton) update(pressed bool) {
	if pressed {
		switch mb.s {
		case None, JustReleased:
			mb.s = JustPressed
		case JustPressed:
			mb.s = Held
		case Held:
			mb.held++
		}
	} else {
		switch mb.s {
		case JustPressed, Held:
			mb.s = JustReleased
			mb.held = 0
		case JustReleased:
			mb.s = None
			mb.held = 0
		}
	}
}

func (m *mouse) update() {
	// Screen position
	m.sX, m.sY = ebiten.CursorPosition()

	// Canvas position
	m.cX, m.cY = display.Window.CanvasPosition(float64(m.sX), float64(m.sY))

	// Buttons
	leftP := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	rightP := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	m.left.update(leftP)
	m.right.update(rightP)
}

func (m *mouse) Position() (float64, float64) {
	return m.cX, m.cY
}

// func (m *mouse) ScreenPosition() (int, int) {
// 	return m.sX, m.sY
// }

func (m *mouse) Is(b ebiten.MouseButton, s Status) bool {
	switch b {
	case ebiten.MouseButtonLeft:
		return m.left.s == s
	case ebiten.MouseButtonRight:
		return m.right.s == s
	}
	return false
}

func (m *mouse) InBounds(x, y, w, h float64) bool {
	return m.cX >= x && m.cX <= x+w && m.cY >= y && m.cY <= y+h
}

func (m *mouse) HoldTime(b ebiten.MouseButton) int {
	switch b {
	case ebiten.MouseButtonLeft:
		return m.left.held
	case ebiten.MouseButtonRight:
		return m.right.held
	}
	return 0
}
