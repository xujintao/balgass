package utils

func Xor(in []byte) []byte {
	var keys = [3]byte{0xFC, 0xCF, 0xAB}
	out := make([]byte, len(in))
	for i := range in {
		out[i] = in[i] ^ keys[i%3]
	}
	return out
}
