package user

type User struct {
	Name     string
	Settings *UserSettings
}

var u = &User{
	Name:     "funyarinpa",
	Settings: NewUserSettings(),
}

func Init() {
	if u == nil {
		u = &User{
			Name:     "funyarinpa",
			Settings: NewUserSettings(),
		}
		u.Settings.Apply()
	}

}

func S() *UserSettings {
	return u.Settings
}
