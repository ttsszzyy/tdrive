package rand

import (
	"math/rand"
	"time"
	"unsafe"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	seed = rand.NewSource(time.Now().UnixNano())
)

func Seed(s int64) {
	seed = rand.NewSource(s)
}

func B(n int) []byte {
	if n <= 0 {
		return nil
	}
	i := 0
	b := make([]byte, n)
	for {
		copy(b[i:], <-bufferChan)
		i += 4
		if i >= n {
			break
		}
	}
	return b
}

// RandString RandStringBytesMaskImprSrcUnsafe
func RandStringEx(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, seed.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = seed.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func RandKey(src string, n int) (string, []int) {
	var rnd = rand.New(seed)
	b := make([]byte, n)
	idx := make([]int, n)
	l := len(src) - 1
	rnd.Read(b)
	for k, v := range b {
		idx[k] = int(v) & l
		b[k] = src[idx[k]]
	}

	return *(*string)(unsafe.Pointer(&b)), idx
}

func RandString(n int) string {
	var rnd = rand.New(seed)
	b := make([]byte, n)
	l := len(letterBytes) - 1
	rnd.Read(b)
	for k, v := range b {
		b[k] = letterBytes[int(v)&l]
	}

	return *(*string)(unsafe.Pointer(&b))
}
