package string

import (
	"testing"
	"time"
)

func TestGetToken(t *testing.T) {
	println(time.Now().UnixNano())
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789!$*."
	size := 12
	t.Log(GetToken(charset, size))
}
