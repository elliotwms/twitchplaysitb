package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTallyVotes(t *testing.T) {
	commands := map[string]*Command{
		"foo": {Description: "click"},
		"bar": {Description: "mouse 100 100"},
		"baz": {Description: "click"},
	}

	results := tallyVotes(commands)

	assert.Len(t, results, 2)

	// todo this test fails sometimes to to map randomisation

	assert.Empty(t, commands)
	assert.Equal(t, results[0].Votes, 2)
	assert.Len(t, results[0].Users, 2)

	assert.Equal(t, results[1].Votes, 1)
	assert.Len(t, results[1].Users, 1)
}
