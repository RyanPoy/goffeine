package goffeine_test

import (
	"github.com/stretchr/testify/assert"
	"goffeine"
	"goffeine/internal/utils"
	"math"
	"strconv"
	"testing"
)

func TestEnsureCapacity_smaller(t *testing.T) {
	sketch := goffeine.NewSketch(512)
	size := len(sketch.Table)
	sketch.EnsureCapacity(size / 2)
	assert.Equal(t, size, len(sketch.Table))
	assert.Equal(t, 10*size, sketch.SampleSize)
	assert.Equal(t, (size>>3)-1, sketch.BlockMask)
}

func TestEnsureCapacity_larger(t *testing.T) {
	sketch := goffeine.NewSketch(512)
	size := len(sketch.Table)
	sketch.EnsureCapacity(2 * size)
	assert.Equal(t, 2*size, len(sketch.Table))
	assert.Equal(t, 10*2*size, sketch.SampleSize)
	assert.Equal(t, ((2*size)>>3)-1, sketch.BlockMask)
}

func TestEnsureCapacity_maximum(t *testing.T) {
	sketch := goffeine.NewSketch(512)
	size := math.MaxInt32/10 + 1
	sketch = goffeine.NewSketch(size)

	assert.Equal(t, math.MaxInt32, sketch.SampleSize)
	assert.Equal(t, utils.CeilingPowerOfTwo32(size), len(sketch.Table))
	assert.Equal(t, (len(sketch.Table)>>3)-1, sketch.BlockMask)
}

func TestIncrement_once(t *testing.T) {
	sketch := goffeine.NewSketch(512)
	item := "key1"
	sketch.Increment(item)
	assert.Equal(t, 1, sketch.Frequency(item))
}

func TestIncrement_max(t *testing.T) {
	sketch := goffeine.NewSketch(512)
	item := "key1"
	for i := 0; i < 20; i++ {
		sketch.Increment(item)
	}
	assert.Equal(t, 15, sketch.Frequency(item))
}

func TestIncrement_distinct(t *testing.T) {
	sketch := goffeine.NewSketch(512)
	sketch.Increment("key1")
	sketch.Increment("key1_1")
	assert.Equal(t, 1, sketch.Frequency("key1"))
	assert.Equal(t, 1, sketch.Frequency("key1_1"))
	assert.Equal(t, 0, sketch.Frequency("key1_2"))
}

func TestIncrement_zero(t *testing.T) {
	sketch := goffeine.NewSketch(512)
	sketch.Increment("")
	assert.Equal(t, 1, sketch.Frequency(""))
}

func TestReset(t *testing.T) {
	reset := false
	sketch := goffeine.NewSketch(64)
	sketch.EnsureCapacity(64)

	for i := 1; i < 20*len(sketch.Table); i++ {
		sketch.Increment(strconv.Itoa(i))
		if sketch.Size != i {
			reset = true
			break
		}
	}
	assert.True(t, reset)
	assert.Greater(t, sketch.SampleSize/2, sketch.Size)
}

func TestFull(t *testing.T) {
	sketch := goffeine.NewSketch(512)
	sketch.SampleSize = math.MaxInt32

	for i := 0; i < 100_000; i++ {
		sketch.Increment(strconv.Itoa(i))
	}
	//for item := range sketch.Table {
	//	assert.Equal(t, 64, goffeine.BitCount64(int64(item)))
	//}
	sketch.Reset()
	for item := range sketch.Table {
		assert.Equal(t, goffeine.ResetMask, int64(item))
	}
}

//  @Test
//  public void heavyHitters() {
//    FrequencySketch<Double> sketch = makeSketch(512);
//    for (int i = 100; i < 100_000; i++) {
//      sketch.increment((double) i);
//    }
//    for (int i = 0; i < 10; i += 2) {
//      for (int j = 0; j < i; j++) {
//        sketch.increment((double) i);
//      }
//    }
//
//    // A perfect popularity count yields an array [0, 0, 2, 0, 4, 0, 6, 0, 8, 0]
//    int[] popularity = new int[10];
//    for (int i = 0; i < 10; i++) {
//      popularity[i] = sketch.frequency((double) i);
//    }
//    for (int i = 0; i < popularity.length; i++) {
//      if ((i == 0) || (i == 1) || (i == 3) || (i == 5) || (i == 7) || (i == 9)) {
//        assertThat(popularity[i]).isAtMost(popularity[2]);
//      } else if (i == 2) {
//        assertThat(popularity[2]).isAtMost(popularity[4]);
//      } else if (i == 4) {
//        assertThat(popularity[4]).isAtMost(popularity[6]);
//      } else if (i == 6) {
//        assertThat(popularity[6]).isAtMost(popularity[8]);
//      }
//    }
//  }
//
//  @DataProvider(name = "sketch")
//  public Object[][] providesSketch() {
//    return new Object[][] {{ makeSketch(512) }};
//  }
//
//  private static <E> FrequencySketch<E> makeSketch(long maximumSize) {
//    var sketch = new FrequencySketch<E>();
//    sketch.ensureCapacity(maximumSize);
//    return sketch;
//  }
//}
//
