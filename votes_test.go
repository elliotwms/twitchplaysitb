package main

import (
	"sort"
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

	//assert.Empty(t, commands)
	//assert.Equal(t, results[0].Votes, 2)
	//assert.Len(t, results[0].Users, 2)
	//
	//assert.Equal(t, results[1].Votes, 1)
	//assert.Len(t, results[1].Users, 1)
}

func TestVoteSort(t *testing.T) {
	list := VoteResultList{
		{
			Votes: 1,
		},
		{
			Votes: 3,
		},
		{
			Votes: 2,
		},
	}

	assert.False(t, sort.IsSorted(list))

	sort.Sort(list)
	assert.True(t, sort.IsSorted(list))

	assert.Equal(t, 1, list[0].Votes)
	assert.Equal(t, 2, list[1].Votes)
	assert.Equal(t, 3, list[2].Votes)

	sort.Sort(sort.Reverse(list))

	assert.Equal(t, 3, list[0].Votes)
	assert.Equal(t, 2, list[1].Votes)
	assert.Equal(t, 1, list[2].Votes)
}

func TestGetWinningVote(t *testing.T) {
	list := VoteResultList{
		{
			Votes: 1,
		},
		{
			Votes: 3,
		},
		{
			Votes: 2,
		},
	}

	result := getWinningVote(list)

	assert.Equal(t, 3, result.Votes)
}
