package internal

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/assets"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/cache"
	"github.com/liqmix/slaptrax/internal/debug"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/state"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/liqmix/slaptrax/internal/user"
	"github.com/tinne26/etxt"
)

const (
	maxStateStackSize = 10
	startingTicks     = 200
)

func getStartingState() *RenderState {
	// startingState := types.GameStatePlay
	// song := types.GetAllSongs()[0]
	// diff := song.GetDifficulties()[0]
	// startingArgs := &state.PlayArgs{
	// Song:       song,
	// Difficulty: diff,
	// }

	// startingState := types.GameStateOffset
	startingState := types.GameStateTitle
	startingArgs := interface{}(nil)
	// startingState := types.GameStateResult
	// startingArgs := &state.ResultStateArgs{
	// 	Score: &types.Score{
	// 		TotalScore: types.MaxScore,
	// 		MaxCombo:   100,
	// 		Difficulty: 7,
	// 		Song:       types.GetAllSongs()[0],
	// 		Rating:     types.RatingS,
	// 		TotalNotes: 100,
	// 		Perfect:    1000, Good: 0, Bad: 0, Miss: 0,
	// 		Combo:      100,
	// 		HitRecords: []*types.HitRecord{},
	// 	},
	// }

	return GetState(startingState, startingArgs)
}

type Game struct {
	debugster *debug.Debugster

	started         bool
	startTicks      int64
	stateStack      []*RenderState
	currentState    *RenderState
	navText         *ui.NavText
	loadingTextOpts *ui.TextOptions
	userHeader      *ui.UserProfile
	background      *ebiten.Image

	// Header caching
	cachedFullHeader   *ebiten.Image
	cachedSmallHeader  *ebiten.Image
	lastHeaderCacheKey string
	needsHeaderRefresh bool
}

