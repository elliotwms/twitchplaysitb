package main

import (
	"regexp"
	"github.com/go-vgo/robotgo"
	"strconv"
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
		robotgo.MoveMouseSmooth(x, y)
	}
}