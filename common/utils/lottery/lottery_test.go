/*
 * Author: 李鸿胤 leeyfann@gmail.com
 * Date: 2023-11-15 14:45:07
 * LastEditors: 李鸿胤 leeyfann@gmail.com
 * Note: Need note condition
 */
package lottery

import (
	"fmt"
	"testing"
)

func TestLottery(t *testing.T) {
	var weights = [][]int{
		{0, 11, 22},
		{133, 11, 23},
		{20, 30, 40},
		{25, 35, 40},
		{1, 1, 1, 1, 1, 1, 22, 0, 0},
	}
	var prizes [][]Prize
	for _, w := range weights {
		var tp []Prize
		for i := range w {
			tp = append(tp, Prize{i + 1, w[i]})
		}
		prizes = append(prizes, tp)
	}
	// prizes := []Prize{
	// 	{ID: 1, Weight: 10},
	// 	{ID: 2, Weight: 20},
	// 	{ID: 3, Weight: 30},
	// 	{ID: 4, Weight: 40},
	// }

	for _, data := range prizes {
		// 进行抽奖
		winner := Draw(data)
		// 输出抽奖结果
		if winner.ID > 0 {
			fmt.Printf("恭喜你中奖了：%v\n", winner.ID)
		} else {
			fmt.Println("很遗憾，未中奖。")
		}
	}

}
