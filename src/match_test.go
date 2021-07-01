package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIndividualMatch(t *testing.T) {
	individualMatch := "A;8;2;AUT-1;Furl;h;M;M;;T;;;BEL-4;Somame"

	match := ParseMatch(individualMatch).(Match)
	assert.IsType(t, Match{}, match)
	assert.Equal(t, "A", match.Shiaijo)
	assert.Equal(t, 8, match.Pool)
	assert.Equal(t, 2, match.Fight)
	assert.Equal(t, "AUT-1", match.NumberTareWhite)
	assert.Equal(t, "Furl", match.NameTareWhite)

}

func TestParseTeamMatch(t *testing.T) {
	teamMatch := "C;5;3;AUT-1;Mayer;h;;M;;K;;;BEL-3;Somame;4;Austria;2;5;1;3;1;Belgium"

	match := ParseMatch(teamMatch).(TeamMatch)
	assert.IsType(t, TeamMatch{}, match)
	assert.Equal(t, "C", match.Shiaijo)
	assert.Equal(t, 5, match.Pool)
	assert.Equal(t, 3, match.Fight)
}
