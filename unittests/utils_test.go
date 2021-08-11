package unittests

import (
	"github.com/stretchr/testify/assert"
	"goffeine/utils"
	"testing"
)

func TestHashCode(t *testing.T) {
	assert := assert.New(t)
	expected := 1134309195

	assert.Equal(expected, utils.HashCode("abc"), "HashCode不准确")

	expected = 2423445863
	assert.Equal(expected, utils.HashCode("中国"))
}
