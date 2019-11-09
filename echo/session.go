package echo

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

func GenerateToken() string {
	uuID := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	return uuID
}
