package user

type UserProfile struct {
	Username string
	Settings UserSettings
}

var Current = UserProfile{
	Username: "Zezima",
	Settings: DefaultSettings,
}
