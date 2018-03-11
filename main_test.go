package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCredentials(t *testing.T) {
	_, _, _, err := getCredentials()

	assert.EqualError(t, err, "missing username")

	os.Setenv("TWITCH_USERNAME", "foo")

	username, _, _, err := getCredentials()

	assert.Equal(t, "foo", username)
	assert.EqualError(t, err, "missing token")

	os.Setenv("TWITCH_TOKEN", "bar")

	username, token, _, err := getCredentials()

	assert.Equal(t, "bar", token)
	assert.EqualError(t, err, "missing channel")

	os.Setenv("TWITCH_CHANNEL", "baz")

	username, token, channel, err := getCredentials()

	assert.Equal(t, "baz", channel)
	assert.NoError(t, err)
}
