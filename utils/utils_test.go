package utils

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestHashCode(t *testing.T) {
	assert := assert.New(t)
	var expected uint32 = 1134309195

	assert.Equal(expected, HashCode("abc"), "HashCode不准确")

	expected = 2423445863
	assert.Equal(expected, HashCode("中国"))
}

func TestNumberOfLeadingZerosForInt(t *testing.T) {
	f, err := os.Open("numberOfLeadingZerosForInt_10w.txt") // 用Caffeine的Java版本生成10w数据
	if err != nil {
		defer f.Close()
	}
	assert := assert.New(t)
	reader := bufio.NewReader(f)
	for ; ; {
		line, ferr := reader.ReadString('\n')
		if ferr != nil && ferr == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		vs := strings.Split(line, ",")

		num, _ := strconv.Atoi(vs[0])
		expected, _ := strconv.Atoi(vs[1])
		relt := NumberOfLeadingZerosForInt(num)
		assert.Equal(expected, relt, fmt.Sprintf("NumberOfLeadingZerosForInt(%d)=%d, 现在=%d", num, expected, relt))
	}
}

func TestCeilingPowerOfTwoForInt(t *testing.T) {
	f, err := os.Open("ceilingPowerOfTwoForInt_10w.txt") // 用Caffeine的Java版本生成10w数据
	if err != nil {
		defer f.Close()
	}
	assert := assert.New(t)
	reader := bufio.NewReader(f)
	for ; ; {
		line, ferr := reader.ReadString('\n')
		if ferr != nil && ferr == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		vs := strings.Split(line, ",")

		num, _ := strconv.Atoi(vs[0])
		expected, _ := strconv.Atoi(vs[1])
		relt := CeilingPowerOfTwoForInt(num)
		assert.Equal(expected, relt, fmt.Sprintf("CeilingPowerOfTwoForInt(%d)=%d, 现在=%d", num, expected, relt))
	}
}


func TestNumberOfLeadingZerosForLong(t *testing.T) {
	f, err := os.Open("numberOfLeadingZerosForLong_10w.txt") // 用Caffeine的Java版本生成10w数据
	if err != nil {
		defer f.Close()
	}
	assert := assert.New(t)
	reader := bufio.NewReader(f)
	for ; ; {
		line, ferr := reader.ReadString('\n')
		if ferr != nil && ferr == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		vs := strings.Split(line, ",")

		num, _ := strconv.Atoi(vs[0])
		expected, _ := strconv.Atoi(vs[1])
		relt := NumberOfLeadingZerosForLong(int64(num))
		assert.Equal(expected, relt, fmt.Sprintf("NumberOfLeadingZerosForLong(%d)=%d, 现在=%d", num, expected, relt))
	}
}

func TestCeilingPowerOfTwoForLong(t *testing.T) {
	f, err := os.Open("ceilingPowerOfTwoForLong_10w.txt") // 用Caffeine的Java版本生成10w数据
	if err != nil {
		defer f.Close()
	}
	assert := assert.New(t)
	reader := bufio.NewReader(f)
	for ; ; {
		line, ferr := reader.ReadString('\n')
		if ferr != nil && ferr == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		vs := strings.Split(line, ",")

		num, _ := strconv.Atoi(vs[0])
		expected, _ := strconv.Atoi(vs[1])
		relt := CeilingPowerOfTwoForLong(int64(num))
		assert.Equal(int64(expected), relt, fmt.Sprintf("CeilingPowerOfTwoForLong(%d)=%d, 现在=%d", num, expected, relt))
	}
}
