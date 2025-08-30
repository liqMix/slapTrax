package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/system"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
)

type ChartSelection struct {
	types.BaseGameState
	group *ui.UIGroup
}

func NewChartSelectionState() *ChartSelection {
	audio.FadeOutBGM()

	s := &ChartSelection{}
	s.SetAction(input.ActionBack, func() {
		audio.PlaySFX(audio.SFXBack)
		s.SetNextState(types.GameStateTitle, nil)
	})

	// Create UI group
	group := ui.NewUIGroup()
	group.SetCenter(ui.Point{
		X: 0.5,
		Y: 0.5,
	})
	group.SetPaneled(true)

	buttonSize := ui.Point{
		X: 0.3,
		Y: 0.1,
	}
	textScale := 1.5
	offset := float64(ui.TextHeight(nil) * 2)

	// Title
	title := ui.NewElement()
	title.SetText("Chart Editor")
	title.SetTextScale(2.0)
	title.SetTextBold(true)
	title.SetCenter(ui.Point{
		X: 0.5,
		Y: 0.2,
	})
	title.SetDisabled(true)

	// Create New Chart button
	center := ui.Point{
		X: 0.5,
		Y: 0.4,
	}
	newChart := ui.NewElement()
	newChart.SetCenter(center)
	newChart.SetSize(buttonSize)
	newChart.SetText("Create New Chart")
	newChart.SetTextScale(textScale)
	newChart.SetTrigger(func() {
		// Open audio file dialog first
		audioPath, err := system.OpenAudioFileDialog()
		if err != nil {
			logger.Error("Failed to open audio file dialog: %v", err)
			return
		}
		if audioPath == "" {
			return // User cancelled
		}

		s.SetNextState(types.GameStateEditor, &EditorArgs{
			Song:      nil,
			ChartPath: "",
			AudioPath: audioPath,
		})
	})
	group.Add(newChart)
	center.Y += offset

	// Edit Existing Chart button
	editChart := ui.NewElement()
	editChart.SetCenter(center)
	editChart.SetSize(buttonSize)
	editChart.SetText("Edit Existing Chart")
	editChart.SetTextScale(textScale)
	editChart.SetTrigger(func() {
		// Open folder selection dialog
		folderPath, err := system.SelectFolderDialog("Select chart folder")
		if err != nil {
			logger.Error("Failed to open folder dialog: %v", err)
			return
		}
		if folderPath == "" {
			return // User cancelled
		}

		s.SetNextState(types.GameStateEditor, &EditorArgs{
			Song:      nil,
			ChartPath: folderPath,
			AudioPath: "",
		})
	})
	group.Add(editChart)
	center.Y += offset

	// Browse Charts button (future expansion)
	browseChart := ui.NewElement()
	browseChart.SetCenter(center)
	browseChart.SetSize(buttonSize)
	browseChart.SetText("Browse Charts")
	browseChart.SetTextScale(textScale)
	browseChart.SetDisabled(true) // Disabled for now
	browseChart.SetTrigger(func() {
		// TODO: Implement chart browser
	})
	group.Add(browseChart)

	group.SetSize(ui.Point{
		X: 0.4,
		Y: center.Y - 0.4 + buttonSize.Y,
	})

	s.group = group
	return s
}

func (s *ChartSelection) Update() error {
	s.BaseGameState.Update()
	s.group.Update()
	return nil
}

func (s *ChartSelection) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.group.Draw(screen, opts)
}