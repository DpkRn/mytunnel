package pkg

import (
	"math/rand"

	"github.com/google/uuid"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyz")

func GenerateID() string {
	length := 6 + rand.Intn(3) // 6–8 chars

	b := make([]rune, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateRequestID() string {
	return "req_" + uuid.New().String()
}
