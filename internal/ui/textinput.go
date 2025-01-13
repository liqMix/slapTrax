package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/types"
)

type TextInput struct {
	Element

	placeholder string
	isPassword  bool
	maxLength   int
}

func NewTextInput(placeholder string) *TextInput {
	ti := &TextInput{
		Element:     *NewElement(),
		placeholder: placeholder,
		maxLength:   8,
	}
	ti.SetTrigger(func() {
		ti.SetFocused(true)
	})
	return ti
}

func (t *TextInput) SetText(s string) {
	t.text = s
}

func (t *TextInput) GetText() string { return t.text }

func (t *TextInput) SetIsPassword(p bool) { t.isPassword = p }

// func (t *TextInput) SetOnChange(f func(string)) { t.onChange = f }

func (t *TextInput) Update() {
	if !t.IsFocused() {
		return
	}

	// Handle keyboard input
	runes := input.K.Runes()
	if len(runes) > 0 && len(t.text) < t.maxLength {
		t.SetText(t.text + string(runes))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(t.text) > 0 {
		t.SetText(t.text[:len(t.text)-1])
	}

	if input.K.AreAny([]ebiten.Key{ebiten.KeyEnter, ebiten.KeyEscape}, input.JustPressed) ||
		input.M.Is(ebiten.MouseButtonLeft, input.JustPressed) {
		t.SetFocused(false)
	}
}

func (t *TextInput) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if t.hidden {
		return
	}

	t.Component.Draw(screen, opts)

	textColor := types.TrackTypeCorner.Color()
	if t.focused {
		textColor = types.TrackTypeCenter.Color()
	}

	displayText := t.text
	if t.isPassword {
		displayText = strings.Repeat("•", len(t.text))
	}
	if displayText == "" {
		displayText = t.placeholder
		textColor = color.RGBA{128, 128, 128, 255}
	}
	textOpts := &TextOptions{
		Align: t.textOptions.Align,
		Scale: t.textOptions.Scale,
		Color: textColor,
	}
	DrawTextAt(screen, displayText, t.center, textOpts, opts)

	if t.IsHovered() || t.IsFocused() {
		textOpts := &TextOptions{
			Align: t.textOptions.Align,
			Scale: t.textOptions.Scale * 2,
			Color: types.TrackTypeCorner.Color(),
		}
		if t.IsFocused() {
			textOpts.Color = types.TrackTypeCenter.Color()
		}
		size := t.size.Translate(0.05, 0)
		DrawHoverMarkersCenteredAt(screen, t.center, &size, textOpts, opts)
	}
}
