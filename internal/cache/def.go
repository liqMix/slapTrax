package cache

func InitCaches() {
	Image = NewImageCache()
}

func Clear() {
	Image.Clear()
}
