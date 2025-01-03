package input

type Status int

const (
	None Status = iota
	JustPressed
	JustReleased
	Held
)
