package input

var m = newMouse()
var k = newKeyboard()
var (
	M      = m
	K      = k
	Update = func() {
		M.update()
		K.update()
	}
)
