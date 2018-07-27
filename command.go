package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elliotwms/twitchplaysitb/commands"
	"github.com/elliotwms/twitchplaysitb/drivers"
)

func Parse(t string) *commands.Command {
	driver := drivers.NewRobotGoDriver(gameWidth, gameHeight, xOffset, yOffset, GetCoordinates)

	// lower case all commands to normalize them
	t = strings.ToLower(t)

	c := &commands.Command{
		Text: t,
	}

	// Simple commands

	// Click
	if match, _ := regexp.MatchString("^click$", t); match {
		c.Description = "Click the mouse"
		c.Actions = []commands.Action{
			driver.Click(),
		}
	}

	// End turn
	if match, _ := regexp.MatchString("^endturn$", t); match {
		c.Description = "End turn"
		c.Actions = []commands.Action{
			driver.PressKey("space"),
		}
	}

	// Undo move
	if match, _ := regexp.MatchString("^undo$", t); match {
		c.Description = "Undo move"
		c.Actions = []commands.Action{
			driver.PressKey("shift"),
		}
	}

	// Reset turn
	if match, _ := regexp.MatchString("^reset$", t); match {
		c.Description = "Reset turn"
		c.Actions = []commands.Action{
			driver.PressKey("backspace"),
		}
	}

	// Deselect weapon
	if match, _ := regexp.MatchString("^deselect|disarm$", t); match {
		c.Description = "Deselect weapon"
		c.Actions = []commands.Action{
			driver.PressKey("q"),
		}
	}

	// Next unit
	if match, _ := regexp.MatchString("^next$", t); match {
		c.Description = "Select next unit"
		c.Actions = []commands.Action{
			driver.PressKey("tab"),
		}
	}

	// Next unit (batch)
	if r := regexp.MustCompile("^next ([2-9])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		count, _ := strconv.Atoi(ss[1])

		c.Description = fmt.Sprintf("Select next unit #%d", count)
		c.Actions = []commands.Action{
			func() {
				for i := 0; i < count; i++ {
					driver.PressKey("tab")()
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
		c.Actions = []commands.Action{
			driver.Mouse(x, y),
		}
	}

	// Click
	// Moves to coordinates and clicks
	if r := regexp.MustCompile("^click ([0-9]*) ([0-9]*)$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x, _ := strconv.Atoi(ss[1])
		y, _ := strconv.Atoi(ss[2])

		c.Description = fmt.Sprintf("Click the mouse at x: %d, y: %d", x, y)
		c.Actions = []commands.Action{
			driver.Mouse(x, y),
			driver.Click(),
		}
	}

	// Mouse grid
	// Moves the mouse to the grid reference
	if r := regexp.MustCompile("^mouse ([a-h])([1-8]*)$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x := strings.ToUpper(ss[1])
		y := ss[2]

		c.Description = fmt.Sprintf("Move the mouse to tile %s%s", x, y)
		c.Actions = []commands.Action{
			driver.MouseGrid(x, y),
		}
	}

	// Click grid
	// Moves the mouse to the grid reference and clicks
	if r := regexp.MustCompile("^click ([a-h])([1-8]*)$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x := strings.ToUpper(ss[1])
		y := ss[2]

		c.Description = fmt.Sprintf("Click the mouse at tile %s%s", x, y)
		c.Actions = []commands.Action{
			driver.MouseGrid(strings.ToUpper(ss[1]), ss[2]),
			driver.Click(),
		}
	}

	// Select unit
	// Select a specific unit by hotkey
	if r := regexp.MustCompile("^select (mech|deployed|mission) ([1-3])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		c.Description = fmt.Sprintf("Select %s unit #%s", ss[1], ss[2])
		c.Actions = []commands.Action{
			driver.PressKey(getUnitKey(ss[1], ss[2])),
		}
	}

	// Move unit
	// Move a specific unit to a given map coordinate
	if r := regexp.MustCompile("^move (mech|deployed|mission) ([1-3]) ([a-h])([1-8])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x := strings.ToUpper(ss[3])
		y := ss[4]

		c.Description = fmt.Sprintf("Move %s unit #%s to %s%s", ss[1], ss[2], x, y)
		c.Actions = []commands.Action{
			driver.PressKey(getUnitKey(ss[1], ss[2])),
			driver.MouseGrid(x, y),
			driver.Click(),
		}
	}

	// Select weapon
	// Select a specific weapon by hotkey
	if r := regexp.MustCompile("^weapon ([1-2])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		c.Description = fmt.Sprintf("Arm weapon #%s", ss[1])
		c.Actions = []commands.Action{
			driver.PressKey(ss[1]),
		}
	}

	// Attack
	// Attack with a unit using a weapon at a given tile
	if r := regexp.MustCompile("^attack (mech|deployed|mission) ([1-3]) ([1-2]) ([a-h])([1-8])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x := strings.ToUpper(ss[4])
		y := ss[5]

		c.Description = fmt.Sprintf("Attacking with %s unit #%s using weapon %s on tile %s%s", ss[1], ss[2], ss[3], x, y)
		c.Actions = []commands.Action{
			driver.PressKey(getUnitKey(ss[1], ss[2])),
			driver.PressKey(ss[3]),
			driver.MouseGrid(x, y),
			driver.Click(),
		}
	}

	// Repair
	// Repair a unit at a given tile
	if r := regexp.MustCompile("^repair ([1-3]) ([a-h])([1-8])$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x := strings.ToUpper(ss[2])
		y := ss[3]

		c.Description = fmt.Sprintf("Repairing mech #%d at tile %d%d", ss[1], x, y)
		c.Actions = []commands.Action{
			driver.PressKey(getUnitKey("mech", ss[1])),
			driver.PressKey("r"),
			driver.MouseGrid(x, y),
			driver.Click(),
		}
	}

	// On-off commands

	// Toggle info
	// Toggle the info tooltip
	if r := regexp.MustCompile("^info (on|off)$"); r.MatchString(t) {
		state := r.FindStringSubmatch(t)[1]
		c.Description = fmt.Sprintf("Turning info tooltip %s", state)
		c.Actions = []commands.Action{
			func() {
				driver.ToggleKey("control", state == "on")
			},
		}
	}

	// Toggle turn order
	// Toggle the turn order tooltips
	if r := regexp.MustCompile("^order (on|off)$"); r.MatchString(t) {
		state := r.FindStringSubmatch(t)[1]
		c.Description = fmt.Sprintf("Turning turn order tooltips %s", state)
		c.Actions = []commands.Action{
			func() {
				driver.ToggleKey("alt", state == "on")
			},
		}
	}

	// Calibration commands

	// Calibrate
	// Moves the mouse from the top left to bottom right corner and around the grid. Used for offset calibration
	if match, _ := regexp.MatchString("^calibrate$", t); match {
		c.Description = "Calibrating"
		c.Actions = []commands.Action{
			driver.Mouse(0, 0),
			driver.Mouse(1280, 720),
			driver.MouseGrid("A", "1"),
			driver.MouseGrid("H", "1"),
			driver.MouseGrid("H", "8"),
			driver.MouseGrid("A", "8"),
		}
	}

	if len(c.Actions) == 0 {
		return nil
	}

	return c
}

// getUnitKey returns a unit key by type and number
func getUnitKey(t, n string) string {
	keys := map[string]string{
		"mech1":     "a",
		"mech2":     "s",
		"mech3":     "d",
		"deployed1": "f",
		"deployed2": "g",
		"deployed3": "h",
		"mission1":  "z",
		"mission2":  "x",
	}

	if val, ok := keys[fmt.Sprintf("%s%s", t, n)]; ok {
		return val
	}

	return "a"
}
