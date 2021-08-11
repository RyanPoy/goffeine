package utils

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"goffeine/utils/integer"
	"goffeine/utils/long"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestHashCode(t *testing.T) {
	assert := assert.New(t)
	expected := 1134309195

	assert.Equal(expected, HashCode("abc"), "HashCode不准确")

	expected = 2423445863
	assert.Equal(expected, HashCode("中国"))
}

func TestNumberOfLeadingZerosForInt(t *testing.T) {
	assert := assert.New(t)
	numbers, expecteds := readLines("test_files/numberOfLeadingZerosForInt_10w.txt")
	length := len(numbers)
	for i := 0; i < length; i++ {
		num, expected := numbers[i], expecteds[i]
		relt := integer.NumberOfLeadingZeros(num)
		assert.Equal(expected, relt, fmt.Sprintf("integer.NumberOfLeadingZeros(%d)=%d, 现在=%d", num, expected, relt))
	}
}

func TestCeilingPowerOfTwoForInt(t *testing.T) {
	assert := assert.New(t)
	numbers, expecteds := readLines("test_files/ceilingPowerOfTwoForInt_10w.txt")
	length := len(numbers)
	for i := 0; i < length; i++ {
		num, expected := numbers[i], expecteds[i]
		relt := integer.CeilingPowerOfTwo(num)
		assert.Equal(expected, relt, fmt.Sprintf("integer.CeilingPowerOfTwo(%d)=%d, 现在=%d", num, expected, relt))
	}
}

func TestNumberOfLeadingZerosForLong(t *testing.T) {
	assert := assert.New(t)
	numbers, expecteds := readLines("test_files/numberOfLeadingZerosForLong_10w.txt")
	length := len(numbers)
	for i := 0; i < length; i++ {
		num, expected := numbers[i], expecteds[i]
		relt := long.NumberOfLeadingZeros(int64(num))
		assert.Equal(expected, relt, fmt.Sprintf("long.NumberOfLeadingZeros(%d)=%d, 现在=%d", num, expected, relt))
	}
}

func TestCeilingPowerOfTwoForLong(t *testing.T) {
	assert := assert.New(t)
	numbers, expecteds := readLines("test_files/ceilingPowerOfTwoForLong_10w.txt")
	length := len(numbers)
	for i := 0; i < length; i++ {
		num, expected := numbers[i], expecteds[i]
		relt := long.CeilingPowerOfTwo(int64(num))
		assert.Equal(int64(expected), relt, fmt.Sprintf("long.CeilingPowerOfTwo(%d)=%d, 现在=%d", num, expected, relt))
	}
}

func readLines(fname string) ([]int, []int) {
	f, err := os.Open(fname) // 用Caffeine的Java版本生成10w数据
	if err != nil {
		defer f.Close()
	}
	reader := bufio.NewReader(f)
	numbers, expecteds := make([]int, 0), make([]int, 0)

	for ; ; {
		line, ferr := reader.ReadString('\n')
		if ferr != nil && ferr == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		vs := strings.Split(line, ",")

		num, _ := strconv.Atoi(vs[0])
		expected, _ := strconv.Atoi(vs[1])

		numbers = append(numbers, num)
		expecteds = append(expecteds, expected)
	}
	return numbers, expecteds
}
