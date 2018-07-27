package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	c := Parse("click")

	assert.NotNil(t, c)
}
