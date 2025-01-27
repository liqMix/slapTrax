package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/l"
)

type NavText struct {
	textOpts *TextOptions
	position *Point
	text     string
}

var spacingText = "      "

func NewNavText() *NavText {
	navText := ""
	for _, a := range []input.Action{
		input.ActionBack,
		input.ActionSelect,
		input.ActionUp,
		input.ActionDown,
		input.ActionLeft,
		input.ActionRight,
	} {
		keys := a.Key()
		keyText := ""
		for i, k := range keys {
			if i > 0 {
				keyText += " or "
			}
			keyText += k.String()
		}
		navText += spacingText + l.String(a.String()) + ": " + keyText
	}
	navText += spacingText
	textOpts := GetDefaultTextOptions()
	textOpts.Scale = 1.0
	return &NavText{
		position: &Point{
			X: 0.5,
			Y: 0.95,
		},
		textOpts: textOpts,
		text:     navText,
	}
}

func (n *NavText) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	DrawTextAt(screen, n.text, n.position, n.textOpts, opts)
}
