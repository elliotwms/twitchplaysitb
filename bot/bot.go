package bot

import (
	"errors"

	"github.com/elliotwms/twitchplaysitb/commands"
	"github.com/elliotwms/twitchplaysitb/drivers"
)

type Bot struct {
	d    drivers.Driver
	dict commands.Dictionary
}

func New(driver drivers.Driver, dictionary commands.Dictionary) *Bot {
	return &Bot{
		d:    driver,
		dict: dictionary,
	}
}

func (b *Bot) ResolveCommand(input string) (*commands.Command, error) {
	for r, c := range b.dict {
		if r.MatchString(input) {
			return c(b.d, r.FindStringSubmatch(input)), nil
		}
	}

	return nil, errors.New("command not found")
}
