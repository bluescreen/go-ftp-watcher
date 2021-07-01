package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	assert.Equal(t, 123, 1234, "they should be equal1")
}

func TestSomethingElse(t *testing.T) {
	assert.Equal(t, 123, 123, "they should be equal")
}
