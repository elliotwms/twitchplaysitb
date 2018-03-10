package main

import (
	"fmt"
	"github.com/gempir/go-twitch-irc"
	"github.com/go-vgo/robotgo"
	"os"
	"time"
	"errors"
)

const xOffset = 0
const yOffset = 45

const gameWidth = 1280
const gameHeight = 720


func main() {
	//setUpGame()

	username, token, channel, err := getCredentials()

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	c := twitch.NewClient(username, token)

	// todo Create and work channel for messages

	c.OnNewMessage(handleMessage(c))

	c.Join(channel)

	err = c.Connect()

	if err != nil {
		panic(err)
	}

}

func handleMessage(c *twitch.Client) func(channel string, user twitch.User, message twitch.Message) {
	return func(channel string, user twitch.User, message twitch.Message) {
		fmt.Printf("%s: %s\n", user.Username, message.Text)

		// todo Check user hasn't already submitted a command

		// todo Parse command
		command := Parse(message.Text)

		// Just do the commands for now. We'll do queueing later...
		if command != nil {
			// Broadcast the command
			c.Say(channel, command.Description)
			fmt.Println(command.Description)

			for _, a := range command.Actions {
				a.Do()
				time.Sleep(3 * time.Second)
			}
		}

		// todo If valid command, add to command queue
	}
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

func setUpGame() {
	ok := robotgo.ShowAlert("Into The Twitch", "Move your app to the top left of the screen and set the resolution to \"Default Windowed\"")

	if ok == 1 {
		os.Exit(1)
	}

	robotgo.ActiveName("IntoTheBreach")

	bmp := robotgo.CaptureScreen(xOffset, yOffset, 1280, 720)
	robotgo.SaveBitmap(bmp, "test.png")

	// Mouse over New Game
	robotgo.MoveMouseSmooth(xOffset+150, yOffset+280)

	// Focus window
	robotgo.Click()
	time.Sleep(1 * time.Second)

	// New Game
	robotgo.Click()
}
