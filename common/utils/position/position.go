/*
 * Author: 李鸿胤 leeyfann@gmail.com
 * Date: 2023-11-20 16:59:53
 * LastEditors: 李鸿胤 leeyfann@gmail.com
 * Note: Need note condition
 */
package position

import (
	"fmt"
	"math"
	"strconv"
)

// 地球半径，单位米
const R = 6367000

// lonA, latA分别为A点的纬度和经度
// lonB, latB分别为B点的纬度和经度
// 返回的距离单位为米/千米
func Sphere(lonA, latA, lonB, latB string) string {
	lna, _ := strconv.ParseFloat(lonA, 64)
	laa, _ := strconv.ParseFloat(latA, 64)
	lnb, _ := strconv.ParseFloat(lonB, 64)
	lab, _ := strconv.ParseFloat(latB, 64)

	c := math.Sin(laa)*math.Sin(lab)*math.Cos(lna-lnb) + math.Cos(laa)*math.Cos(lab)
	res := R * math.Acos(c) * math.Pi / 180
	if res > 1000 {
		return fmt.Sprintf("%.1fkm", res/1000)
	} else {
		return fmt.Sprintf("%.1fm", res)
	}
}
