package main

import (
	"regexp"
	"github.com/go-vgo/robotgo"
	"strconv"
	"strings"
	"time"
)

type Command struct {
	Text    string
	Actions []Action
}

type Action struct {
	Do func()
}

func Parse(t string) *Command {
	c := &Command{
		Text: t,
	}

	// Simple commands

	// Click
	if match, _ := regexp.MatchString("^click$", t); match {
		c.Actions = []Action{
			{
				Do: click(),
			},
		}

		return c
	}

	// End turn
	if match, _ := regexp.MatchString("^endturn$", t); match {
		c.Actions = []Action{
			{
				Do: func() {
					robotgo.KeyTap("space")
				},
			},
		}

		return c
	}

	// Undo move
	if match, _ := regexp.MatchString("^undo$", t); match {
		c.Actions = []Action{
			{
				Do: func() {
					robotgo.KeyTap("shift")
				},
			},
		}

		return c
	}

	// Reset turn
	if match, _ := regexp.MatchString("^reset$", t); match {
		c.Actions = []Action{
			{
				Do: func() {
					robotgo.KeyTap("backspace")
				},
			},
		}

		return c
	}

	// Deselect weapon
	if match, _ := regexp.MatchString("^deselect|disarm$", t); match {
		c.Actions = []Action{
			{
				Do: func() {
					robotgo.KeyTap("q")
				},
			},
		}

		return c
	}

	// Next unit
	if match, _ := regexp.MatchString("^next$", t); match {
		c.Actions = []Action{
			{
				Do: func() {
					robotgo.KeyTap("tab")
				},
			},
		}

		return c
	}

	// Next unit (batch
	if r := regexp.MustCompile("^next ([2-9])$"); r.MatchString(t) {

		ss := r.FindStringSubmatch(t)

		count, _ := strconv.Atoi(ss[1])

		c.Actions = []Action{
			{
				Do: func() {
					for i := 0; i < count; i++ {
						robotgo.KeyTap("tab")
						time.Sleep(1 * time.Second)
					}
				},
			},
		}

		return c
	}

	// Less simple commands

	// Mouse
	// Moves the mouse to the given coordinates
	if r := regexp.MustCompile("^mouse ([0-9]*) ([0-9]*)$"); r.MatchString(t) {

		ss := r.FindStringSubmatch(t)

		x, _ := strconv.Atoi(ss[1])
		y, _ := strconv.Atoi(ss[2])

		c.Actions = []Action{
			{
				Do: mouse(x, y),
			},
		}

		return c
	}

	// Click
	// Moves to coordinates and clicks
	if r := regexp.MustCompile("^click ([0-9]*) ([0-9]*)$"); r.MatchString(t) {
		ss := r.FindStringSubmatch(t)

		x, _ := strconv.Atoi(ss[1])
		y, _ := strconv.Atoi(ss[2])

		c.Actions = []Action{
			{
				Do: mouse(x, y),
			},
			{
				Do: click(),
			},
		}

		return c
	}

	// Mouse grid
	// Moves the mouse to the grid reference
	if r := regexp.MustCompile("^mouse ([A-H])([1-8]*)$"); r.MatchString(t) {

		ss := r.FindStringSubmatch(t)

		c.Actions = []Action{
			{
				Do: mouseGrid(strings.ToUpper(ss[1]), ss[2]),
			},
		}

		return c
	}

	// Click grid
	// Moves the mouse to the grid reference and clicks
	if r := regexp.MustCompile("^click ([A-H])([1-8]*)$"); r.MatchString(t) {

		ss := r.FindStringSubmatch(t)

		c.Actions = []Action{
			{
				Do: mouseGrid(strings.ToUpper(ss[1]), ss[2]),
			},
			{
				Do: click(),
			},
		}

		return c
	}

	// Select unit
	// Select a specific unit by hotkey
	if r := regexp.MustCompile("^select (mech|deployed|mission) ([1-3])$"); r.MatchString(t) {

		ss := r.FindStringSubmatch(t)

		c.Actions = []Action{
			{
				Do: selectUnit(ss[1], ss[2]),
			},
		}

		return c
	}

	// Select weapon
	// Select a specific weapon by hotkey
	if r := regexp.MustCompile("^weapon ([1-2])$"); r.MatchString(t) {

		ss := r.FindStringSubmatch(t)

		c.Actions = []Action{
			{
				Do: selectWeapon(ss[1]),
			},
		}

		return c
	}

	// Attack
	// Attack with a unit using a weapon at a given tile
	if r := regexp.MustCompile("^attack (mech|deployed|mission) ([1-3]) ([1-2]) ([A-H])([1-8])$"); r.MatchString(t) {

		ss := r.FindStringSubmatch(t)

		c.Actions = []Action{
			{
				Do: selectUnit(ss[1], ss[2]),
			},
			{
				Do: selectWeapon(ss[3]),
			},
			{
				Do: mouseGrid(ss[4], ss[5]),
			},
			{
				Do: click(),
			},
		}

		return c
	}

	// Repair
	// Repair a unit at a given tile
	if r := regexp.MustCompile("^repair ([1-3]) ([A-H])([1-8])$"); r.MatchString(t) {

		ss := r.FindStringSubmatch(t)

		c.Actions = []Action{
			{
				Do: selectUnit("mech", ss[1]),
			},
			{
				Do: repair(),
			},
			{
				Do: mouseGrid(ss[2], ss[3]),
			},
			{
				Do: click(),
			},
		}

		return c
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

		c.Actions = []Action{
			{
				Do: func() {
					robotgo.KeyToggle("control", state)
				},
			},
		}

		return c
	}

	// Toggle turn order
	// Toggle the turn order tooltips
	if r := regexp.MustCompile("^order (on|off)$"); r.MatchString(t) {

		ss := r.FindStringSubmatch(t)

		state := "up"

		if ss[1] == "on" {
			state = "down"
		}

		c.Actions = []Action{
			{
				Do: func() {
					robotgo.KeyToggle("alt", state)
				},
			},
		}

		return c
	}

	// Calibration commands

	// Calibrate
	// Moves the mouse from the top left to bottom right corner and around the grid. Used for offset calibration
	if match, _ := regexp.MatchString("^calibrate$", t); match {
		c.Actions = []Action{
			{
				Do: mouse(0, 0),
			},
			{
				Do: mouse(1280, 720),
			},
			{
				Do: mouseGrid("A", "1"),
			},
			{
				Do: mouseGrid("H", "1"),
			},
			{
				Do: mouseGrid("H", "8"),
			},
			{
				Do: mouseGrid("A", "8"),
			},
		}

		return c
	}

	return nil
}

// Commands
// These commands wrap the action into a callable zero-argument function

func click() func() {
	return func() {
		robotgo.Click()
	}
}

func mouse(x int, y int) func() {
	return func() {
		robotgo.MoveMouseSmooth(x+xOffset, y+yOffset)
	}
}

func mouseGrid(a string, n string) func() {
	x, y := GetCoordinates(a, n)

	return mouse(x, y)
}

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

func selectWeapon(n string) func() {
	return func() {
		robotgo.KeyTap(n)
	}
}

func repair() func() {
	return func() {
		robotgo.KeyTap("r") // Repair
	}
}
