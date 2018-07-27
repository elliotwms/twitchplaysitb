package commands

import (
	"crypto/md5"
	"fmt"
	"regexp"

	"github.com/elliotwms/twitchplaysitb/drivers"
)

type Command struct {
	Text        string   // The raw command input
	Description string   // The command description
	Actions     []Action // The actions to perform when executing the command
}

type Action func()
type CommandBuilder func(d drivers.Driver, a []string) *Command
type Dictionary map[*regexp.Regexp]CommandBuilder

// GetHash hashes the command description (which should be unique depending on the command arguments)
func (c *Command) GetHash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(c.Description)))
}
