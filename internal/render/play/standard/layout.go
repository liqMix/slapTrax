package standard

import (
	"image/color"

	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

type LaneConfig struct {
	Left           ui.NormPoint // Bottom left point (0-1 space)
	Center         ui.NormPoint // Bottom center point (0-1 space)
	Right          ui.NormPoint // Bottom right point (0-1 space)
	VanishingPoint ui.NormPoint // Convergence point (0-1 space)
	CurveAmount    float64      // 0 = right angle at center, 1 = straight line
	NoteWidth      float32      // Width of notes as a ratio of track width
	NoteColor      color.RGBA
}

var (
	offset      = 0.05
	spacing     = 0.0
	cornerCurve = 0.7

	// Main
	mainHeight       = 0.60
	mainWidth        = 0.0
	centerTrackWidth = mainWidth / 2
	mainLeft         = offset
	mainRight        = offset + mainWidth
	mainTop          = 1 - (offset + mainHeight)
	mainBottom       = 1 - offset
	mainCenter       = ui.NormPoint{
		X: mainRight - (mainWidth / 2),
		Y: mainBottom - (mainHeight / 2),
	}
	mainTopSpacing    = spacing / 2
	mainBottomSpacing = (centerTrackWidth / 2) + spacing

	mainScoreScale = 0.035

	// Edge
	edgeLeft   = mainRight + offset
	edgeRight  = 1 - offset
	edgeWidth  = edgeRight - edgeLeft
	edgeHeight = 0.8 * mainHeight
	edgeTop    = 1 - (offset + edgeHeight)
	edgeBottom = 1 - offset
	edgeCenter = ui.NormPoint{
		X: edgeLeft + (edgeWidth / 2),
		Y: mainCenter.Y,
	}
	edgeTapWidth = edgeWidth / 3

	// Colors
	black     = color.RGBA{0, 0, 0, 255}
	gray      = color.RGBA{100, 100, 100, 255}
	white     = color.RGBA{200, 200, 200, 255}
	orange    = color.RGBA{255, 165, 0, 255}
	blue      = color.RGBA{0, 0, 255, 255}
	lightBlue = color.RGBA{173, 216, 230, 255}
	yellow    = color.RGBA{255, 255, 0, 255}

	// Note Width (as ratio to track width)
	targetThick = 4
)

var laneConfigs = map[song.TrackName]*LaneConfig{}

// Modifies layout to account for edge tracks
func SetLayout(displayEdgeTracks bool) {
	var width float64
	if displayEdgeTracks {
		width = 0.6
	} else {
		width = 1 - (offset * 2)
	}

	if width == mainWidth {
		return
	}

	mainWidth = width

	// Set all parameters that depend on mainWidth
	centerTrackWidth = mainWidth / 2
	mainRight = offset + mainWidth
	mainCenter = ui.NormPoint{
		X: mainRight - (mainWidth / 2),
		Y: mainBottom - (mainHeight / 2),
	}
	mainBottomSpacing = (centerTrackWidth / 2) + spacing

	// Edge
	edgeLeft = mainRight + offset
	edgeWidth = edgeRight - edgeLeft
	edgeCenter = ui.NormPoint{
		X: edgeLeft + (edgeWidth / 2),
		Y: mainCenter.Y,
	}
	edgeTapWidth = edgeWidth / 3

	// Update lane configs
	laneConfigs = map[song.TrackName]*LaneConfig{
		song.LeftBottom: {
			CurveAmount:    cornerCurve,
			VanishingPoint: mainCenter,
			Left: ui.NormPoint{
				X: mainLeft,
				Y: mainCenter.Y + spacing,
			},
			Center: ui.NormPoint{
				X: mainLeft,
				Y: mainBottom,
			},
			Right: ui.NormPoint{
				X: mainCenter.X - mainBottomSpacing,
				Y: mainBottom,
			},
		},
		song.LeftTop: {
			CurveAmount:    cornerCurve,
			VanishingPoint: mainCenter,
			Left: ui.NormPoint{
				X: mainCenter.X - mainTopSpacing,
				Y: mainTop,
			},
			Center: ui.NormPoint{
				X: mainLeft,
				Y: mainTop,
			},
			Right: ui.NormPoint{
				X: mainLeft,
				Y: mainCenter.Y - spacing,
			},
		},

		song.RightBottom: {
			CurveAmount:    cornerCurve,
			VanishingPoint: mainCenter,
			Left: ui.NormPoint{
				X: mainCenter.X + mainBottomSpacing,
				Y: mainBottom,
			},
			Center: ui.NormPoint{
				X: mainRight,
				Y: mainBottom,
			},
			Right: ui.NormPoint{
				X: mainRight,
				Y: mainCenter.Y + spacing,
			},
		},
		song.RightTop: {
			CurveAmount:    cornerCurve,
			VanishingPoint: mainCenter,
			Left: ui.NormPoint{
				X: mainRight,
				Y: mainCenter.Y - spacing,
			},
			Center: ui.NormPoint{
				X: mainRight,
				Y: mainTop,
			},
			Right: ui.NormPoint{
				X: mainCenter.X + mainTopSpacing,
				Y: mainTop,
			},
		},

		song.Center: {
			CurveAmount:    0,
			VanishingPoint: mainCenter,
			Left: ui.NormPoint{
				X: mainCenter.X - centerTrackWidth/2,
				Y: mainBottom,
			},
			Center: ui.NormPoint{
				X: mainCenter.X,
				Y: mainBottom,
			},
			Right: ui.NormPoint{
				X: mainCenter.X + centerTrackWidth/2,
				Y: mainBottom,
			},
		},

		song.EdgeTop: {
			CurveAmount:    0,
			VanishingPoint: edgeCenter,
			Left: ui.NormPoint{
				X: edgeLeft,
				Y: edgeTop,
			},
			Center: ui.NormPoint{
				X: edgeCenter.X,
				Y: edgeTop,
			},
			Right: ui.NormPoint{
				X: edgeRight,
				Y: edgeTop,
			},
		},
		song.EdgeTap1: {
			CurveAmount:    0,
			VanishingPoint: edgeCenter,
			Left: ui.NormPoint{
				X: edgeLeft,
				Y: edgeBottom,
			},
			Center: ui.NormPoint{
				X: (edgeLeft + edgeTapWidth/2),
				Y: edgeBottom,
			},
			Right: ui.NormPoint{
				X: edgeLeft + edgeTapWidth,
				Y: edgeBottom,
			},
		},

		song.EdgeTap2: {
			CurveAmount:    0,
			VanishingPoint: edgeCenter,
			Left: ui.NormPoint{
				X: edgeLeft + edgeTapWidth,
				Y: edgeBottom,
			},
			Center: ui.NormPoint{
				X: (edgeLeft + edgeTapWidth + edgeTapWidth/2),
				Y: edgeBottom,
			},
			Right: ui.NormPoint{
				X: edgeLeft + edgeTapWidth*2,
				Y: edgeBottom,
			},
		},

		song.EdgeTap3: {
			CurveAmount:    0,
			VanishingPoint: edgeCenter,
			Left: ui.NormPoint{
				X: edgeLeft + edgeTapWidth*2,
				Y: edgeBottom,
			},
			Center: ui.NormPoint{
				X: edgeLeft + edgeTapWidth*2 + edgeTapWidth/2,
				Y: edgeBottom,
			},
			Right: ui.NormPoint{
				X: edgeLeft + edgeTapWidth*3,
				Y: edgeBottom,
			},
		},
	}
}
