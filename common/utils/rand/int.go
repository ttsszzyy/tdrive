/*
 * Author: lihy lihy@zhiannet.com
 * Date: 2023-12-08 10:37:05
 * LastEditors: lihy lihy@zhiannet.com
 * Note: Need note condition
 */
package rand

import (
	"math/rand"
	"sync"
	"time"
)

// 生成随机种子
var source = rand.NewSource(time.Now().UnixNano())
var locker = &sync.Mutex{}

// 随机获取一个Int数字
func Int(min int, max int) int {
	if min > max {
		min, max = max, min
	}
	r := max - min + 1
	if r == 0 {
		return min
	}

	locker.Lock()
	result := min + int(source.Int63()%int64(r))
	locker.Unlock()
	return result
}
func RandomInt64(min int64, max int64) int64 {
	if min > max {
		min, max = max, min
	}
	r := max - min + 1
	if r == 0 {
		return min
	}

	locker.Lock()
	result := min + source.Int63()%r
	locker.Unlock()
	return result
}

// 随机获取一个Int64数字
func Int64() int64 {
	locker.Lock()
	result := source.Int63()
	locker.Unlock()
	return result
}
