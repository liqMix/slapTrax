package state

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/cache"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/l"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/liqmix/slaptrax/internal/user"
)

var (
	settingsCenter = ui.Point{X: 0.5, Y: 0.5}
	settingsSize   = ui.Point{X: 0.5, Y: 0.8}
	tabStart       = ui.Point{X: 0.25, Y: 0.15}
	tabOffset      = 0.1

	optionsStart  = ui.Point{X: 0.5, Y: 0.28}
	optionsOffset = 0.1
)

type Settings struct {
	types.BaseGameState

	clearingCache bool
	tabs          *ui.Tabs
}

func NewSettingsState() *Settings {
	s := &Settings{}
	s.Refresh()
	return s
}

func (s *Settings) Refresh() {
	s.tabs = ui.NewTabs()

	tabCenter := tabStart
	var g *ui.UIGroup

	// Gameplay Tab
	g = ui.NewUIGroup()
	s.createGameplayOptions(g)
	s.tabs.Add(l.String(l.SETTINGS_GAME), g)
	tabCenter.X += tabOffset

	// Graphics Tab
	g = ui.NewUIGroup()
	s.createGraphicsOptions(g)
	s.tabs.Add(l.String(l.SETTINGS_GFX), g)
	tabCenter.X += tabOffset

	// Audio Tab
	g = ui.NewUIGroup()
	s.createAudioOptions(g)
	s.tabs.Add(l.String(l.SETTINGS_AUDIO), g)
	tabCenter.X += tabOffset

	// Accessibility Tab
	g = ui.NewUIGroup()
	g.SetCenter(tabCenter)
	s.createAccessOptions(g)
	s.tabs.Add(l.String(l.SETTINGS_ACCESS), g)

	s.tabs.SetCenter(
		ui.Point{X: 0.5, Y: 0.15},
		ui.Point{X: 0.15, Y: 0.2},
		0.07,
	)
}

