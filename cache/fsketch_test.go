package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFSketchWithNegativeOrZero(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, NewFSketch(-1).length)
	assert.Equal(1, NewFSketch(0).length)
}

func TestNewFSketch(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(8, NewFSketch(16).length)
	assert.Equal(32, NewFSketch(100).length)
}

