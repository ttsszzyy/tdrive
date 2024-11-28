/*
 * Author: lihy lihy@zhiannet.com
 * Date: 2023-12-13 10:33:24
 * LastEditors: lihy lihy@zhiannet.com
 * Note: Need note condition
 */
package utils

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/syncx"
	"log"
	"testing"
)

func Test_DistinctSlice(t *testing.T) {
	arr := []int64{1, 2, 3, 2, 3, 5, 1, 2, 3, 5, 67, 7, 7, 34, 5, 3423}
	fmt.Println(DistinctSlice[int64](arr))
}

func TestName(t *testing.T) {
	md5 := MD5("21ops.com")
	fmt.Println(md5)
	password := GenPassword(md5, "T.drive")
	fmt.Println(password)
}

func TestGetFileName(t *testing.T) {
	limit := syncx.NewLimit(1)

	if limit.TryBorrow() {
		fmt.Println("Successfully borrowed the first token.")
	} else {
		fmt.Println("Failed to borrow the first token.")
	}

	if limit.TryBorrow() {
		fmt.Println("Successfully borrowed the second token.")
	} else {
		fmt.Println("Failed to borrow the second token.")
	}

	if limit.TryBorrow() {
		fmt.Println("Successfully borrowed the third token.")
	} else {
		fmt.Println("Failed to borrow the third token.")
	}

	err := limit.Return()
	if err != nil {
		log.Println("Error returning token:", err)
	} else {
		fmt.Println("Successfully returned a token.")
	}

	err = limit.Return()
	if err != nil {
		log.Println("Error returning token:", err)
	} else {
		fmt.Println("Successfully returned a token.")
	}
}
