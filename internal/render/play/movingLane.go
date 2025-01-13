package play

import "github.com/liqmix/slaptrax/internal/ui"

func OscillateWindowOffset(currentTime int64) {
	ChangeWindowOffset(0.05 + (0.01 * float64(currentTime%1000) / 1000))
}

func ChangeWindowOffset(offset float64) {
	windowOffset = windowOffset + offset/100
	minX = 0.0 - windowOffset
	minY = 0.0 - windowOffset
	maxY = 1.0 + windowOffset
	maxX = 1.0 + windowOffset
	centerX = minX + ((maxX - minX) / 2)
	centerY = minY + ((maxY - minY) / 2)

	// Header
	headerWidth = 0
	headerHeight = 0

	// Play Area
	playWidth = maxX
	playHeight = maxY
	playBottom = maxY
	playLeft = minX
	playTop = minY
	playRight = maxX

	playCenterX = playLeft + (playWidth / 2)
	playCenterY = playTop + (playHeight / 2)
	playCenterPoint = ui.Point{
		X: playCenterX,
		Y: playCenterY,
	}
}
