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

func TestNumberOfLeadingZeros(t *testing.T) {
	f, err := os.Open("numberOfLeadingZeros_10w.txt") // 用Caffeine的Java版本生成10w数据
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
		relt := NumberOfLeadingZeros(num)
		assert.Equal(expected, relt, fmt.Sprintf("NumberOfLeadingZeros(%d)=%d, 现在=%d", num, expected, relt))
	}
}

func TestCeilingPowerOfTwo(t *testing.T) {
	f, err := os.Open("ceilingPowerOfTwo_10w.txt") // 用Caffeine的Java版本生成10w数据
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
		relt := CeilingPowerOfTwo(num)
		assert.Equal(expected, relt, fmt.Sprintf("CeilingPowerOfTwo(%d)=%d, 现在=%d", num, expected, relt))
	}
}


//func (t *testing.T) {
//	assert := assert.New(t)
//	assert.Equal(32, CeilingPowerOfTwo(30))
//	assert.Equal(128, CeilingPowerOfTwo(90))
//	//for i := 1; i < math.MaxInt32; i++ {
//	//	fmt.Println(i, CeilingPowerOfTwo(i))
//	//}
//}
