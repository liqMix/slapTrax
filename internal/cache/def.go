package cache

func InitCaches(renderWidth, renderHeight int) {
	Image = NewImageCache(renderWidth, renderHeight)
	Path = NewPathCache(renderWidth, renderHeight, 100)
}
