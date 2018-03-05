package main

import (
	"fmt"
	"github.com/gempir/go-twitch-irc"
	"github.com/go-vgo/robotgo"
	"os"
	"time"
)

const xOffset = 0
const yOffset = 45

func main() {
	ok := robotgo.ShowAlert("Into The Twitch", "Move your app to the top left of the screen and set the resolution to \"Default Windowed\"")

	if ok == 1 {
		os.Exit(1)
	}

	setUpGame()
	play()
}

func setUpGame() {
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

func play() {
	username := os.Getenv("USERNAME")
	password := os.Getenv("TOKEN")
	channel := os.Getenv("CHANNEL")

	client := twitch.NewClient(username, password)

	client.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		fmt.Println(message.Text)
	})

	err := client.Connect()

	if err != nil {
		panic(err)
	}

	client.Join(channel)
	client.Say(channel, "Hello World")
}
