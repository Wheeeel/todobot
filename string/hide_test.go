package string

import (
	"testing"
)

func TestMaskASCII(t *testing.T) {
	input := "A"
	mask := "***"
	t.Logf("input=%s, mask=%s, result=%s", input, mask, Hide(input, mask))
	input = "AA"
	t.Logf("input=%s, mask=%s, result=%s", input, mask, Hide(input, mask))
	input = "ABC"
	t.Logf("input=%s, mask=%s, result=%s", input, mask, Hide(input, mask))
	input = "ABCD"
	t.Logf("input=%s, mask=%s, result=%s", input, mask, Hide(input, mask))
	input = "AB C"
	t.Logf("input=%s, mask=%s, result=%s", input, mask, Hide(input, mask))
	input = "家豪 李"
	t.Logf("input=%s, mask=%s, result=%s", input, mask, Hide(input, mask))
}
