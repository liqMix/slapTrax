package cache

func InitCaches() {
	Image = NewImageCache()
	Path = NewPathCache()
}

func Clear() {
	Image.Clear()
	Path.Clear()
}
