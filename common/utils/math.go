/*
 * @Author: Young
 * @Date: 2022-08-13 14:48:16
 * @LastEditors: Young
 * @LastEditTime: 2022-08-13 14:48:16
 * @FilePath: /buyday/common/utils/math.go
 */

package utils

// Abs 取int64的绝对值
func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
