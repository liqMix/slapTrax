package user

type UserProfile struct {
	Name     string
	Settings UserSettings
}

var Current = UserProfile{
	Name:     "Zezima",
	Settings: DefaultSettings,
}

func Settings() UserSettings {
	return Current.Settings
}
