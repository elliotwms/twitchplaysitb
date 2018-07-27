package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/elliotwms/twitchplaysitb/commands"
	"github.com/gempir/go-twitch-irc"
	"github.com/go-vgo/robotgo"
)

const xOffset = 0
const yOffset = 45

const gameWidth = 1280
const gameHeight = 720

var verbose = false
var turnTime = 60
var pid int32

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Missing PID argument 1")
		os.Exit(1)
	}

	if p, err := strconv.ParseInt(os.Args[1], 10, 32); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	} else {
		pid = int32(p)
	}

	if ok, err := robotgo.PidExists(pid); !ok {
		fmt.Printf("Could not find process with ID %s", err.Error())
		os.Exit(1)
	}

	// Make the game window active
	robotgo.ActivePID(int32(pid))

	username, token, channel, err := getCredentials()

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	c := twitch.NewClient(username, token)

	cq := make(map[string]*commands.Command)

	go workCommands(cq, c, channel, pid)

	c.OnNewMessage(handleMessage(c, cq))

	c.Join(channel)

	err = c.Connect()

	if err != nil {
		panic(err)
	}
}

func handleMessage(c *twitch.Client, cq map[string]*commands.Command) func(channel string, user twitch.User, message twitch.Message) {
	return func(channel string, user twitch.User, message twitch.Message) {
		fmt.Printf("%s: %s\n", user.Username, message.Text)

		command := Parse(message.Text)

		if command != nil {
			// Add the command to the command queue
			cq[user.Username] = command
		}
	}
}

func workCommands(cq map[string]*commands.Command, c *twitch.Client, channel string, pid int32) {
	for {
		if ok, _ := robotgo.PidExists(pid); !ok {
			break
		}

		count := len(cq)

		if len(cq) > 0 {
			results := tallyVotes(cq)

			if verbose {
				c.Say(channel, fmt.Sprintf("Processing %d votes for %d actions\n", count, len(results)))

				for _, r := range results {
					c.Say(channel, fmt.Sprintf("%s: %d votes\n", r.Command.Description, r.Votes))
				}
			}

			result := getWinningVote(results)

			c.Say(channel, fmt.Sprintf("Result: %s\n", result.Command.Description))

			for _, action := range result.Command.Actions {
				action()
			}

		} else {
			c.Say(channel, "No votes to process")
		}

		t := turnTime // todo adjust time

		c.Say(channel, fmt.Sprintf("Next command in %d seconds", t))

		time.Sleep(time.Duration(t) * time.Second)
	}

	c.Say(channel, fmt.Sprintf("You broke it!"))
	os.Exit(1)
}

func getCredentials() (username string, token string, channel string, err error) {
	username, found := os.LookupEnv("TWITCH_USERNAME")

	if !found {
		return username, token, channel, errors.New("missing username")
	}

	token, found = os.LookupEnv("TWITCH_TOKEN")

	if !found {
		return username, token, channel, errors.New("missing token")
	}

	channel, found = os.LookupEnv("TWITCH_CHANNEL")

	if !found {
		return username, token, channel, errors.New("missing channel")
	}

	return username, token, channel, nil
}
