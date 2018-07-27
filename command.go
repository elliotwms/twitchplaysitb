package main

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Command struct {
	Text        string   // The raw command input
	Description string   // The command description
	Actions     []Action // The actions to perform when executing the command
}

// GetHash hashes the command description (which should be unique depending on the command arguments)
func (c *Command) GetHash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(c.Description)))
}

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
			pressKey("space"),
		}
	}

	// Undo move
	if match, _ := regexp.MatchString("^undo$", t); match {
		c.Description = "Undo move"
		c.Actions = []Action{
			pressKey("shift"),
		}
	}

	// Reset turn
	if match, _ := regexp.MatchString("^reset$", t); match {
		c.Description = "Reset turn"
		c.Actions = []Action{
			pressKey("backspace"),
		}
	}

	// Deselect weapon
	if match, _ := regexp.MatchString("^deselect|disarm$", t); match {
		c.Description = "Deselect weapon"
		c.Actions = []Action{
			pressKey("q"),
		}
	}

	// Next unit
	if match, _ := regexp.MatchString("^next$", t); match {
		c.Description = "Select next unit"
		c.Actions = []Action{
			pressKey("tab"),
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
					pressKey("tab")()
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
		state := r.FindStringSubmatch(t)[1]
		c.Description = fmt.Sprintf("Turning info tooltip %s", state)
		c.Actions = []Action{
			func() {
				toggleKey("control", state == "on")
			},
		}
	}

	// Toggle turn order
	// Toggle the turn order tooltips
	if r := regexp.MustCompile("^order (on|off)$"); r.MatchString(t) {
		state := r.FindStringSubmatch(t)[1]
		c.Description = fmt.Sprintf("Turning turn order tooltips %s", state)
		c.Actions = []Action{
			func() {
				toggleKey("alt", state == "on")
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
