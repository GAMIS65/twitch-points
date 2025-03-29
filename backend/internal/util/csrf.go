package util

import (
	"crypto/rand"
	"fmt"
)

func GenerateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