func (s *Settings) createGraphicsOptions(group *ui.UIGroup) {
	optionPos := ui.Point{
		X: optionsStart.X,
		Y: optionsStart.Y,
	}

	// screenSizeButton := ui.NewValueElement()
	// if user.S().Fullscreen {
	// 	screenSizeButton.SetDisabled(true)
	// }

	fixedRender := ui.NewValueElement()
	fixedRender.SetLabel(l.String(l.SETTINGS_GFX_FIXEDRENDER))
	fixedRender.SetHidden(!display.Window.IsFullscreen())
	fixedRender.SetGetValueText(func() string {
		if display.Window.IsFixedRenderScale() {
			return l.String(l.ON)
		}
		return l.String(l.OFF)
	})
	fixedRender.SetTrigger(func() {
		isFixed := !display.Window.IsFixedRenderScale()
		display.Window.SetFixedRenderScale(isFixed)
		user.S().FixedRenderScale = isFixed
		cache.Clear()
		s.clearingCache = false
	})
	// Render size
	renderSize := ui.NewValueElement()
	renderSize.SetCenter(optionPos)
	renderSize.SetLabel(l.String(l.SETTINGS_GFX_RENDERSIZE))

	renderSizes := []display.RenderSize{
		display.RenderSizeTiny,
		display.RenderSizeSmall,
		display.RenderSizeMedium,
		display.RenderSizeLarge,
		display.RenderSizeMax,
	}
	currentRenderIdx := 0
	for i, size := range renderSizes {
		w, h := size.Value()
		if user.S().RenderWidth == w && user.S().RenderHeight == h {
			currentRenderIdx = i
			break
		}
	}
	renderSize.SetGetValueText(func() string {
		size := renderSizes[currentRenderIdx]
		return l.String(size.String())
	})
	renderSize.SetTrigger(func() {
		currentRenderIdx = (currentRenderIdx + 1) % len(renderSizes)
		size := renderSizes[currentRenderIdx]
		w, h := size.Value()

		user.S().RenderWidth = w
		user.S().RenderHeight = h
		display.Window.SetRenderSize(w, h)
		cache.Clear()
		s.clearingCache = false
	})
	group.Add(renderSize)
	optionPos.Y += optionsOffset

	// Fullscreen
	b := ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_GFX_FULLSCREEN))
	b.SetGetValueText(func() string {
		if user.S().Fullscreen {
			return l.String(l.ON)
		}
		return l.String(l.OFF)
	})
	b.SetTrigger(func() {
		user.S().Fullscreen = !user.S().Fullscreen
		display.Window.SetFullscreen(user.S().Fullscreen)
		fixedRender.SetHidden(!user.S().Fullscreen)

		cache.Clear()
		s.clearingCache = false
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	fixedRender.SetCenter(optionPos)
	group.Add(fixedRender)
	optionPos.Y += optionsOffset

	centerNoteColorB := ui.NewElement()
	cornerNoteColorB := ui.NewElement()
	noteColorThemes := types.AllNoteColorThemes()
	currentNoteColorThemeIdx := 0
	for i, theme := range noteColorThemes {
		if types.NoteColorTheme(user.S().NoteColorTheme) == theme {
			currentNoteColorThemeIdx = i
			break
		}
	}
	currentTheme := noteColorThemes[currentNoteColorThemeIdx]
	b = ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_GFX_NOTECOLOR))
	b.SetGetValueText(func() string {
		return l.String(string(noteColorThemes[currentNoteColorThemeIdx]))
	})
	b.SetTrigger(func() {
		currentNoteColorThemeIdx = (currentNoteColorThemeIdx + 1) % len(noteColorThemes)
		current := noteColorThemes[currentNoteColorThemeIdx]
		user.S().NoteColorTheme = string(current)

		if current == types.NoteColorThemeCustom {
			centerNoteColorB.SetHidden(false)
			cornerNoteColorB.SetHidden(false)
		} else {
			centerNoteColorB.SetHidden(true)
			cornerNoteColorB.SetHidden(true)
		}

		s.clearingCache = true
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	startingCustom := currentTheme == types.NoteColorThemeCustom
	colors := types.AllGameColors
	centerIdx := -1
	cornerIdx := -1
	currentCenter := types.GameColorFromHex(user.S().CenterNoteColor)
	currentCorner := types.GameColorFromHex(user.S().CornerNoteColor)
	for i, color := range colors {
		if color == currentCenter {
			centerIdx = i
		}
		if color == currentCorner {
			cornerIdx = i
		}
	}
	centerNoteColorB.SetCenter(optionPos)
	centerNoteColorB.SetText(l.String(l.CENTER))
	centerNoteColorB.InvertHoverColor()
	centerNoteColorB.SetTextColor(types.ColorFromHex(user.S().CenterNoteColor))
	centerNoteColorB.SetHidden(!startingCustom)
	centerNoteColorB.SetTrigger(func() {
		if centerIdx == -1 {
			centerIdx = 0
		}
		centerIdx = (centerIdx + 1) % len(colors)
		user.S().CenterNoteColor = types.HexFromColor(colors[centerIdx].C())
		centerNoteColorB.SetTextColor(types.ColorFromHex(user.S().CenterNoteColor))

		s.clearingCache = true
	})
	group.Add(centerNoteColorB)
	optionPos.Y += optionsOffset

	cornerNoteColorB.SetCenter(optionPos)
	cornerNoteColorB.SetText(l.String(l.CORNER))
	cornerNoteColorB.InvertHoverColor()
	cornerNoteColorB.SetTextColor(types.ColorFromHex(user.S().CornerNoteColor))
	cornerNoteColorB.SetHidden(!startingCustom)
	cornerNoteColorB.SetTrigger(func() {
		if cornerIdx == -1 {
			cornerIdx = 0
		}
		cornerIdx = (cornerIdx + 1) % len(colors)
		user.S().CornerNoteColor = types.HexFromColor(colors[cornerIdx].C())
		cornerNoteColorB.SetTextColor(types.ColorFromHex(user.S().CornerNoteColor))

		s.clearingCache = true

	})
	group.Add(cornerNoteColorB)
	optionPos.Y += optionsOffset

	// // Screen size
	// screenSizeButton.SetLabel(l.String(types.l.SETTINGS_GFX_SCREENSIZE))
	// screenSizeButton.SetCenter(optionPos)

	// resolutions := []struct{ w, h int }{
	// 	{1280, 720}, {1920, 1080}, {2560, 1440},
	// }
	// currentRes := 0
	// screenSizeButton.SetGetValueText(func() string {
	// 	return fmt.Sprintf("%dx%d", resolutions[currentRes].w, resolutions[currentRes].h)
	// })
	// screenSizeButton.SetTrigger(func() {
	// 	currentRes = (currentRes + 1) % len(resolutions)
	// 	res := resolutions[currentRes]

	// 	user.S().ScreenWidth = res.w
	// 	user.S().ScreenHeight = res.h
	// 	display.Window.SetDisplaySize(float64(res.w), float64(res.h))
	// })
	// group.Add(b)
	// optionPos.Y += optionsOffset
}

func (s *Settings) createAudioOptions(group *ui.UIGroup) {
	volumes := []struct {
		key    string
		value  *float64
		step   float64
		button *ui.ValueElement
	}{
		{l.SETTINGS_AUDIO_BGMVOLUME, &user.S().BGMVolume, 0.1, ui.NewValueElement()},
		{l.SETTINGS_AUDIO_SFXVOLUME, &user.S().SFXVolume, 0.1, ui.NewValueElement()},
		{l.SETTINGS_AUDIO_SONGVOLUME, &user.S().SongVolume, 0.1, ui.NewValueElement()},
	}

	optionPos := ui.Point{
		X: optionsStart.X,
		Y: optionsStart.Y,
	}
	for i := range volumes {
		v := &volumes[i] // Create pointer to avoid closure issue
		v.button.SetCenter(optionPos)
		v.button.SetLabel(l.String(v.key))
		v.button.SetGetValueText(func() string {
			return fmt.Sprintf("%.0f%%", *v.value*100)
		})
		v.button.SetTrigger(func() {
			*v.value += v.step
			if *v.value > 1.0 {
				*v.value = 0
			}
			if i == 0 {
				audio.SetBGMVolume(*v.value)
			} else if i == 1 {
				audio.SetSFXVolume(*v.value)
			} else if i == 2 {
				audio.SetSongVolume(*v.value)
			}
		})
		group.Add(v.button)
		optionPos.Y += optionsOffset
	}
}

func (s *Settings) createGameplayOptions(group *ui.UIGroup) {
	optionPos := ui.Point{
		X: optionsStart.X,
		Y: optionsStart.Y,
	}
	var b *ui.ValueElement

	keyConfig := ui.NewValueElement()
	keyConfig.SetCenter(optionPos)
	keyConfig.SetLabel(l.String(l.SETTINGS_GAME_KEY_CONFIG))
	keyConfig.SetGetValueText(func() string {
		return l.String(input.TrackKeyConfig(user.S().KeyConfig).String())
	})
	keyConfig.SetTrigger(func() {
		s.SetNextState(types.GameStateKeyConfig, &FloatStateArgs{
			Cb: func() {
				keyConfig.Refresh()
			},
		})
	})
	group.Add(keyConfig)
	optionPos.Y += optionsOffset

	// Note width
	noteWidths := []float32{1.0, 1.25, 1.5, 2.0, 0.5, 0.75}
	currentWidthIdx := 0
	for i, width := range noteWidths {
		if user.S().NoteWidth == width {
			currentWidthIdx = i
			break
		}
	}
	noteWidth := ui.NewValueElement()
	noteWidth.SetCenter(optionPos)
	noteWidth.SetLabel(l.String(l.SETTINGS_GAME_NOTEWIDTH))
	noteWidth.SetGetValueText(func() string {
		size := noteWidths[currentWidthIdx]
		return fmt.Sprintf("%.2fx", size)
	})
	noteWidth.SetTrigger(func() {
		currentWidthIdx = (currentWidthIdx + 1) % len(noteWidths)
		user.S().NoteWidth = noteWidths[currentWidthIdx]
		s.clearingCache = true
	})
	group.Add(noteWidth)
	optionPos.Y += optionsOffset

	// Lane Speed
	laneSpeeds := []float64{0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0, 5.5, 6.0, 6.5, 7.0, 7.5, 8.0, 8.5, 9.0, 9.5, 10.0}
	currentSpeedIdx := 0
	for i, speed := range laneSpeeds {
		if user.S().LaneSpeed == speed {
			currentSpeedIdx = i
			break
		}
	}
	b = ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_GAME_LANESPEED))
	b.SetGetValueText(func() string {
		size := laneSpeeds[currentSpeedIdx]
		return fmt.Sprintf("%.1fx", size)
	})
	b.SetTrigger(func() {
		currentSpeedIdx = (currentSpeedIdx + 1) % len(laneSpeeds)
		user.S().LaneSpeed = laneSpeeds[currentSpeedIdx]
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	// Edge play area
	b = ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_GAME_EDGEPLAYAREA))
	b.SetGetValueText(func() string {
		if user.S().EdgePlayArea {
			return l.String(l.ON)
		}
		return l.String(l.OFF)
	})
	b.SetTrigger(func() {
		user.S().EdgePlayArea = !user.S().EdgePlayArea
		s.clearingCache = true
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	// Offsets
	audioOffset := ui.NewValueElement()
	inputOffset := ui.NewValueElement()
	audioGetValue := func() string {
		return fmt.Sprintf("%dms", user.S().AudioOffset)
	}
	inputGetValue := func() string {
		return fmt.Sprintf("%dms", user.S().InputOffset)
	}

	audioOffset.SetCenter(optionPos)
	audioOffset.SetLabel(l.String(l.SETTINGS_GAME_AUDIOOFFSET))
	audioOffset.SetGetValueText(audioGetValue)
	audioOffset.SetTrigger(func() {
		s.SetNextState(types.GameStateOffset, &FloatStateArgs{
			Cb: func() {
				audioOffset.Refresh()
				inputOffset.Refresh()
			},
		})
	})
	optionPos.Y += optionsOffset

	inputOffset.SetCenter(optionPos)
	inputOffset.SetLabel(l.String(l.SETTINGS_GAME_INPUTOFFSET))
	inputOffset.SetGetValueText(inputGetValue)
	inputOffset.SetTrigger(func() {
		s.SetNextState(types.GameStateOffset, &FloatStateArgs{
			Cb: func() {
				audioOffset.Refresh()
				inputOffset.Refresh()
			},
		})
	})
	group.Add(audioOffset)
	group.Add(inputOffset)
}

func (s *Settings) createAccessOptions(group *ui.UIGroup) {
	optionPos := ui.Point{
		X: optionsStart.X,
		Y: optionsStart.Y,
	}
	var b *ui.ValueElement

	b = ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_ACCESS_NOHOLDNOTES))
	b.SetGetValueText(func() string {
		if user.S().DisableHoldNotes {
			return l.String(l.ON)
		}
		return l.String(l.OFF)
	})
	b.SetTrigger(func() {
		user.S().DisableHoldNotes = !user.S().DisableHoldNotes
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	b = ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_ACCESS_NOHITEFFECT))
	b.SetGetValueText(func() string {
		if user.S().DisableHitEffects {
			return l.String(l.ON)
		}
		return l.String(l.OFF)
	})
	b.SetTrigger(func() {
		user.S().DisableHitEffects = !user.S().DisableHitEffects
		s.clearingCache = true
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	b = ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_ACCESS_NOLANEEFFECT))
	b.SetGetValueText(func() string {
		if user.S().DisableLaneEffects {
			return l.String(l.ON)
		}
		return l.String(l.OFF)
	})
	b.SetTrigger(func() {
		user.S().DisableLaneEffects = !user.S().DisableLaneEffects
		s.clearingCache = true
	})
	group.Add(b)
	optionPos.Y += optionsOffset
}

func (s *Settings) Update() error {
	s.BaseGameState.Update()

	if input.K.Is(ebiten.KeyEscape, input.JustPressed) || input.K.Is(ebiten.KeyF1, input.JustPressed) {
		user.Save()

		if s.clearingCache {
			cache.Clear()
			s.clearingCache = false
		}
		s.SetNextState(types.GameStateBack, nil)
	}
	s.tabs.Update()
	return nil
}

func (s *Settings) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	ui.DrawNoteThemedRect(screen, &settingsCenter, &settingsSize)
	s.tabs.Draw(screen, opts)
}
