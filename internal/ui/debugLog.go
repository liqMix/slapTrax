package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/logger"
)

const visibleLines = 10

type DebugLog struct {
	messages []*logger.Message
}

func NewDebugLog() *DebugLog {
	return &DebugLog{}
}

func (d *DebugLog) Update() {
	pending := logger.GetMessages()
	if len(pending) == 0 {
		return
	}
	d.messages = append(d.messages, pending...)
}

func (d *DebugLog) Draw(screen *ebiten.Image) {
	// visible := int(math.Min(float64(visibleLines), float64(len(d.messages))))
	// if visible == 0 {
	// 	return
	// }
	// visibleMessages := d.messages[len(d.messages)-visible:]

	// p := &Point{X: 0.5, Y: 0.1}
	// height := TextHeight()
	// for _, m := range visibleMessages {
	// 	if m == nil || m.Message == "" {
	// 		continue
	// 	}
	// 	DrawTextAt(screen, m.Message, p, &TextOptions{Align: etxt.Left, Scale: 0.5, Color: m.Color})
	// 	p.Y += height
	// }
}
