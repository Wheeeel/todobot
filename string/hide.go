package string

func Hide(input string, mask string) string {
	mlen := len([]rune(mask))
	ilen := len([]rune(input))
	rinput := []rune(input)
	rmask := []rune(mask)
	if ilen < 2 {
		return input
	}
	for i := 0; i < (ilen-2)/2; i++ {
		rinput[i*2+1] = rune(rmask[i%mlen])
	}
	return string(rinput)
}
