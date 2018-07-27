package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	c := Parse("click")

	assert.NotNil(t, c)
}

func TestGetUnitKey(t *testing.T) {
	tests := []struct {
		description string
		t           string
		n           string
		expected    string
	}{
		{description: "Valid type and number", t: "mech", n: "1", expected: "a"},
		{description: "Another valid type and number", t: "deployed", n: "2", expected: "g"},
		{description: "Valid type only", t: "mission", n: "3", expected: "a"},
		{description: "Empty type and number", t: "", n: "", expected: "a"},
	}

	for k, tt := range tests {
		t.Logf("Test #%d: %s", k+1, tt.description)
		assert.Equal(t, tt.expected, getUnitKey(tt.t, tt.n))
	}
}
