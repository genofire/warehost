package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert := assert.New(t)
	err := Error("hallo")

	assert.Equal("hallo", err.Error())

}
