/*
 * @Author: Young
 * @Date: 2022-07-26 11:17:56
 * @LastEditors: Young
 * @LastEditTime: 2022-08-05 20:37:44
 * @FilePath: /buyday/common/utils/rand/rand_test.go
 */

package rand

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"unsafe"
)

func rdstr(rd rand.Source, n int) string {
	var rnd = rand.New(rd)
	b := make([]byte, n)
	l := len(letterBytes) - 1
	rnd.Read(b)
	fmt.Println(b)
	for k, v := range b {
		b[k] = letterBytes[int(v)&l]
	}

	return *(*string)(unsafe.Pointer(&b))
}

func Test_RandString(t *testing.T) {
	randomInt64 := RandomInt64(6, 10)
	fmt.Println(randomInt64)

}

func TestXxx(t *testing.T) {
	count := 1000.0
	for i := 1; i < 1001; i += 3 {
		t.Log(math.Mod(count-float64(i), count))
	}
}
