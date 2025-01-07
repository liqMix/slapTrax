package ui

import "github.com/hajimehoshi/ebiten/v2"

type Tabs struct {
	labels *UIGroup
	group  *UIGroup

	horizontal bool
}

func NewTabs() *Tabs {
	labels := NewUIGroup()
	labels.SetHorizontal()
	labels.SetTriggerOnHover(true)
	return &Tabs{
		labels:     labels,
		horizontal: true,
	}
}

func (t *Tabs) SetVertical() {
	t.labels.SetVertical()

	// lol, lmao
	currentGroup := t.group
	t.group = nil
	for _, item := range t.labels.items {
		item.Trigger()
		t.group.SetHorizontal()
	}
	t.group = currentGroup
}

func (t *Tabs) Add(label string, group *UIGroup) {
	l := NewElement()
	l.SetText(label)
	l.SetTrigger(func() {
		t.group = group
	})
	t.labels.Add(l)

	if t.group == nil {
		t.group = group
	}

	if t.horizontal {
		t.group.SetVertical()
	} else {
		t.group.SetHorizontal()
	}
}

func (t *Tabs) SetCenter(center Point, size Point, spacing float64) {
	t.labels.SetCenter(center)
	t.labels.SetSize(size)

	totalItems := len(t.labels.items)
	if totalItems == 0 {
		return
	}

	// Calculate total width using float64 throughout
	var totalItemWidth float64
	for _, item := range t.labels.items {
		totalItemWidth += TextWidth(nil, item.GetText())
	}

	totalSpacing := spacing * float64(totalItems-1)
	totalWidth := totalItemWidth + totalSpacing

	if t.horizontal {
		startX := center.X - totalWidth/2
		currentX := startX

		for _, item := range t.labels.items {
			itemWidth := TextWidth(nil, item.GetText())
			itemCenter := currentX + itemWidth/2
			item.SetCenter(Point{X: itemCenter, Y: center.Y})
			currentX += itemWidth + spacing
		}
	} else {
		startY := center.Y - totalWidth/2
		currentY := startY

		for _, item := range t.labels.items {
			itemHeight := TextWidth(nil, item.GetText()) // Assuming this is the height in vertical mode
			itemCenter := currentY + itemHeight/2
			item.SetCenter(Point{X: center.X, Y: itemCenter})
			currentY += itemHeight + spacing
		}
	}
}

func (t *Tabs) Update() {
	if len(t.labels.items) == 0 {
		return
	}

	t.labels.Update()
	if t.group != nil {
		t.group.Update()
	}
}

func (t *Tabs) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if len(t.labels.items) == 0 {
		return
	}

	t.labels.Draw(screen, opts)
	if t.group != nil {
		t.group.Draw(screen, opts)
	}
}
