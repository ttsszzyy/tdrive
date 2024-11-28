package base64

import "encoding/base64"

// Encode encodes bytes with BASE64 algorithm.
func Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

// Decode decodes bytes with BASE64 algorithm.
func Decode(data []byte) ([]byte, error) {
	src := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(src, data)
	return src[:n], err
}
