package drivers

import (
	"github.com/elliotwms/twitchplaysitb/commands"
	"github.com/go-vgo/robotgo"
)

type RobotGoDriver struct {
	gameWidth          int
	gameHeight         int
	xOffset            int
	yOffset            int
	coordinateResolver CoordinateResolver
}

// New creates a new instance of the RobotGoDriver
func NewRobotGoDriver(gameWidth, gameHeight, xOffset, yOffset int, resolver CoordinateResolver) *RobotGoDriver {
	return &RobotGoDriver{
		gameWidth:          gameWidth,
		gameHeight:         gameHeight,
		xOffset:            xOffset,
		yOffset:            yOffset,
		coordinateResolver: resolver,
	}
}

// click clicks the mouse
func (d *RobotGoDriver) Click() commands.Action {
	return func() {
		robotgo.Click()
	}
}

// mouse moves the mouse to a given set of pixel coordinates accounting for the x and y offset
func (d *RobotGoDriver) Mouse(x int, y int) commands.Action {
	// Prevent moving out of bounds on the x axis
	if x > d.gameWidth {
		x = d.gameWidth - 100 // with a safety buffer of 100px
	}

	// Do the same on the y axis
	if y > d.gameHeight {
		y = d.gameHeight - 100
	}

	return func() {
		robotgo.MoveMouseSmooth(x+d.xOffset, y+d.yOffset)
	}
}

// mouseGrid moves the mouse to a given set of map coordinates
func (d *RobotGoDriver) MouseGrid(a string, n string) commands.Action {
	x, y := d.coordinateResolver(a, n)

	return d.Mouse(x, y)
}

// pressKey taps a key
func (d *RobotGoDriver) PressKey(k string) commands.Action {
	return func() {
		robotgo.KeyTap(k)
	}
}

// toggleKey toggles a key up or down
func (d *RobotGoDriver) ToggleKey(k string, down bool) commands.Action {
	state := "up"
	if down {
		state = "down"
	}

	return func() {
		robotgo.KeyToggle(k, state)
	}
}
