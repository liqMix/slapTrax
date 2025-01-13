package play

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

func getNoteWidth() float32 {
	renderWidth, _ := display.Window.RenderSize()

	// scale note width based on render width
	// default defined for 1280
	return noteWidth * (float32(renderWidth) / 1280)
}

func getJudgementWidth() float32 {
	noteWidth := getNoteWidth()
	return noteWidth / 8
}

func SmoothProgress(progress float64) float32 {
	if progress >= 1 {
		return 1
	} else if progress <= 0 {
		return 0
	}
	return float32(minT / (minT + (1-minT)*(1-progress)))
}

func GetFadeAlpha(progress float32, max uint8) uint8 {
	if progress < fadeInThreshold {
		return 0
	}
	if progress > fadeOutThreshold {
		return max
	}
	fadeProgress := (progress - fadeInThreshold) / fadeRange
	return uint8(float32(max) * fadeProgress)
}

func GetNoteFadeAlpha(progress float32) uint8 {
	return GetFadeAlpha(progress, noteMaxAlpha)
}

func notePts(length float64) [][]*ui.Point {
	centerLength := length * centerNoteLengthRatio

	return [][]*ui.Point{
		// LeftBottom
		{
			&ui.Point{X: playLeft, Y: playBottom - length},
			&ui.Point{X: playLeft, Y: playBottom},
			&ui.Point{X: playLeft + length, Y: playBottom},
		},
		// LeftTop
		{
			&ui.Point{X: playLeft, Y: playTop + length},
			&ui.Point{X: playLeft, Y: playTop},
			&ui.Point{X: playLeft + length, Y: playTop},
		},
		// CenterBottom
		{
			&ui.Point{X: playCenterX - centerLength, Y: playBottom},
			&ui.Point{X: playCenterX, Y: playBottom},
			&ui.Point{X: playCenterX + centerLength, Y: playBottom},
		},
		// CenterTop
		{
			&ui.Point{X: playCenterX - centerLength, Y: playTop},
			&ui.Point{X: playCenterX, Y: playTop},
			&ui.Point{X: playCenterX + centerLength, Y: playTop},
		},
		// RightBottom
		{
			&ui.Point{X: playRight - length, Y: playBottom},
			&ui.Point{X: playRight, Y: playBottom},
			&ui.Point{X: playRight, Y: playBottom - length},
		},
		// RightTop
		{
			&ui.Point{X: playRight - length, Y: playTop},
			&ui.Point{X: playRight, Y: playTop},
			&ui.Point{X: playRight, Y: playTop + length},
		},
	}
}
