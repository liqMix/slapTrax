package user

import "github.com/liqmix/ebiten-holiday-2024/internal/external"

var (
	Current = external.M.GetUser
	Init    = external.M.Initialize
	S       = external.M.GetSettings
	Save    = external.M.SaveSettings
)
