package shm

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	chars = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func randomName() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 6)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), string(b))
}
