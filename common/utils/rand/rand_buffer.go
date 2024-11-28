package rand

import (
	"math/rand"
	"time"
)

var (
	bufferChan = make(chan []byte, 10000)
	src        = rand.NewSource(time.Now().UnixNano())
)

func init() {
	go func() {
		var step int
		rnd := rand.New(src)
		for {
			buffer := make([]byte, 1024)
			if n, err := rnd.Read(buffer); err != nil {
				panic(err)
			} else {
				// The random buffer from system is very expensive,
				// so fully reuse the random buffer by changing
				// the step with a different number can
				// improve the performance a lot.
				// for _, step = range []int{4, 5, 6, 7} {
				for _, step = range []int{4} {
					for i := 0; i <= n-4; i += step {
						bufferChan <- buffer[i : i+4]
					}
				}
			}
		}
	}()
}
