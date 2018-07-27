package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_GetHash(t *testing.T) {
	cases := []struct {
		i string
		o string
	}{
		{i: "test", o: "098f6bcd4621d373cade4e832627b4f6"},
		{i: "this is an actual description", o: "f24e393ecc19b8121e4ea4ed0059df62"},
		{i: "Attacking with mech unit #1 using weapon 1 on tile A1", o: "da271dbfee3ebc189803e3f6743b3cc0"},
	}

	for _, tt := range cases {
		c := Command{
			Description: tt.i,
		}

		assert.Equal(t, tt.o, c.GetHash())
	}
}
