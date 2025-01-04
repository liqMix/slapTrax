package render

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/tinne26/etxt"
)

var (
	linePosition = ui.Point{
		X: 0.5,
		Y: 0.65,
	}
	lineLength = 0.75
	lineThick  = float32(3)
	lineColor  = types.Gray

	noteWidth = 0.1
	noteThick = float32(10)
	noteColor = types.Yellow

	textPosition = ui.Point{
		X: 0.5,
		Y: 0.25,
	}
	textOptions = ui.TextOptions{
		Align: etxt.Center,
		Color: types.White,
		Scale: 1.0,
	}
	offsetDisplayPosition = ui.Point{
		X: 0.5,
		Y: 0.85,
	}
)

type Offset struct {
	types.BaseRenderer
	state *state.Offset
}

func NewOffsetRender(s state.State) *Offset {
	o := &Offset{state: s.(*state.Offset)}
	o.BaseRenderer.Init(o.static)

	return o
}

func (o *Offset) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	o.BaseRenderer.Draw(screen, opts)

	o.drawNote(screen)
}

func (o *Offset) static(img *ebiten.Image) {
	center := ui.Point{X: 0.5, Y: 0.5}
	size := ui.Point{X: 0.9, Y: 0.75}
	ui.DrawFilledRect(img, &center, &size, types.Gray)
	size.X -= 0.025
	size.Y -= 0.025
	ui.DrawFilledRect(img, &center, &size, types.Black)

	o.drawLine(img)

	// Draw instructions centered at top
	ui.DrawTextAt(
		img,
		"Up/Down: Audio Offset | Left/Right: Input Offset | Space: Test Hit\nA: Autoset Audio Offset | I: Autoset Input Offset | Enter: Save | Esc: Exit",
		&textPosition,
		&textOptions,
	)
}

func (o *Offset) drawNote(img *ebiten.Image) {
	lineStart := linePosition.X - lineLength/2

	noteX := lineStart + (lineLength * o.state.NoteProgress)

	// Draw the moving note vertically centered on the line
	note := ui.GetVectorPath([]*ui.Point{
		{X: noteX, Y: linePosition.Y - noteWidth},
		{X: noteX, Y: linePosition.Y + noteWidth},
	})

	if o.state.NoteProgress >= 0.49 && o.state.NoteProgress <= 0.51 {
		noteColor = types.Green
	} else {
		noteColor = types.Yellow
	}
	note.Draw(img, noteThick, noteColor)

	// Draw current offset values and hit difference at bottom
	offsetText := fmt.Sprintf(
		"Audio: %dms | Input: %dms | Last Hit: %dms",
		o.state.AudioOffset,
		o.state.InputOffset,
		o.state.HitDiff,
	)
	ui.DrawTextAt(
		img,
		offsetText,
		&offsetDisplayPosition,
		&textOptions,
	)

	// Draw auto/manual indicator
	autoText := "MANUAL"
	autoPosition := ui.Point{X: 0.5, Y: 0.95} // Bottom center
	if o.state.AutoAdjustAudio || o.state.AutoAdjustInput {
		key := "A"
		if o.state.AutoAdjustInput {
			key = "I"
		}
		autoText = fmt.Sprintf("AUTO (Press %s to toggle)", key)
		textOptions.Color = types.Red
	} else {
		textOptions.Color = types.White
	}
	ui.DrawTextAt(
		img,
		autoText,
		&autoPosition,
		&textOptions,
	)
}

// The line the note travels on (horizontal)
func (o *Offset) drawLine(img *ebiten.Image) {
	// Draw the baseline
	line := ui.GetVectorPath([]*ui.Point{
		{X: linePosition.X - lineLength/2, Y: linePosition.Y},
		{X: linePosition.X + lineLength/2, Y: linePosition.Y},
	})

	// Left hit line
	leftHitLine := ui.GetVectorPath([]*ui.Point{
		{X: linePosition.X - lineLength/2, Y: linePosition.Y - noteWidth/2},
		{X: linePosition.X - lineLength/2, Y: linePosition.Y + noteWidth/2},
	})

	// Center hit line
	centerHitLine := ui.GetVectorPath([]*ui.Point{
		{X: linePosition.X, Y: linePosition.Y - noteWidth/2},
		{X: linePosition.X, Y: linePosition.Y + noteWidth/2},
	})

	// Right hit line
	rightHitLine := ui.GetVectorPath([]*ui.Point{
		{X: linePosition.X + lineLength/2, Y: linePosition.Y - noteWidth/2},
		{X: linePosition.X + lineLength/2, Y: linePosition.Y + noteWidth/2},
	})

	// Draw all lines
	line.Draw(img, lineThick, lineColor)
	leftHitLine.Draw(img, noteThick, lineColor)
	centerHitLine.Draw(img, noteThick, lineColor)
	rightHitLine.Draw(img, noteThick, lineColor)
}
