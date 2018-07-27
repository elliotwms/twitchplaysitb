package drivers

type NoOpDriver struct {
}

func (*NoOpDriver) Click() func() {
	return func() {

	}
}

func (*NoOpDriver) Mouse(x, y int) func() {
	return func() {

	}
}

func (*NoOpDriver) PressKey(k string) func() {
	return func() {

	}
}

func (*NoOpDriver) ToggleKey(k string, down bool) func() {
	return func() {

	}
}
