package user

import "github.com/liqmix/slaptrax/internal/external"

var (
	Current = external.M.GetUser
	Init    = external.M.Initialize
	S       = external.M.GetSettings
	Save    = external.M.SaveSettings
)
