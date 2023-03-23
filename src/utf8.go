package runtime

const (
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
)

const runeError = '\uFFFD'

func encoderune(p []byte, r rune) int {
	return 0
}

// builtin
func decoderune(s string, k int) (r rune, pos int) {
	return
}
