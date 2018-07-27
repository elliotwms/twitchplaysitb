package drivers

type Driver interface {
	Click() func()
	Mouse(x, y int) func()
	PressKey(k string) func()
	ToggleKey(k string, down bool) func()
}
