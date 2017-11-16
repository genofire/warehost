package lib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ValideDeprecatedPW = "pbkdf2_sha1$10000$a5viM+Paz3o=$orD4shu1Ss+1wPAhAt8hkZ/fH7Y="

func TestPassword(t *testing.T) {
	assert := assert.New(t)

	// wrong count of hashparts
	ok, err := Validate("$$$$$", "wrong-password")
	assert.False(ok)
	assert.Equal(errorNoValideHashFormat, err)

	// hash interations not readable
	ok, err = Validate("a$a$a$a", "wrong-password")
	assert.False(ok)
	assert.Equal(errorNoValideHashFormat, err)

	// wrong password
	ok, err = Validate(ValideDeprecatedPW, "wrong-password")
	assert.False(ok)
	assert.Equal(errorHashDeprecated, err)

	ok, err = Validate(ValideDeprecatedPW, "root")
	assert.True(ok)
	assert.Equal(errorHashDeprecated, err)

	// random internal hash -> no other validation exists
	hash := NewHash("foobar")
	assert.True(strings.Contains(hash, hashfunc))
	ok, err = Validate(hash, "foobar")
	assert.True(ok)
	assert.NoError(err)
}
