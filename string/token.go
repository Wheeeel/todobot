package string

import (
	"math/rand"
	"time"
)

func GetToken(chars string, length int) (randstr string) {
	randstr = ""
	rand.Seed(time.Now().UnixNano())
	charlen := len(chars)
	for i := 0; i < length; i++ {
		x := rand.Intn(charlen)
		randstr += string(chars[x])
	}
	return randstr
}
