package lib

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltLength              = 8
	hashLength              = 20
	interations             = 10000
	hashfunc                = "sha256"
	errorNoValideHashFormat = Error("no valide hash format")
	errorHashDeprecated     = Error("deprecated hash function used, replace hashed password")
)

var hashlib = map[string]func() hash.Hash{
	"sha1":   sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
}

// Validate a password and a hash
func Validate(hash, password string) (bool, error) {

	parts := strings.Split(hash, "$")
	if len(parts) != 4 {
		return false, errorNoValideHashFormat
	}
	curIter, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, errorNoValideHashFormat
	}
	hashfuncUsed := strings.Split(parts[0], "_")[1]

	dk := pbkdf2.Key([]byte(password), []byte(parts[2]), curIter, len(parts[3])-8, hashlib[hashfuncUsed])
	x := fmt.Sprintf("pbkdf2_%s$%s$%s$%s", hashfuncUsed, parts[1], parts[2], base64.StdEncoding.EncodeToString(dk))

	if hashfuncUsed != hashfunc {
		err = errorHashDeprecated
	} else {
		err = nil
	}
	return (x == hash), err
}

// GenerateRandomString by length for key
func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// NewHash of given password
func NewHash(password string) string {
	salt, _ := generateRandomString(saltLength)
	dk := pbkdf2.Key([]byte(password), []byte(salt), interations, hashLength, hashlib[hashfunc])
	return fmt.Sprintf("pbkdf2_%s$%d$%s$%s", hashfunc, interations, salt, base64.StdEncoding.EncodeToString(dk))
}
