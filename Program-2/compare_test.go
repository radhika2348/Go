package main

import "testing"
import "github.com/stretchr/testify/assert"


func TestMax(t *testing.T) {
	assert.Equal(t, true, Max(9, 4))

	assert.Equal(t, false, Max(2, 7))

	assert.Equal(t, false, Max(5, 5))
}
