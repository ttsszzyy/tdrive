/*
 * Author: 李鸿胤 leeyfann@gmail.com
 * Date: 2023-11-14 21:52:55
 * LastEditors: 李鸿胤 leeyfann@gmail.com
 * Note: Need note condition
 */
package locker

import (
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func BenchmarkXxx(b *testing.B) {
	c := redis.RedisConf{
		Host: "localhost:6379",
	}

	key1 := "tst1"
	re := c.NewRedis()
	b.ResetTimer()

	timeout := 10 * time.Second

	// var concurrency atomic.Int64

	// for i := 0; i < b.N; i++ {
	// 	go func() {
	// 		lk := NewRedisLocker(re, key1)
	// 		lk.Spin(timeout)
	// 		lk.Unlock()
	// 	}()

	// concurrency.Add(1)
	// }
	// b.Log(concurrency)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			lk := NewRedisLocker(re, key1)
			lk.Spin(timeout)
			lk.Unlock()
		}
	})
}
