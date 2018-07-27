package main

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

type Command struct {
	Text        string // The raw command input
	Description string // The command description
	Actions     []Action
}

// GetHash hashes the command description (which should be unique depending on the command arguments)
func (c *Command) GetHash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(c.Description)))
}

type Action func()

func Parse(t string) *Command {
	// lower case all commands to normalize them
	t = strings.ToLower(t)

	c := &Command{
		Text: t,
	}

	// Simple commands

	// Click
	if match, _ := regexp.MatchString("^click$", t); match {
		c.Description = "Click the mouse"
		c.Actions = []Action{
			click(),
		}
	}

	// End turn
	if match, _ := regexp.MatchString("^endturn$", t); match {
		c.Description = "End turn"
		c.Actions = []Action{
			func() {
				robotgo.KeyTap("space")
			},
		}
	}

	// Undo move
	if match, _ := regexp.MatchString("^undo$", t); match {
		c.Description = "Undo move"
		c.Actions = []Action{
			func() {
				robotgo.KeyTap("shift")
			},
		}
	}

	// Reset turn
	if match, _ := regexp.MatchString("^reset$", t); match {
		c.Description = "Reset turn"
		c.Actions = []Action{
			func() {
				robotgo.KeyTap("backspace")
			},
		}
	}

	// Deselect weapon
	if match, _ := regexp.MatchString("^deselect|disarm$", t); match {
		c.Description = "Deselect weapon"
		c.Actions = []Action{
			func() {
				robotgo.KeyTap("q")
			},
		}
	}

	// Next unit
	if match, _ := regexp.MatchString("^next$", t); match {
		c.Description = "Select next unit"
		c.Actions = []Action{
			func() {
				robotgo.KeyTap("tab")
			},
		}
	}

	// Next unit (batch)
	if r := regexp.MustCompile("^next ([2-9])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		count, _ := strconv.Atoi(ss[1])

		c.Description = fmt.Sprintf("Select next unit #%d", count)
		c.Actions = []Action{
			func() {
				for i := 0; i < count; i++ {
					robotgo.KeyTap("tab")
					time.Sleep(1 * time.Second)
				}
			},
		}
	}

	// Less simple commands

	// Mouse
	// Moves the mouse to the given coordinates
	if r := regexp.MustCompile("^mouse ([0-9]*) ([0-9]*)$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x, _ := strconv.Atoi(ss[1])
		y, _ := strconv.Atoi(ss[2])

		c.Description = fmt.Sprintf("Move the mouse to x: %d, y: %d", x, y)
		c.Actions = []Action{
			mouse(x, y),
		}
	}

	// Click
	// Moves to coordinates and clicks
	if r := regexp.MustCompile("^click ([0-9]*) ([0-9]*)$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x, _ := strconv.Atoi(ss[1])
		y, _ := strconv.Atoi(ss[2])

		c.Description = fmt.Sprintf("Click the mouse at x: %d, y: %d", x, y)
		c.Actions = []Action{
			mouse(x, y),
			click(),
		}
	}

	// Mouse grid
	// Moves the mouse to the grid reference
	if r := regexp.MustCompile("^mouse ([a-h])([1-8]*)$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x := strings.ToUpper(ss[1])
		y := ss[2]

		c.Description = fmt.Sprintf("Move the mouse to tile %s%s", x, y)
		c.Actions = []Action{
			mouseGrid(x, y),
		}
	}

	// Click grid
	// Moves the mouse to the grid reference and clicks
	if r := regexp.MustCompile("^click ([a-h])([1-8]*)$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x := strings.ToUpper(ss[1])
		y := ss[2]

		c.Description = fmt.Sprintf("Click the mouse at tile %s%s", x, y)
		c.Actions = []Action{
			mouseGrid(strings.ToUpper(ss[1]), ss[2]),
			click(),
		}
	}

	// Select unit
	// Select a specific unit by hotkey
	if r := regexp.MustCompile("^select (mech|deployed|mission) ([1-3])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		c.Description = fmt.Sprintf("Select %s unit #%s", ss[1], ss[2])
		c.Actions = []Action{
			selectUnit(ss[1], ss[2]),
		}
	}

	// Move unit
	// Move a specific unit to a given map coordinate
	if r := regexp.MustCompile("^move (mech|deployed|mission) ([1-3]) ([a-h])([1-8])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x := strings.ToUpper(ss[3])
		y := ss[4]

		c.Description = fmt.Sprintf("Move %s unit #%s to %s%s", ss[1], ss[2], x, y)
		c.Actions = []Action{
			selectUnit(ss[1], ss[2]),
			mouseGrid(x, y),
			click(),
		}
	}

	// Select weapon
	// Select a specific weapon by hotkey
	if r := regexp.MustCompile("^weapon ([1-2])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		c.Description = fmt.Sprintf("Arm weapon #%s", ss[1])
		c.Actions = []Action{
			selectWeapon(ss[1]),
		}
	}

	// Attack
	// Attack with a unit using a weapon at a given tile
	if r := regexp.MustCompile("^attack (mech|deployed|mission) ([1-3]) ([1-2]) ([a-h])([1-8])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x := strings.ToUpper(ss[4])
		y := ss[5]

		c.Description = fmt.Sprintf("Attacking with %s unit #%s using weapon %s on tile %s%s", ss[1], ss[2], ss[3], x, y)
		c.Actions = []Action{
			selectUnit(ss[1], ss[2]),
			selectWeapon(ss[3]),
			mouseGrid(x, y),
			click(),
		}
	}

	// Repair
	// Repair a unit at a given tile
	if r := regexp.MustCompile("^repair ([1-3]) ([a-h])([1-8])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x := strings.ToUpper(ss[2])
		y := ss[3]

		c.Description = fmt.Sprintf("Repairing mech #%d at tile %d%d", ss[1], x, y)
		c.Actions = []Action{
			selectUnit("mech", ss[1]),
			repair(),
			mouseGrid(x, y),
			click(),
		}
	}

	// On-off commands

	// Toggle info
	// Toggle the info tooltip
	if r := regexp.MustCompile("^info (on|off)$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		state := "up"
		if ss[1] == "on" {
			state = "down"
		}

		c.Description = fmt.Sprintf("Turning info tooltip %s", ss[1])
		c.Actions = []Action{
			func() {
				robotgo.KeyToggle("control", state)
			},
		}
	}

	// Toggle turn order
	// Toggle the turn order tooltips
	if r := regexp.MustCompile("^order (on|off)$"); r.MatchString(t) {

		ss := r.FindStringSubmatch(t)

		state := "up"
		if ss[1] == "on" {
			state = "down"
		}

		c.Description = fmt.Sprintf("Turning turn order tooltips %s", ss[1])
		c.Actions = []Action{
			func() {
				robotgo.KeyToggle("alt", state)
			},
		}
	}

	// Calibration commands

	// Calibrate
	// Moves the mouse from the top left to bottom right corner and around the grid. Used for offset calibration
	if match, _ := regexp.MatchString("^calibrate$", t); match {
		c.Description = "Calibrating"
		c.Actions = []Action{
			mouse(0, 0),
			mouse(1280, 720),
			mouseGrid("A", "1"),
			mouseGrid("H", "1"),
			mouseGrid("H", "8"),
			mouseGrid("A", "8"),
		}
	}

	if len(c.Actions) == 0 {
		return nil
	}

	return c
}

// Commands
// These commands wrap the action into a callable zero-argument function

// click clicks the mouse
func click() func() {
	return func() {
		robotgo.Click()
	}
}

// mouse moves the mouse to a given set of pixel coordinates accounting for the x and y offset
func mouse(x int, y int) func() {
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
func mouseGrid(a string, n string) func() {
	x, y := GetCoordinates(a, n)

	return mouse(x, y)
}

// selectUnit selects a unit by type and number
func selectUnit(t string, n string) func() {
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
func selectWeapon(n string) func() {
	return func() {
		robotgo.KeyTap(n)
	}
}

// repair hits the repair shortcut
func repair() func() {
	return func() {
		robotgo.KeyTap("r") // Repair
	}
}
