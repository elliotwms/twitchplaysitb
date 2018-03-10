package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetCoordinates(t *testing.T) {

	cases := []struct {
		a string
		n string
		x int
		y int
	}{
		{"A", "1", 640, 667},
		{"H", "1", 1032, 380},
		{"A", "8", 248, 380},
		{"H", "8", 640, 93},
	}

	for _, tt := range cases {
		x, y := GetCoordinates(tt.a, tt.n)

		assert.Equal(t, tt.x, x)
		assert.Equal(t, tt.y, y)
	}
}
