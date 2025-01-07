package play

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

var (
	windowOffset = 0.025
	minX         = 0.0 + windowOffset
	minY         = 0.0 + windowOffset
	maxY         = 1.0 - (windowOffset * 2)
	maxX         = 1.0 - (windowOffset * 2)
	centerX      = minX + ((maxX - minX) / 2)
	centerY      = minY + ((maxY - minY) / 2)

	// Header
	headerWidth       = maxX * 0.8
	headerHeight      = maxY * 0.15
	headerLeft        = centerX - (headerWidth / 2)
	headerTop         = minY
	headerRight       = centerX + (headerWidth / 2)
	headerBottom      = headerTop + headerHeight
	headerCenterPoint = ui.Point{
		X: headerLeft + (headerWidth / 2),
		Y: headerTop + (headerHeight / 2),
	}

	// Play Area
	playWidth  = maxX * 0.75
	playHeight = 0.6
	playBottom = maxY - windowOffset
	playLeft   = centerX - (playWidth / 2)
	playTop    = playBottom - playHeight
	playRight  = centerX + (playWidth / 2)

	playCenterX     = playLeft + (playWidth / 2)
	playCenterY     = playTop + (playHeight / 2)
	playCenterPoint = ui.Point{
		X: playCenterX,
		Y: playCenterY,
	}
	centerComboSize = ui.Point{
		X: playWidth * 0.05,
		Y: playHeight * 0.05,
	}

	// Notes
	noteLength            = playWidth * 0.25
	noteWidth             = float32(40)
	noteComboRatio        = float32(1.5)
	noteMaxAlpha          = uint8(255)
	centerNoteLengthRatio = 0.5

	fadeInThreshold  = SmoothProgress(0.4)
	fadeOutThreshold = SmoothProgress(0.7)
	fadeRange        = fadeOutThreshold - fadeInThreshold
	minT             = 0.01 // Small value for vanishing point calculation

	// Judgement
	judgementWidth        = noteWidth / 8
	judgementPressedRatio = float32(3.0)

	// Vecs
	markerTopLeft = &ui.Point{
		X: playLeft,
		Y: playTop,
	}
	markerTopRight = &ui.Point{
		X: playRight,
		Y: playTop,
	}
	markerBottomLeft = &ui.Point{
		X: playLeft,
		Y: playBottom,
	}
	markerBottomRight = &ui.Point{
		X: playRight,
		Y: playBottom,
	}
	markerWidth         = float32(8)
	beatMarkerAlpha     = uint8(100)
	notePoints          = notePts(noteLength)
	measureMarkerPoints = []*ui.Point{
		markerTopLeft,
		markerTopRight,
		markerBottomRight,
		markerBottomLeft,
	}

	// hit brightness
	// hitBrightnessScale = float32(0.05)
	// comboBrightness    = float32(0.3)
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
func fullscreenLayout() {
	playWidth = 1.0
	playHeight = 1.0
	playBottom = 1.0
	playLeft = 0.0
	playTop = 0.0
	playRight = 1.0

	playCenterX = playLeft + (playWidth / 2)
	playCenterY = playTop + (playHeight / 2)
	playCenterPoint = ui.Point{
		X: playCenterX,
		Y: playCenterY,
	}

	judgementWidth = noteWidth
	judgementPressedRatio = float32(2.0)
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
