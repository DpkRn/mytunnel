package id

import "math/rand"

var charset = []rune("abcdefghijklmnopqrstuvwxyz")

// Generate returns a random subdomain id (6–8 lowercase letters).
func Generate() string {
	length := 6 + rand.Intn(3)
	b := make([]rune, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
