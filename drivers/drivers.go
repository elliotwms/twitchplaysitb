package drivers

import "github.com/elliotwms/twitchplaysitb/commands"

type Driver interface {
	Click() commands.Action
	Mouse(x, y int) commands.Action
	MouseGrid(a, n string) commands.Action
	PressKey(k string) commands.Action
	ToggleKey(k string, down bool) commands.Action
}

type CoordinateResolver func(a, n string) (x, y int)
