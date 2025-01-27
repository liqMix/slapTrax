package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/input"
)

type ValueElement struct {
	Element

	label        string
	getValueText func() string
}

func NewValueElement() *ValueElement {
	return &ValueElement{
		Element: *NewElement(),
	}
}

func (b *ValueElement) SetLabel(label string) {
	b.label = label
}

func (b *ValueElement) Refresh() {
	b.SetText(fmt.Sprintf("%s: %s", b.label, b.getValueText()))
}
func (b *ValueElement) SetGetValueText(getValueText func() string) {
	b.getValueText = getValueText
	b.Refresh()
}

func (b *ValueElement) SetTrigger(trigger func()) {
	b.Element.SetTrigger(func() {
		trigger()
		b.Refresh()
	})
}

type KeyboardInputElement struct {
	ValueElement

	editing bool
	input   string
}

func NewKeyboardInputElement() *KeyboardInputElement {
	kie := KeyboardInputElement{
		ValueElement: *NewValueElement(),
	}
	kie.Element.SetTrigger(func() {
		kie.editing = !kie.editing
	})
	return &kie
}

func (kie *KeyboardInputElement) SetTrigger(func()) {}

func (kie *KeyboardInputElement) Update() {
	if kie.editing {
		justPressed := input.K.Get(input.JustPressed)
		if len(justPressed) > 0 {
			key := justPressed[0]
			k := key.String()
			if len(k) == 1 {
				kie.input += k
			} else if key == ebiten.KeyBackspace {
				kie.input = kie.input[:len(kie.input)-1]
			} else if key == ebiten.KeyEnter {
				kie.editing = false
			} else if key == ebiten.KeyEscape || key == ebiten.KeyF1 {
				kie.editing = false
				kie.input = ""
			}
		}
	}
}