func NewGame() *Game {
	return &Game{
		started:            false,
		startTicks:         0,
		currentState:       getStartingState(),
		debugster:          debug.NewDebugster(),
		navText:            ui.NewNavText(),
		userHeader:         ui.NewUserProfile(),
		needsHeaderRefresh: true,
		loadingTextOpts: &ui.TextOptions{
			Align: etxt.Center,
			Scale: 1.5,
			Color: types.White.C(),
		},
		background: assets.GetImage("background.png"),
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	panic("nope")
}

func (g *Game) LayoutF(displayWidth, displayHeight float64) (float64, float64) {
	display.Window.SetDisplaySize(displayWidth, displayHeight)
	return displayWidth, displayHeight
}

func (g *Game) handleStateTransition(nextState types.GameState, nextArgs interface{}) error {
	if nextState == types.GameStateBack {
		audio.PlaySFX(audio.SFXBack)
		return g.popState()
	}

	next := GetState(nextState, nextArgs)
	if next.state.Floats() {
		if len(g.stateStack) >= maxStateStackSize {
			return fmt.Errorf("state stack overflow: max size %d reached", maxStateStackSize)
		}
		g.currentState.Freeze()
		g.stateStack = append(g.stateStack, g.currentState)
	} else {
		g.stateStack = nil
	}
	g.currentState = next
	return nil
}

func (g *Game) popState() error {
	if len(g.stateStack) == 0 {
		return errors.New("cannot pop state: stack is empty")
	}
	g.currentState = g.stateStack[len(g.stateStack)-1]
	g.stateStack = g.stateStack[:len(g.stateStack)-1]
	g.currentState.Unfreeze()
	return nil
}

func (g *Game) Update() error {
	if !g.started {
		g.startTicks++
		if g.startTicks >= startingTicks {
			g.started = true
		}
	}

	audio.Update()
	input.Update()

	if input.JustActioned(input.ActionToggleDebug) {
		g.debugster.Toggle()
	}

	if g.currentState == nil {
		return nil
	}

	gs := g.currentState.state
	if gs == nil {
		return nil
	}

	action := gs.CheckActions()
	if action != input.ActionUnknown {
		sfx := audio.ActionSFX(action)
		if sfx != audio.SFXNone {
			audio.PlaySFX(sfx)
		}
	}

	if err := gs.Update(); err != nil {
		return err
	}

	if gs.HasNextState() {
		nextState, nextArgs := gs.GetNextState()
		gs.SetNextState(types.GameStateNone, nil)

		if err := g.handleStateTransition(nextState, nextArgs); err != nil {
			return fmt.Errorf("state transition failed: %w", err)
		}
	}

	g.userHeader.Update()
	g.debugster.Update()
	return nil
}

// Create canvas at render size
func (g *Game) GetCanvasImage() *ebiten.Image {
	canvas, ok := cache.Image.Get("canvas")
	if !ok {
		canvas = display.NewRenderImage()
		cache.Image.Set("canvas", canvas)
	}
	canvas.Clear()
	return canvas
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.Fill(color.Black)
	canvas := g.GetCanvasImage()

	bgOpts := &ebiten.DrawImageOptions{}

	// Scale background to fit render size
	renderWidth, renderHeight := display.Window.RenderSize()
	if g.background != nil {
		scaleX := float64(renderWidth) / float64(g.background.Bounds().Dx())
		scaleY := float64(renderHeight) / float64(g.background.Bounds().Dy())
		bgOpts.GeoM.Scale(scaleX, scaleY)
	}

	bgOpts.ColorScale.Scale(0.25, 0.25, 0.25, 0.25)
	canvas.DrawImage(g.background, bgOpts)

	if g.currentState != nil {
		for i, s := range g.stateStack {
			opts := &ebiten.DrawImageOptions{}
			a := float32(0.25) * float32(i+1)
			opts.ColorScale.Scale(a, a, a, a)
			s.Draw(canvas, opts)
		}
		g.currentState.Draw(canvas, nil)

		// Draw nav action bar if navigable
		if g.currentState.state.IsNavigable() {
			g.navText.Draw(canvas, nil)
		}
	}

	// Draw unified header system
	g.drawUnifiedHeader(canvas)
	opts := display.Window.GetScreenDrawOptions()
	if !g.started {
		scale := float32(g.startTicks) / float32(startingTicks)
		opts.ColorScale.ScaleAlpha(scale)
	}
	screen.DrawImage(canvas, opts)

	if logger.IsDebugEnabled() {
		g.debugster.Draw(screen)
	}
}

// generateHeaderCacheKey creates a cache key based on current game state and content
func (g *Game) generateHeaderCacheKey() string {
	playState := g.getPlayState()
	isPlayState := playState != nil
	isPaused := len(g.stateStack) > 0
	isEdgePlayArea := user.S().EdgePlayArea

	key := fmt.Sprintf("play:%v|pause:%v|edge:%v", isPlayState, isPaused, isEdgePlayArea)

	if playState != nil {
		// Include song-specific data that affects header (excluding score which changes frequently)
		key += fmt.Sprintf("|song:%s", playState.Song.Title)
	}

	// Include user data that affects header
	key += fmt.Sprintf("|user:%s|rank:%d", user.Current().Username, user.Current().Rank)

	return key
}

// invalidateHeaderCache marks the header cache as needing refresh
func (g *Game) invalidateHeaderCache() {
	g.needsHeaderRefresh = true
}

// drawUnifiedHeader renders all header components consistently based on game state
func (g *Game) drawUnifiedHeader(canvas *ebiten.Image) {
	// Check if we need to refresh cached headers
	currentKey := g.generateHeaderCacheKey()
	if g.needsHeaderRefresh || g.lastHeaderCacheKey != currentKey {
		g.refreshHeaderCaches()
		g.lastHeaderCacheKey = currentKey
		g.needsHeaderRefresh = false
	}

	// Detect current state type
	playState := g.getPlayState()
	isPlayState := playState != nil
	isMenuState := !isPlayState

	// Draw appropriate cached header
	if isMenuState {
		if g.cachedSmallHeader != nil {
			canvas.DrawImage(g.cachedSmallHeader, &ebiten.DrawImageOptions{})
		}
	} else {
		if g.cachedFullHeader != nil {
			canvas.DrawImage(g.cachedFullHeader, &ebiten.DrawImageOptions{})
		}
		
		// Draw score separately (not cached since it changes frequently)
		if playState != nil {
			isPaused := len(g.stateStack) > 0
			isEdgePlayArea := user.S().EdgePlayArea
			
			headerOpts := &ebiten.DrawImageOptions{}
			if isPaused {
				headerOpts.ColorScale.Scale(0.25, 0.25, 0.25, 0.25)
			} else if isEdgePlayArea {
				headerOpts.ColorScale.Scale(0.4, 0.4, 0.4, 0.6)
			}
			
			g.drawScore(canvas, headerOpts, playState, 0.1)
		}
	}
}

// refreshHeaderCaches regenerates both cached header images
func (g *Game) refreshHeaderCaches() {
	displayW, displayH := display.Window.DisplaySize()

	// Create fresh header images
	g.cachedFullHeader = ebiten.NewImage(int(displayW), int(displayH))
	g.cachedSmallHeader = ebiten.NewImage(int(displayW), int(displayH))

	// Generate full header (play state)
	g.renderFullHeaderToCache()

	// Generate small header (menu state)
	g.renderSmallHeaderToCache()
}

// renderFullHeaderToCache renders the full header to cache
func (g *Game) renderFullHeaderToCache() {
	playState := g.getPlayState()
	isPaused := len(g.stateStack) > 0
	isEdgePlayArea := user.S().EdgePlayArea

	// Header constants
	const (
		headerHeight  = 0.2
		headerCenterY = 0.1
		headerRight   = 1.0
		equalMargin   = 0.025
		artSize       = headerHeight * 0.4
		textMargin    = 0.02
	)

	// Calculate content dimming for pause
	headerOpts := &ebiten.DrawImageOptions{}
	if isPaused {
		// When paused, dim content and borders (fixed amount, not based on stack depth)
		headerOpts.ColorScale.Scale(0.25, 0.25, 0.25, 0.25)
	} else if isEdgePlayArea {
		// In edge play area mode, apply subtle dimming for less distraction
		headerOpts.ColorScale.Scale(0.4, 0.4, 0.4, 0.6)
	}

	// Draw full header background with dimming applied
	if !isEdgePlayArea {
		g.drawFullHeaderBackground(g.cachedFullHeader, headerOpts)
	}

	// Draw all header content to cache (excluding score which changes frequently)
	if playState != nil {
		g.drawSongDetails(g.cachedFullHeader, headerOpts, playState, headerHeight, headerCenterY, headerRight, equalMargin, artSize, textMargin)
	}
	g.userHeader.Draw(g.cachedFullHeader, headerOpts)
}

// renderSmallHeaderToCache renders the small header to cache
func (g *Game) renderSmallHeaderToCache() {
	// Menu states: Small header box with user profile only
	g.drawSmallHeaderBox(g.cachedSmallHeader)
	g.userHeader.Draw(g.cachedSmallHeader, nil) // No dimming in menu
}

// drawFullHeaderBackground draws the full-width header panel and border for play state
func (g *Game) drawFullHeaderBackground(canvas *ebiten.Image, opts *ebiten.DrawImageOptions) {
	const (
		headerWidth   = 1.0
		headerHeight  = 0.2
		headerCenterX = 0.5
		headerCenterY = 0.1
		headerBottom  = 0.2
		borderHeight  = 0.003
	)

	center := ui.Point{X: headerCenterX, Y: headerCenterY}
	size := ui.Point{X: headerWidth, Y: headerHeight}

	// Draw main background panel (always full brightness)
	ui.DrawFilledRect(canvas, &center, &size, color.RGBA{R: 0, G: 0, B: 0, A: 120})

	// Draw bottom border with dimming applied
	borderColor := ui.CornerTrackColor()
	if opts != nil {
		borderColor = ui.ApplyColorScale(borderColor, opts.ColorScale.R(), opts.ColorScale.G(), opts.ColorScale.B(), opts.ColorScale.A())
	}
	borderCenter := ui.Point{X: headerCenterX, Y: headerBottom - (borderHeight / 2)}
	borderSize := ui.Point{X: headerWidth, Y: borderHeight}
	ui.DrawFilledRect(canvas, &borderCenter, &borderSize, borderColor)
}

// drawSmallHeaderBox draws a small header box for menu states (user profile area only)
func (g *Game) drawSmallHeaderBox(canvas *ebiten.Image) {
	const (
		equalMargin  = 0.025
		headerHeight = 0.2
		borderWidth  = 0.003

		// Box extends past screen edges - negative margin to go off-screen
		boxWidth  = 0.2 + equalMargin              // Increased width to better fit content
		boxHeight = headerHeight*1.0 + equalMargin // Full header height for better fit

		// Position so box extends past left and top edges
		boxCenterX = boxWidth/2 - equalMargin  // Shift left to go past edge
		boxCenterY = boxHeight/2 - equalMargin // Shift up to go past edge
	)

	center := ui.Point{X: boxCenterX, Y: boxCenterY}
	size := ui.Point{X: boxWidth, Y: boxHeight}

	// Draw main box background
	ui.DrawFilledRect(canvas, &center, &size, color.RGBA{R: 0, G: 0, B: 0, A: 120})

	// Draw right border (vertical)
	rightBorderCenter := ui.Point{
		X: boxCenterX + boxWidth/2 - borderWidth/2,
		Y: boxCenterY,
	}
	rightBorderSize := ui.Point{X: borderWidth, Y: boxHeight}
	ui.DrawFilledRect(canvas, &rightBorderCenter, &rightBorderSize, ui.CornerTrackColor())

	// Draw bottom border (horizontal)
	bottomBorderCenter := ui.Point{
		X: boxCenterX,
		Y: boxCenterY + boxHeight/2 - borderWidth/2,
	}
	bottomBorderSize := ui.Point{X: boxWidth, Y: borderWidth}
	ui.DrawFilledRect(canvas, &bottomBorderCenter, &bottomBorderSize, ui.CornerTrackColor())
}

// drawSongDetails draws song information on the right side of header
func (g *Game) drawSongDetails(canvas *ebiten.Image, opts *ebiten.DrawImageOptions, playState *state.Play, headerHeight, headerCenterY, headerRight, equalMargin, artSize, textMargin float64) {
	// Album art
	art := ui.NewElement()
	art.SetDisabled(true)
	art.SetSize(ui.Point{X: artSize, Y: artSize})
	art.SetCenter(ui.Point{
		X: headerRight - equalMargin - (artSize / 2),
		Y: headerCenterY,
	})
	art.SetImage(playState.Song.Art)

	artGroup := ui.NewUIGroup()
	artGroup.SetDisabled(true)
	artGroup.SetCenter(ui.Point{X: headerRight - equalMargin - (artSize / 2), Y: headerCenterY})
	artGroup.SetSize(ui.Point{X: artSize, Y: artSize})
	artGroup.Add(art)
	artGroup.Draw(canvas, opts)

	// Song text
	textOpts := ui.GetDefaultTextOptions()
	textOpts.Align = etxt.Right
	textOpts.Scale = 1.3

	// Apply text dimming if in edge play area mode (since opts won't affect text color)
	if user.S().EdgePlayArea && len(g.stateStack) == 0 {
		textOpts.Color = color.RGBA{
			R: uint8(float32(textOpts.Color.R) * 0.4),
			G: uint8(float32(textOpts.Color.G) * 0.4),
			B: uint8(float32(textOpts.Color.B) * 0.4),
			A: uint8(float32(textOpts.Color.A) * 0.6),
		}
	}

	// Song title (bold)
	textCenter := &ui.Point{
		X: headerRight - equalMargin - artSize - textMargin,
		Y: headerCenterY - 0.03,
	}
	ui.DrawTextAt(canvas, playState.Song.Title, textCenter, textOpts, opts)
	textCenter.X += 0.001
	ui.DrawTextAt(canvas, playState.Song.Title, textCenter, textOpts, opts)
	textCenter.X -= 0.001

	// Artist
	textOpts.Scale = 1.0
	textCenter.Y += 0.04
	ui.DrawTextAt(canvas, playState.Song.Artist, textCenter, textOpts, opts)

	// Album
	textCenter.Y += 0.03
	textOpts.Scale = 0.9
	ui.DrawTextAt(canvas, playState.Song.Album, textCenter, textOpts, opts)
}

// drawScore draws the score in the center of the header
func (g *Game) drawScore(canvas *ebiten.Image, opts *ebiten.DrawImageOptions, playState *state.Play, headerCenterY float64) {
	center := &ui.Point{X: 0.5, Y: headerCenterY}
	textOpts := &ui.TextOptions{
		Align: etxt.Center,
		Scale: 4.0,
		Color: types.White.C(),
	}

	// Apply text dimming if in edge play area mode (since opts won't affect text color)
	if user.S().EdgePlayArea && len(g.stateStack) == 0 {
		textOpts.Color = color.RGBA{
			R: uint8(float32(textOpts.Color.R) * 0.4),
			G: uint8(float32(textOpts.Color.G) * 0.4),
			B: uint8(float32(textOpts.Color.B) * 0.4),
			A: uint8(float32(textOpts.Color.A) * 0.6),
		}
	}

	totalScore := fmt.Sprintf("%d", playState.Score.TotalScore)
	ui.DrawTextAt(canvas, totalScore, center, textOpts, opts)
}

// getPlayState returns the current Play state if available (checks both current state and state stack for paused Play)
func (g *Game) getPlayState() *state.Play {
	// Check current state first
	if g.currentState != nil && g.currentState.state != nil {
		if playState, ok := g.currentState.state.(*state.Play); ok {
			return playState
		}
	}

	// Check state stack (for when Play is frozen under pause)
	for _, s := range g.stateStack {
		if s != nil && s.state != nil {
			if playState, ok := s.state.(*state.Play); ok {
				return playState
			}
		}
	}

	return nil
}
