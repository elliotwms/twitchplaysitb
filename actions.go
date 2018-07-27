package main

import "github.com/go-vgo/robotgo"

// Actions
// These commands wrap the action into a callable zero-argument function

type Action func()

// click clicks the mouse
func click() Action {
	return func() {
		robotgo.Click()
	}
}

// mouse moves the mouse to a given set of pixel coordinates accounting for the x and y offset
func mouse(x int, y int) Action {
	// Prevent moving out of bounds on the x axis
	if x > gameWidth {
		x = gameWidth - 100 // with a safety buffer of 100px
	}

	// Do the same on the y axis
	if y > gameHeight {
		y = gameHeight - 100
	}

	return func() {
		robotgo.MoveMouseSmooth(x+xOffset, y+yOffset)
	}
}

// mouseGrid moves the mouse to a given set of map coordinates
func mouseGrid(a string, n string) Action {
	x, y := GetCoordinates(a, n)

	return mouse(x, y)
}

// selectUnit selects a unit by type and number
func selectUnit(t string, n string) Action {
	k := "a"

	switch t {
	case "mech":
		switch n {
		case "1":
			k = "a"
		case "2":
			k = "s"
		case "3":
			k = "d"
		}
	case "deployed":
		switch n {
		case "1":
			k = "f"
		case "2":
			k = "g"
		case "3":
			k = "h"
		}
	case "mission":
		switch n {
		case "1":
			k = "z"
		case "2":
			k = "x"
		}
	}

	return func() {
		robotgo.KeyTap(k)
	}
}

// selectWeapon selects a weapon by number
func selectWeapon(n string) Action {
	return func() {
		robotgo.KeyTap(n)
	}
}

// repair hits the repair shortcut
func repair() Action {
	return func() {
		robotgo.KeyTap("r") // Repair
	}
}

// pressKey taps a key
func pressKey(k string) Action {
	return func() {
		robotgo.KeyTap(k)
	}
}

// toggleKey toggles a key up or down
func toggleKey(k string, down bool) Action {
	state := "up"
	if down {
		state = "down"
	}

	return func() {
		robotgo.KeyToggle(k, state)
	}
}
