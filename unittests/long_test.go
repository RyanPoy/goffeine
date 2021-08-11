package unittests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"goffeine/utils/long"
	"testing"
)

func TestNumberOfLeadingZerosForLong(t *testing.T) {
	assert := assert.New(t)
	numbers, expecteds := readLines("test_files/for_long/10w_numberOfLeadingZeros.txt")
	length := len(numbers)
	for i := 0; i < length; i++ {
		num, expected := numbers[i], expecteds[i]
		relt := long.NumberOfLeadingZeros(int64(num))
		assert.Equal(expected, relt, fmt.Sprintf("long.NumberOfLeadingZeros(%d)=%d, 现在=%d", num, expected, relt))
	}
}

func TestCeilingPowerOfTwoForLong(t *testing.T) {
	assert := assert.New(t)
	numbers, expecteds := readLines("test_files/for_long/10w_ceilingPowerOfTwo.txt")
	length := len(numbers)
	for i := 0; i < length; i++ {
		num, expected := numbers[i], expecteds[i]
		relt := long.CeilingPowerOfTwo(int64(num))
		assert.Equal(int64(expected), relt, fmt.Sprintf("long.CeilingPowerOfTwo(%d)=%d, 现在=%d", num, expected, relt))
	}
}
