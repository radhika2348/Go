package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMultipleOf7(t *testing.T) {
	actualTrue := MultipleOf7(14)
	assert.Equal(t, 98, actualTrue) 

	actualFalse := MultipleOf7(10)
	assert.Equal(t, 70, actualFalse) 
}
