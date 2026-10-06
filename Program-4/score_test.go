package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestIsGreaterThan40(t *testing.T) {
	assert.Equal(t, true, IsGreaterThan40(45))

	//assert.Equal(t, false, IsGreaterThan40(20))

	//assert.Equal(t, true, IsGreaterThan40(40))
}

func TestIsLesserThan40(t *testing.T) {
	assert.Equal(t, false, IsGreaterThan40(20))


}

func TestIsEqualTo40(t *testing.T) {
	assert.Equal(t, true, IsGreaterThan40(40))


}