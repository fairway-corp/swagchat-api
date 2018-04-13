package utils

import (
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// AppendStrings is append strings
func AppendStrings(strings ...string) string {
	buf := make([]byte, 0)
	for _, str := range strings {
		buf = append(buf, str...)
	}
	return string(buf)
}

// IsValidID is valid ID
func IsValidID(ID string) bool {
	r := regexp.MustCompile(`(?m)^[0-9a-zA-Z-]+$`)
	return r.MatchString(ID)
}

var token68Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~+/")

// GenerateClientSecret is generate clientSecret
func GenerateClientSecret(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = token68Letters[rand.Intn(len(token68Letters))]
	}
	return string(b)
}

// GenerateUUID is generate UUID
func GenerateUUID() string {
	uuid := uuid.NewV4().String()
	return uuid
}

// GenerateClientID is generate clientID
func GenerateClientID() string {
	uuid := uuid.NewV4().String()
	return strings.Replace(uuid, "-", "", -1)
}

// GetFileNameWithoutExt is get filename without extention
func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
