// Written by Gon Yi
package atimer

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randChars(length int) string {
	randNum := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randNum.Intn(62)] // 62 = len(charset)
	}
	return string(b)
}
