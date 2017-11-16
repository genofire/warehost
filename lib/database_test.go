package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	assert := assert.New(t)
	a := JSONNullInt64{}
	a.Scan(300)

	assert.Equal(int64(300), a.Int64)

	a.UnmarshalJSON([]byte{'2', '3'})
	assert.Equal(int64(23), a.Int64)

	err := a.UnmarshalJSON([]byte{})
	assert.Error(err)
	a.UnmarshalJSON([]byte{'n', 'u', 'l', 'l'})
	assert.False(a.Valid)

	value, _ := a.MarshalJSON()
	assert.Equal([]byte{'n', 'u', 'l', 'l'}, value)

	a.UnmarshalJSON([]byte{'1', '4'})
	value, _ = a.MarshalJSON()
	assert.Equal([]byte{'1', '4'}, value)
}
