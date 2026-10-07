package main

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnSquareOfPositiveNumber(t *testing.T) {
	randomNum := rand.Intn(100)
	expected := randomNum * randomNum
	actual := square(randomNum)

	assert.Equal(t, expected, actual)
}

func TestShouldReturnSquareOfNegativeNumber(t *testing.T) {
	randomNum := -rand.Intn(100)
	expected := randomNum * randomNum
	actual := square(randomNum)

	assert.Equal(t, expected, actual)
}

func TestShouldReturnSquareOfZero(t *testing.T) {
	randomNum := 0
	expected := 0
	actual := square(randomNum)

	assert.Equal(t, expected, actual)
}