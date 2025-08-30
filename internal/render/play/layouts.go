package play

import "github.com/liqmix/slaptrax/internal/ui"

var (
	//// Overall
	windowOffset = 0.025
	minX         = 0.0 + windowOffset
	minY         = 0.0 + windowOffset
	maxY         = 1.0 - (windowOffset * 2)
	maxX         = 1.0 - (windowOffset * 2)
	centerX      = minX + ((maxX - minX) / 2)
	centerY      = minY + ((maxY - minY) / 2)

	//// Header
	headerWidth  = maxX * 0.3
	headerHeight = maxY * 0.065
	headerTop    = minY + windowOffset
	headerRight  = maxX
	headerLeft   = headerRight - headerWidth
	headerBottom = headerTop + headerHeight
	headerCenter = ui.Point{
		X: headerLeft + (headerWidth / 2),
		Y: headerTop + (headerHeight / 2),
	}

	//// Play Area
	playWidth   = maxX * 0.75
	playCenterX = centerX

	playHeight = 0.6
	playBottom = maxY - windowOffset
	playTop    = playBottom - playHeight

	playLeft  = centerX - (playWidth / 2)
	playRight = centerX + (playWidth / 2)

	playCenterY     = playTop + (playHeight / 2)
	playCenterPoint = ui.Point{
		X: playCenterX,
		Y: playCenterY,
	}

	// Notes
	noteLength            = playWidth * 0.25
	noteWidth             = float32(30)
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
	comboCenter = ui.Point{
		X: playCenterX,
		Y: playCenterY - 0.1,
	}
	// hit brightness
	// hitBrightnessScale = float32(0.05)
	// comboBrightness    = float32(0.3)
)

func applyLayout() {
	//// Overall
	windowOffset = 0.025
	minX = 0.0 + windowOffset
	minY = 0.0 + windowOffset
	maxY = 1.0 - (windowOffset * 2)
	maxX = 1.0 - (windowOffset * 2)
	centerX = minX + ((maxX - minX) / 2)
	centerY = minY + ((maxY - minY) / 2)

	//// Header
	headerWidth = maxX * 0.3
	headerHeight = maxY * 0.065
	headerTop = minY + windowOffset
	headerRight = maxX
	headerLeft = headerRight - headerWidth
	headerBottom = headerTop + headerHeight
	headerCenter = ui.Point{
		X: headerLeft + (headerWidth / 2),
		Y: headerTop + (headerHeight / 2),
	}

	// Notes
	noteLength = playWidth * 0.25
	noteWidth = float32(40)
	noteComboRatio = float32(1.5)
	noteMaxAlpha = uint8(255)
	centerNoteLengthRatio = 0.5

	fadeInThreshold = SmoothProgress(0.4)
	fadeOutThreshold = SmoothProgress(0.7)
	fadeRange = fadeOutThreshold - fadeInThreshold
	minT = 0.01 // Small value for vanishing point calculation

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
	markerWidth = float32(8)
	beatMarkerAlpha = uint8(100)
	notePoints = notePts(noteLength)
	measureMarkerPoints = []*ui.Point{
		markerTopLeft,
		markerTopRight,
		markerBottomRight,
		markerBottomLeft,
	}
	comboCenter = ui.Point{
		X: playCenterX,
		Y: playCenterY + playHeight/4,
	}
}

func applyDefaultLayout() {
	playWidth = maxX * 0.75
	playCenterX = centerX

	playHeight = 0.6
	playBottom = maxY - windowOffset
	playTop = playBottom - playHeight

	playLeft = centerX - (playWidth / 2)
	playRight = centerX + (playWidth / 2)

	playCenterY = playTop + (playHeight / 2)
	playCenterPoint = ui.Point{
		X: playCenterX,
		Y: playCenterY,
	}

	// Judgement
	judgementWidth = noteWidth / 8
	judgementPressedRatio = float32(3.0)
	applyLayout()
}

func applyEdgeLayout() {
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

	judgementWidth = noteWidth / 4
	// judgementPressedRatio = float32(2.0)
	applyLayout()
}
