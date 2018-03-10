package main

import (
	"regexp"
	"github.com/go-vgo/robotgo"
	"strconv"
	"strings"
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

	// Click
	if match, _ := regexp.MatchString("^click$", t); match {
		c.Actions = []Action{
			{
				Do: click(),
			},
		}

		return c
	}

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
