package state

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

var (
	settingsCenter = ui.Point{X: 0.5, Y: 0.5}
	settingsSize   = ui.Point{X: 0.5, Y: 0.8}
	tabStart       = ui.Point{X: 0.25, Y: 0.15}
	tabOffset      = 0.1

	optionsStart  = ui.Point{X: 0.5, Y: 0.25}
	optionsOffset = 0.1
)

type Settings struct {
	types.BaseGameState

	tabs *ui.Tabs
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
	s.tabs.SetCenter(ui.Point{X: 0.5, Y: 0.15}, ui.Point{X: 0.15, Y: 0.2}, 0.05)
}

func (s *Settings) createGraphicsOptions(group *ui.UIGroup) {
	optionPos := ui.Point{
		X: optionsStart.X,
		Y: optionsStart.Y,
	}

	// screenSizeButton := ui.NewValueElement()
	// if user.S.Fullscreen {
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
		user.S.FixedRenderScale = isFixed
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
		if user.S.RenderWidth == w && user.S.RenderHeight == h {
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

		user.S.RenderWidth = w
		user.S.RenderHeight = h
		display.Window.SetRenderSize(w, h)
	})
	group.Add(renderSize)
	optionPos.Y += optionsOffset

	// Fullscreen
	b := ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_GFX_FULLSCREEN))
	b.SetGetValueText(func() string {
		if ebiten.IsFullscreen() {
			return l.String(l.ON)
		}
		return l.String(l.OFF)
	})
	b.SetTrigger(func() {
		display.Window.SetFullscreen(!ebiten.IsFullscreen())
		user.S.Fullscreen = ebiten.IsFullscreen()

		fixedRender.SetHidden(!ebiten.IsFullscreen())
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
		if types.NoteColorTheme(user.S.NoteColorTheme) == theme {
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
		user.S.NoteColorTheme = string(current)
		display.RebuildCaches()
		if current == types.NoteColorThemeCustom {
			centerNoteColorB.SetHidden(false)
			cornerNoteColorB.SetHidden(false)
		} else {
			centerNoteColorB.SetHidden(true)
			cornerNoteColorB.SetHidden(true)
		}
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	startingCustom := currentTheme == types.NoteColorThemeCustom
	colors := types.AllGameColors
	centerIdx := -1
	cornerIdx := -1
	currentCenter := types.GameColorFromHex(user.S.CenterNoteColor)
	currentCorner := types.GameColorFromHex(user.S.CornerNoteColor)
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
	centerNoteColorB.ForceTextColor()
	centerNoteColorB.SetTextColor(types.ColorFromHex(user.S.CenterNoteColor))
	centerNoteColorB.SetHidden(!startingCustom)
	centerNoteColorB.SetTrigger(func() {
		if centerIdx == -1 {
			centerIdx = 0
		}
		centerIdx = (centerIdx + 1) % len(colors)
		user.S.CenterNoteColor = types.HexFromColor(colors[centerIdx].C())
		centerNoteColorB.SetTextColor(types.ColorFromHex(user.S.CenterNoteColor))
		display.RebuildCaches()
	})
	group.Add(centerNoteColorB)
	optionPos.Y += optionsOffset

	cornerNoteColorB.SetCenter(optionPos)
	cornerNoteColorB.SetText(l.String(l.CORNER))
	cornerNoteColorB.ForceTextColor()
	cornerNoteColorB.SetTextColor(types.ColorFromHex(user.S.CornerNoteColor))
	cornerNoteColorB.SetHidden(!startingCustom)
	cornerNoteColorB.SetTrigger(func() {
		if cornerIdx == -1 {
			cornerIdx = 0
		}
		cornerIdx = (cornerIdx + 1) % len(colors)
		user.S.CornerNoteColor = types.HexFromColor(colors[cornerIdx].C())
		cornerNoteColorB.SetTextColor(types.ColorFromHex(user.S.CornerNoteColor))
		display.RebuildCaches()
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

	// 	user.S.ScreenWidth = res.w
	// 	user.S.ScreenHeight = res.h
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
		{l.SETTINGS_AUDIO_BGMVOLUME, &user.S.BGMVolume, 0.1, ui.NewValueElement()},
		{l.SETTINGS_AUDIO_SFXVOLUME, &user.S.SFXVolume, 0.1, ui.NewValueElement()},
		{l.SETTINGS_AUDIO_SONGVOLUME, &user.S.SongVolume, 0.1, ui.NewValueElement()},
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
			audio.SetVolume(&audio.Volume{
				BGM:  user.S.BGMVolume,
				SFX:  user.S.SFXVolume,
				Song: user.S.SongVolume,
			})
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

	// Lane Speed
	laneSpeeds := []float64{0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0}
	currentSpeedIdx := 0
	for i, speed := range laneSpeeds {
		if user.S.LaneSpeed == speed {
			currentSpeedIdx = i
			break
		}
	}
	b := ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_GAME_LANESPEED))
	b.SetGetValueText(func() string {
		size := laneSpeeds[currentSpeedIdx]
		return fmt.Sprintf("%.1fx", size)
	})
	b.SetTrigger(func() {
		currentSpeedIdx = (currentSpeedIdx + 1) % len(laneSpeeds)
		user.S.LaneSpeed = laneSpeeds[currentSpeedIdx]
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	// Offsets
	b = ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_GAME_AUDIOOFFSET))
	b.SetGetValueText(func() string {
		return fmt.Sprintf("%dms", user.S.AudioOffset)
	})
	b.SetTrigger(func() {
		s.SetNextState(types.GameStateOffset, &FloatStateArgs{
			onClose: func() {
				b.SetLabel(fmt.Sprintf("%s: %dms", l.String(l.SETTINGS_GAME_AUDIOOFFSET), user.S.AudioOffset))
			},
		})
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	b = ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_GAME_INPUTOFFSET))
	b.SetGetValueText(func() string {
		return fmt.Sprintf("%dms", user.S.InputOffset)
	})
	b.SetTrigger(func() {
		s.SetNextState(types.GameStateOffset, &FloatStateArgs{
			onClose: func() {
				b.SetLabel(fmt.Sprintf("%s: %dms", l.String(l.SETTINGS_GAME_INPUTOFFSET), user.S.InputOffset))
			},
		})
	})
	group.Add(b)
}

func (s *Settings) createAccessOptions(group *ui.UIGroup) {
	optionPos := ui.Point{
		X: optionsStart.X,
		Y: optionsStart.Y,
	}
	b := ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_ACCESS_NOHOLDNOTES))
	b.SetGetValueText(func() string {
		if user.S.DisableHoldNotes {
			return l.String(l.ON)
		}
		return l.String(l.OFF)
	})
	b.SetTrigger(func() {
		user.S.DisableHoldNotes = !user.S.DisableHoldNotes
		user.Save()
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	b = ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_ACCESS_NOHITEFFECT))
	b.SetGetValueText(func() string {
		if user.S.DisableHitEffects {
			return l.String(l.ON)
		}
		return l.String(l.OFF)
	})
	b.SetTrigger(func() {
		user.S.DisableHitEffects = !user.S.DisableHitEffects
		user.Save()
	})
	group.Add(b)
	optionPos.Y += optionsOffset

	b = ui.NewValueElement()
	b.SetCenter(optionPos)
	b.SetLabel(l.String(l.SETTINGS_ACCESS_NOLANEEFFECT))
	b.SetGetValueText(func() string {
		if user.S.DisableLaneEffects {
			return l.String(l.ON)
		}
		return l.String(l.OFF)
	})
	b.SetTrigger(func() {
		user.S.DisableLaneEffects = !user.S.DisableLaneEffects
		user.Save()
	})
	group.Add(b)
}

func (s *Settings) Update() error {
	if input.K.Is(ebiten.KeyEscape, input.JustPressed) {
		user.Save()
		s.SetNextState(types.GameStateBack, nil)
	}
	s.tabs.Update()
	return nil
}

func (s *Settings) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	ui.DrawNoteThemedRect(screen, &settingsCenter, &settingsSize)
	s.tabs.Draw(screen, opts)
}
