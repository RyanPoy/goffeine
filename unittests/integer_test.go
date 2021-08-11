package unittests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"goffeine/utils/integer"
	"testing"
)

func TestNumberOfLeadingZerosForInteger(t *testing.T) {
	assert := assert.New(t)
	numbers, expecteds := readLines("test_files/for_integer/10w_numberOfLeadingZeros.txt")
	length := len(numbers)
	for i := 0; i < length; i++ {
		num, expected := numbers[i], expecteds[i]
		relt := integer.NumberOfLeadingZeros(num)
		assert.Equal(expected, relt, fmt.Sprintf("integer.NumberOfLeadingZeros(%d)=%d, 现在=%d", num, expected, relt))
	}
}

func TestCeilingPowerOfTwoForInteger(t *testing.T) {
	assert := assert.New(t)
	numbers, expecteds := readLines("test_files/for_integer/10w_ceilingPowerOfTwo.txt")
	length := len(numbers)
	for i := 0; i < length; i++ {
		num, expected := numbers[i], expecteds[i]
		relt := integer.CeilingPowerOfTwo(num)
		assert.Equal(expected, relt, fmt.Sprintf("integer.CeilingPowerOfTwo(%d)=%d, 现在=%d", num, expected, relt))
	}
}

