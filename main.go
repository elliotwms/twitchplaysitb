package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gempir/go-twitch-irc"
	"github.com/go-vgo/robotgo"
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

	cq := make(map[string]*Command)

	go workCommands(cq)

	c.OnNewMessage(handleMessage(c, cq))

	c.Join(channel)

	err = c.Connect()

	if err != nil {
		panic(err)
	}

}

func handleMessage(c *twitch.Client, cq map[string]*Command) func(channel string, user twitch.User, message twitch.Message) {
	return func(channel string, user twitch.User, message twitch.Message) {
		fmt.Printf("%s: %s\n", user.Username, message.Text)

		command := Parse(message.Text)

		if command != nil {
			// Add the command to the command queue
			cq[user.Username] = command
		}
	}
}

func workCommands(cq map[string]*Command) {
	for {
		count := len(cq)

		if len(cq) > 0 {
			results := tallyVotes(cq)

			fmt.Printf("Processing %d votes for %d actions", count, len(results))

			for _, r := range results {
				fmt.Printf("%s: %d votes\n", r.Command.Description, r.Votes)
			}

			result := getWinningVote(results)

			for _, action := range result.Command.Actions {
				action.Do()
			}

		} else {
			fmt.Println("No votes to process")
		}

		time.Sleep(30 * time.Second)
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
