package cache

import (
	"encoding/binary"
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
	assert.Equal(2, NewFSketch(16).length)
	assert.Equal(8, NewFSketch(100).length)
	assert.Equal(1048576, NewFSketch(10000000).length)
}

func TestFrequnceIsZeroWhenNotExistKey(t *testing.T) {
	assert := assert.New(t)
	key := []byte("123中国")
	sketch := NewFSketch(20)
	assert.Equal(0, sketch.Frequency(key))
}

func TestFrequnceAfterIncrement(t *testing.T) {
	assert := assert.New(t)
	key := []byte("123中国")
	sketch := NewFSketch(20)
	sketch.Increment(key)
	assert.Equal(1, sketch.Frequency(key))
}

func TestMaxFrequnce(t *testing.T) {
	assert := assert.New(t)
	key := []byte("123中国")
	sketch := NewFSketch(20)
	for i := 0; i < 20; i++ {
		sketch.Increment(key)
	}
	assert.Equal(15, sketch.Frequency(key))
}

func TestReset(t *testing.T) {
	assert := assert.New(t)
	key := []byte("123中国")
	sketch := NewFSketch(1)
	n := sketch.threshold * 3 / 2 // <=> sketch.threshold * 1.5
	reset := false
	for i := 0; i < n; i++ {
		// 执行完这个循环，counter>=sketch.threshold 到了
		sketch.Increment(key)
		if sketch.counter != i {
			reset = true
		}

	}
	assert.Equal(true, reset)
	assert.LessOrEqual(sketch.counter, sketch.threshold/2)
}

func TestHeavyHitters(t *testing.T) {
	assert := assert.New(t)
	sketch := NewFSketch(512)

	key := make([]byte, 8)
	for i := 100; i < 100_000; i++ {
		binary.PutVarint(key, int64(i))
		sketch.Increment(key)
	}
	for i := 0; i < 10; i += 2 {
		for j := 0; j < i; j++ {
			binary.PutVarint(key, int64(i))
			sketch.Increment(key)
		}
	}

	// A perfect popularity count yields an array [0, 0, 2, 0, 4, 0, 6, 0, 8, 0]
	popularity := make([]int, 10)
	for i := 0; i < 10; i++ {
		binary.PutVarint(key, int64(i))
		popularity[i] = sketch.Frequency(key)
	}
	for i := 0; i < 10; i++ {
		if (i == 0) || (i == 1) || (i == 3) || (i == 5) || (i == 7) || (i == 9) {
			assert.LessOrEqual(popularity[i], popularity[2])
		} else if i == 2 {
			assert.LessOrEqual(popularity[2], popularity[4])
		} else if i == 4 {
			assert.LessOrEqual(popularity[4], popularity[6])
		} else if i == 6 {
			assert.LessOrEqual(popularity[6], popularity[8])
		}
	}
}

func TestIncrementAt(t *testing.T) {
	assert := assert.New(t)
	sketch := NewFSketch(20)

	sketch.incrementAt(0, 0)
	sketch.incrementAt(0, 4)
	sketch.incrementAt(0, 8)
	sketch.incrementAt(0, 12)
	assert.Equal(uint64(0x0001000100010001), sketch.table[0])

	for i := 0; i < 10; i++ {
		sketch.incrementAt(0, 0)
		sketch.incrementAt(0, 4)
		sketch.incrementAt(0, 8)
		sketch.incrementAt(0, 12)
	}
	assert.Equal(uint64(0x000B000B000B000B), sketch.table[0])
}
