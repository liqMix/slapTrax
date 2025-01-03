package play

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

var (
	windowOffset = 0.05
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
	// lil space
	spacing = 0.1

	// Play Area
	playWidth  = maxX * 0.8
	playHeight = maxX - spacing - headerHeight
	playLeft   = centerX - (playWidth / 2)
	playTop    = headerBottom + spacing
	playRight  = centerX + (playWidth / 2)
	playBottom = playTop + playHeight

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

	// Markers
	markerTopLeft = ui.Point{
		X: playLeft,
		Y: playTop,
	}
	markerTopRight = ui.Point{
		X: playRight,
		Y: playTop,
	}
	markerBottomLeft = ui.Point{
		X: playLeft,
		Y: playBottom,
	}
	markerBottomRight = ui.Point{
		X: playRight,
		Y: playBottom,
	}

	// TrackNotes
	// trackSpacing = 0.0

	// trackHeight          = 1.0
	noteLength          = playWidth * 0.25
	noteWidth           = float32(3)
	cornerNoteCurve     = 0.7
	judgementLineLength = playWidth * 0.25
	judgementWidth      = noteWidth * 3
)

var points = []*ui.Point{
	{},
	{},
	{},
}

func GetJudgementPoints(track types.TrackName) []*ui.Point {
	return notePoints(track, judgementLineLength)
}

func GetNotePoints(track types.TrackName) []*ui.Point {
	return notePoints(track, noteLength)
}

func notePoints(track types.TrackName, length float64) []*ui.Point {
	centerLength := length / 2
	switch track {
	case types.LeftTop:
		points[0].X, points[0].Y = playLeft, playTop+length
		points[1].X, points[1].Y = playLeft, playTop
		points[2].X, points[2].Y = playLeft+length, playTop

	case types.LeftBottom:
		points[0].X, points[0].Y = playLeft, playBottom-length
		points[1].X, points[1].Y = playLeft, playBottom
		points[2].X, points[2].Y = playLeft+length, playBottom

	case types.CenterTop:
		points[0].X, points[0].Y = playCenterX-centerLength, playTop
		points[1].X, points[1].Y = playCenterX, playTop
		points[2].X, points[2].Y = playCenterX+centerLength, playTop

	case types.CenterBottom:
		points[0].X, points[0].Y = playCenterX-centerLength, playBottom
		points[1].X, points[1].Y = playCenterX, playBottom
		points[2].X, points[2].Y = playCenterX+centerLength, playBottom

	case types.RightTop:
		points[0].X, points[0].Y = playRight-length, playTop
		points[1].X, points[1].Y = playRight, playTop
		points[2].X, points[2].Y = playRight, playTop+length

	case types.RightBottom:
		points[0].X, points[0].Y = playRight-length, playBottom
		points[1].X, points[1].Y = playRight, playBottom
		points[2].X, points[2].Y = playRight, playBottom-length
	}

	return points
}

func GetMeasureMarkerPoints() []*ui.Point {
	return []*ui.Point{
		&markerTopLeft,
		&markerTopRight,
		&markerBottomRight,
		&markerBottomLeft,
	}
}
