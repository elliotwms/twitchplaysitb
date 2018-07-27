package bot

import (
	"regexp"
	"testing"

	"github.com/elliotwms/twitchplaysitb/commands"
	"github.com/elliotwms/twitchplaysitb/drivers"
	"github.com/stretchr/testify/assert"
)

func TestBot_ResolveCommand(t *testing.T) {
	actionsCalled := 0

	dict := commands.Dictionary{
		regexp.MustCompile("^foo$"): func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Actions: []commands.Action{
					func() {
						assert.Len(t, a, 1)
						actionsCalled++
					},
				},
			}
		},
		regexp.MustCompile("^bar (baz|bux)$"): func(d drivers.Driver, a []string) *commands.Command {
			return &commands.Command{
				Actions: []commands.Action{
					func() {
						assert.Len(t, a, 2)
						actionsCalled++
					},
				},
			}
		},
	}

	b := New(&drivers.NoOpDriver{}, dict)

	c, err := b.ResolveCommand("foo")

	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Len(t, c.Actions, 1)

	for _, a := range c.Actions {
		a()
	}

	c, err = b.ResolveCommand("bar")

	assert.Error(t, err)
	assert.Nil(t, c)

	c, err = b.ResolveCommand("bar baz")

	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Len(t, c.Actions, 1)
	for _, a := range c.Actions {
		a()
	}

	// Check all the assertions in the callbacks were actually called
	assert.Equal(t, 2, actionsCalled)
}
