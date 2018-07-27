package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elliotwms/twitchplaysitb/bot"
	"github.com/elliotwms/twitchplaysitb/commands"
	"github.com/elliotwms/twitchplaysitb/drivers"
)

// buildDictionary builds a dictionary of command definitions for interacting with Into The Breach
func buildDictionary(driver drivers.Driver) commands.Dictionary {
	list := map[string]commands.CommandBuilder{
		"^click$": func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Description: "Click the mouse",
				Actions: []commands.Action{
					driver.Click(),
				},
			}
		},
		"^endturn$": func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Description: "End turn",
				Actions: []commands.Action{
					driver.PressKey("space"),
				},
			}
		},
		"^undo$": func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Description: "Undo move",
				Actions: []commands.Action{
					driver.PressKey("shift"),
				},
			}
		},
		"^reset$": func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Description: "Reset turn",
				Actions: []commands.Action{
					driver.PressKey("backspace"),
				},
			}
		},
		"^deselect|disarm$": func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Description: "Deselect weapon",
				Actions: []commands.Action{
					driver.PressKey("q"),
				},
			}
		},
		"^next$": func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Description: "Select next unit",
				Actions: []commands.Action{
					driver.PressKey("tab"),
				},
			}
		},
		"^next ([2-9])$": func(d drivers.Driver, a []string) *commands.Command {
			count, _ := strconv.Atoi(a[1])
			return &commands.Command{
				Description: fmt.Sprintf("Select next unit #%d", count),
				Actions: []commands.Action{
					func() {
						for i := 0; i < count; i++ {
							driver.PressKey("tab")()
							time.Sleep(1 * time.Second)
						}
					},
				},
			}
		},
		"^mouse ([0-9]*) ([0-9]*)$": func(d drivers.Driver, a []string) *commands.Command {
			x, _ := strconv.Atoi(a[1])
			y, _ := strconv.Atoi(a[2])

			return &commands.Command{
				Description: fmt.Sprintf("Move the mouse to x: %d, y: %d", x, y),
				Actions: []commands.Action{
					driver.Mouse(x, y),
				},
			}
		},
		"^click ([0-9]*) ([0-9]*)$": func(d drivers.Driver, a []string) *commands.Command {
			x, _ := strconv.Atoi(a[1])
			y, _ := strconv.Atoi(a[2])

			return &commands.Command{
				Description: fmt.Sprintf("Click the mouse at x: %d, y: %d", x, y),
				Actions: []commands.Action{
					driver.Mouse(x, y),
					driver.Click(),
				},
			}
		},
		"^mouse ([a-h])([1-8]*)$": func(d drivers.Driver, a []string) *commands.Command {
			x := strings.ToUpper(a[1])
			y := a[2]

			return &commands.Command{
				Description: fmt.Sprintf("Move the mouse to tile %s%s", x, y),
				Actions: []commands.Action{
					driver.Mouse(GetCoordinates(x, y)),
				},
			}
		},
		"^click ([a-h])([1-8]*)$": func(d drivers.Driver, a []string) *commands.Command {
			x := strings.ToUpper(a[1])
			y := a[2]

			return &commands.Command{
				Description: fmt.Sprintf("Click the mouse at tile %s%s", x, y),
				Actions: []commands.Action{
					driver.Mouse(GetCoordinates(strings.ToUpper(x), y)),
					driver.Click(),
				},
			}
		},
		"^select (mech|deployed|mission) ([1-3])$": func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Description: fmt.Sprintf("Select %s unit #%s", a[1], a[2]),
				Actions: []commands.Action{
					driver.PressKey(getUnitKey(a[1], a[2])),
				},
			}
		},
		"^move (mech|deployed|mission) ([1-3]) ([a-h])([1-8])$": func(d drivers.Driver, a []string) *commands.Command {
			x := strings.ToUpper(a[3])
			y := a[4]

			return &commands.Command{
				Description: fmt.Sprintf("Move %s unit #%s to %s%s", a[1], a[2], x, y),
				Actions: []commands.Action{
					driver.PressKey(getUnitKey(a[1], a[2])),
					driver.Mouse(GetCoordinates(x, y)),
					driver.Click(),
				},
			}
		},
		"^weapon ([1-2])$": func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Description: fmt.Sprintf("Arm weapon #%s", a[1]),
				Actions: []commands.Action{
					driver.PressKey(a[1]),
				},
			}
		},
		"^attack (mech|deployed|mission) ([1-3]) ([1-2]) ([a-h])([1-8])$": func(d drivers.Driver, a []string) *commands.Command {
			x := strings.ToUpper(a[4])
			y := a[5]
			return &commands.Command{
				Description: fmt.Sprintf("Attacking with %s unit #%s using weapon %s on tile %s%s", a[1], a[2], a[3], x, y),
				Actions: []commands.Action{
					driver.PressKey(getUnitKey(a[1], a[2])),
					driver.PressKey(a[3]),
					driver.Mouse(GetCoordinates(x, y)),
					driver.Click(),
				},
			}
		},
		"^repair ([1-3]) ([a-h])([1-8])$": func(d drivers.Driver, a []string) *commands.Command {
			x := strings.ToUpper(a[2])
			y := a[3]
			return &commands.Command{
				Description: fmt.Sprintf("Repairing mech #%d at tile %d%d", a[1], x, y),
				Actions: []commands.Action{
					driver.PressKey(getUnitKey("mech", a[1])),
					driver.PressKey("r"),
					driver.Mouse(GetCoordinates(x, y)),
					driver.Click(),
				},
			}
		},
		"^info (on|off)$": func(d drivers.Driver, a []string) *commands.Command {
			state := a[1]
			return &commands.Command{
				Description: fmt.Sprintf("Turning info tooltip %s", state),
				Actions: []commands.Action{
					func() {
						driver.ToggleKey("control", state == "on")
					},
				},
			}
		},
		"^order (on|off)$": func(d drivers.Driver, a []string) *commands.Command {
			state := a[1]
			return &commands.Command{
				Description: fmt.Sprintf("Turning turn order tooltips %s", state),
				Actions: []commands.Action{
					func() {
						driver.ToggleKey("alt", state == "on")
					},
				},
			}
		},
		"^calibrate$": func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Description: "Calibrating",
				Actions: []commands.Action{
					driver.Mouse(0, 0),
					driver.Mouse(1280, 720),
					driver.Mouse(GetCoordinates("A", "1")),
					driver.Mouse(GetCoordinates("H", "1")),
					driver.Mouse(GetCoordinates("H", "8")),
					driver.Mouse(GetCoordinates("A", "8")),
				},
			}
		},
	}

	dict := make(commands.Dictionary)

	for key, value := range list {
		dict[regexp.MustCompile(key)] = value
	}

	return dict
}

// Parse parses the text input and returns a command if a match is found or nil
func Parse(t string) *commands.Command {
	driver := drivers.NewRobotGoDriver(gameWidth, gameHeight, xOffset, yOffset)
	// lower case all commands to normalize them
	t = strings.ToLower(t)
	b := bot.New(driver, buildDictionary(driver))
	c, err := b.ResolveCommand(t)

	if err != nil {
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
