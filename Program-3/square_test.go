package main

import "testing"
import "github.com/stretchr/testify/assert"


func TestSquare(t *testing.T) {
	assert.Equal(t, 25, Square(5))

	assert.Equal(t, 9, Square(-3))

	assert.Equal(t, 0, Square(0))
}
