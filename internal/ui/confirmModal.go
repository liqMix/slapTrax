package ui

import (
	"image/color"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/types"
)

type ConfirmModalOptions struct {
	Title    string
	Message  string
	YesText  string
	NoText   string
	OnYes    func()
	OnNo     func()
	OnCancel func() // Called when user presses Escape
}

type ConfirmModal struct {
	options        ConfirmModalOptions
	selectedButton int // 0 = No/Cancel, 1 = Yes
}

func NewConfirmModal(opts ConfirmModalOptions) *ConfirmModal {
	// Set default texts if not provided
	if opts.YesText == "" {
		opts.YesText = "Yes"
	}
	if opts.NoText == "" {
		opts.NoText = "No"
	}
	if opts.OnCancel == nil {
		opts.OnCancel = opts.OnNo // Default cancel behavior to No
	}
	
	return &ConfirmModal{
		options:        opts,
		selectedButton: 0, // Start with No/Cancel selected for safety
	}
}

func (cm *ConfirmModal) Update() {
	// Handle navigation (only on key press, not hold)
	if input.JustActioned(input.ActionLeft) && cm.selectedButton > 0 {
		cm.selectedButton--
	}
	if input.JustActioned(input.ActionRight) && cm.selectedButton < 1 {
		cm.selectedButton++
	}
	
	// Handle selection
	if input.JustActioned(input.ActionSelect) {
		if cm.selectedButton == 0 {
			if cm.options.OnNo != nil {
				cm.options.OnNo()
			}
		} else {
			if cm.options.OnYes != nil {
				cm.options.OnYes()
			}
		}
	}
	
	// Handle cancel/back
	if input.JustActioned(input.ActionBack) {
		if cm.options.OnCancel != nil {
			cm.options.OnCancel()
		}
	}
}

func (cm *ConfirmModal) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// Semi-transparent backdrop
	backdropColor := color.RGBA{R: 0, G: 0, B: 0, A: 128}
	DrawFilledRect(screen, &Point{X: 0.5, Y: 0.5}, &Point{X: 1.0, Y: 1.0}, backdropColor)
	
	// Modal dimensions and position (normalized coordinates)
	modalCenter := &Point{X: 0.5, Y: 0.5}
	modalSize := &Point{X: 0.4, Y: 0.3}
	
	// Modal colors
	modalBg := color.RGBA{R: 40, G: 40, B: 50, A: 255}
	modalBorder := color.RGBA{R: 100, G: 100, B: 120, A: 255}
	
	// Draw modal with border
	DrawBorderedFilledRect(screen, modalCenter, modalSize, modalBg, modalBorder, 0.005)
	
	// Text options
	titleOpts := &TextOptions{
		Scale: 1.2,
		Color: types.White.C(),
	}
	messageOpts := &TextOptions{
		Scale: 1.0,
		Color: color.RGBA{R: 200, G: 200, B: 200, A: 255},
	}
	buttonOpts := &TextOptions{
		Scale: 1.0,
		Color: color.RGBA{R: 180, G: 180, B: 180, A: 255},
	}
	selectedButtonOpts := &TextOptions{
		Scale: 1.0,
		Color: color.RGBA{R: 255, G: 255, B: 100, A: 255},
	}
	
	// Draw title
	if cm.options.Title != "" {
		titlePos := &Point{X: modalCenter.X, Y: modalCenter.Y - modalSize.Y/4}
		DrawTextAt(screen, cm.options.Title, titlePos, titleOpts, opts)
	}
	
	// Draw message
	if cm.options.Message != "" {
		messagePos := &Point{X: modalCenter.X, Y: modalCenter.Y - 0.05}
		DrawTextAt(screen, cm.options.Message, messagePos, messageOpts, opts)
	}
	
	// Draw buttons
	buttonY := modalCenter.Y + modalSize.Y/4
	noButtonPos := &Point{X: modalCenter.X - 0.08, Y: buttonY}
	yesButtonPos := &Point{X: modalCenter.X + 0.08, Y: buttonY}
	
	// No button
	noOpts := buttonOpts
	if cm.selectedButton == 0 {
		noOpts = selectedButtonOpts
		// Draw selection markers
		DrawHoverMarkersCenteredAt(screen, noButtonPos, &Point{X: 0.08, Y: 0.04}, noOpts, opts)
	}
	DrawTextAt(screen, cm.options.NoText, noButtonPos, noOpts, opts)
	
	// Yes button  
	yesOpts := buttonOpts
	if cm.selectedButton == 1 {
		yesOpts = selectedButtonOpts
		// Draw selection markers
		DrawHoverMarkersCenteredAt(screen, yesButtonPos, &Point{X: 0.08, Y: 0.04}, yesOpts, opts)
	}
	DrawTextAt(screen, cm.options.YesText, yesButtonPos, yesOpts, opts)
	
	// Draw navigation hint
	hintPos := &Point{X: modalCenter.X, Y: modalCenter.Y + modalSize.Y/2 - 0.03}
	hintOpts := &TextOptions{
		Scale: 0.8,
		Color: color.RGBA{R: 150, G: 150, B: 150, A: 255},
	}
	DrawTextAt(screen, "Use arrow keys to navigate, Enter to select, Esc to cancel", hintPos, hintOpts, opts)
}