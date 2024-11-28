package lottery

import (
	"errors"
	"math/rand"
	"sort"
	"time"
)

type Prize struct {
	ID     int
	Weight int
}

var (
	ErrRangeOut = errors.New("total range can not over 100")
)

func Draw(prizes []Prize) Prize {
	rand.Seed(time.Now().UnixNano()) // 使用当前时间作为随机种子

	// 计算总权重
	totalWeight := 0
	var realPrizes []Prize
	for _, p := range prizes {
		if p.Weight > 0 {
			realPrizes = append(realPrizes, p)
			totalWeight += p.Weight
		}
	}

	if len(realPrizes) == 0 {
		return Prize{}
	}

	sort.Slice(realPrizes, func(i, j int) bool { return realPrizes[i].Weight < realPrizes[j].Weight })

	// 生成随机数
	randomNum := rand.Intn(totalWeight)

	curr := 0
	for _, v := range realPrizes {
		curr += v.Weight
		if randomNum < curr {
			return v
		}
	}
	return Prize{}

}
