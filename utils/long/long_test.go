package long

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

func TestCeilingPowerOfTwoForLong(t *testing.T) {
	assert := assert.New(t)
	numbers, expecteds := readLines("10w_ceilingPowerOfTwo.txt")
	length := len(numbers)
	for i := 0; i < length; i++ {
		num, expected := numbers[i], expecteds[i]
		relt := CeilingPowerOfTwo(int64(num))
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