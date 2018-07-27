package drivers

import "github.com/elliotwms/twitchplaysitb/commands"

type Driver interface {
	Click() commands.Action
	Mouse(x, y int) commands.Action
	PressKey(k string) commands.Action
	ToggleKey(k string, down bool) commands.Action
}
