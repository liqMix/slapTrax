package cache

// LayoutReinitRequested flag to trigger layout reinitialization from settings
var LayoutReinitRequested = false

func InitCaches() {
	Image = NewImageCache()
}

func Clear() {
	Image.Clear()
}

// RequestLayoutReinit requests that layouts be reinitialized on next draw
func RequestLayoutReinit() {
	LayoutReinitRequested = true
}
