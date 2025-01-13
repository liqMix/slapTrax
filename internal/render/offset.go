package render

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/tinne26/etxt"
)

var (
	offsetDisplayPosition = ui.Point{
		X: 0.5,
		Y: 0.8,
	}
	lineThick    = float32(3)
	lineColor    = types.Gray.C()
	linePosition = ui.Point{
		X: 0.5,
		Y: 0.65,
	}
	lineLength = 0.5

	noteWidth   = 0.1
	noteThick   = float32(10)
	textOptions = ui.TextOptions{
		Align: etxt.Center,
		Color: types.White.C(),
		Scale: 1,
	}
)

type Offset struct {
	display.BaseRenderer
	state *state.Offset
}

func NewOffsetRender(s state.State) *Offset {
	o := &Offset{state: s.(*state.Offset)}
	o.BaseRenderer.Init(o.static)

	return o
}

func (o *Offset) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	o.BaseRenderer.Draw(screen, opts)
	o.drawNote(screen, opts)
}

func (o *Offset) static(img *ebiten.Image, opts *ebiten.DrawImageOptions) {
	textPosition := ui.Point{
		X: 0.33,
		Y: 0.2,
	}
	inputTextPosition := ui.Point{
		X: 0.75,
		Y: 0.25,
	}

	center := ui.Point{X: 0.5, Y: 0.5}
	size := ui.Point{X: 0.9, Y: 0.75}
	ui.DrawFilledRect(img, &center, &size, types.Gray.C())
	size.X -= 0.025
	size.Y -= 0.025
	ui.DrawFilledRect(img, &center, &size, types.Black.C())

	o.drawLine(img)
	text := l.String(l.OFFSET_INSTRUCTIONS)
	for _, line := range strings.Split(text, "\n") {
		bold := false
		if len(line) >= 3 {
			if line[0:3] == "<b>" {
				bold = true
				line = line[3:]
			}
		}

		textPosition.Y += 0.02
		ui.DrawTextAt(
			img,
			line,
			&textPosition,
			&textOptions,
			opts,
		)
		if bold {
			textPosition.X += 0.001
			ui.DrawTextAt(
				img,
				line,
				&textPosition,
				&textOptions,
				opts,
			)
			textPosition.X -= 0.001
		}
	}
	inputText := l.String(l.OFFSET_INPUT)
	for _, line := range strings.Split(inputText, "\n") {
		inputTextPosition.Y += 0.02
		ui.DrawTextAt(
			img,
			line,
			&inputTextPosition,
			&textOptions,
			opts,
		)
	}
}

func (o *Offset) drawNote(img *ebiten.Image, opts *ebiten.DrawImageOptions) {

	lineStart := linePosition.X - lineLength/2

	noteX := lineStart + (lineLength * o.state.NoteProgress)

	// Draw the moving note vertically centered on the line
	note := ui.GetVectorPath([]*ui.Point{
		{X: noteX, Y: linePosition.Y - noteWidth},
		{X: noteX, Y: linePosition.Y + noteWidth},
	})

	var noteColor color.RGBA
	if o.state.NoteProgress >= 0.5-o.state.CenterWindow && o.state.NoteProgress <= 0.5+o.state.CenterWindow {
		noteColor = types.Green.C()
	} else {
		noteColor = ui.CenterTrackColor()
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
		opts,
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
	leftHitLine.Draw(img, noteThick, ui.CornerTrackColor())
	centerHitLine.Draw(img, noteThick, lineColor)
	rightHitLine.Draw(img, noteThick, ui.CornerTrackColor())
}
