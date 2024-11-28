/*
 * Author: 李鸿胤 leeyfann@gmail.com
 * Date: 2023-11-20 17:06:58
 * LastEditors: 李鸿胤 leeyfann@gmail.com
 * Note: Need note condition
 */
package position

import (
	"fmt"
	"strings"
	"testing"
)

func TestCalDistance(t *testing.T) {
	strArr := []string{
		"118.80507,31.97087,118.80743,31.97194",
		"118.80776,31.97208,118.80743,31.97194",
		"138.81064,31.97207,118.80743,31.97194",
	}

	for _, v := range strArr {
		sarr := strings.Split(v, ",")

		fmt.Println(Sphere(sarr[0], sarr[1], sarr[2], sarr[3]))
	}
}
