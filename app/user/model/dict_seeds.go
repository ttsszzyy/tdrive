/*
 * Author: lihy lihy@zhiannet.com
 * Date: 2024-01-05 16:41:24
 * LastEditors: lihy lihy@zhiannet.com
 * Note: Need note condition
 */
package model

import "time"

var FixedDicts = []Dict{
	RegisterYearDict,
	RegisterTGDict,
	RegisterRewardDict,
	ActivityLevelDict,
	RecommendDict,
	RecommendVipDict,
	RecommendOneDict,
	RecommendThreeDict,
	RecommendTenDict,
	RecommendTwentyFiveDict,
	RecommendFiftyDict,
	RecommendHundredDict,
	RecommendFiveHundredDict,
	RecommendThousandDict,
	RecommendTenThousandDict,
	RecommendOneHundredThousandDict,
	ActivityRuleDict,
	SignDict,
	SpaceMerchantExchangeDict,
}

// 注册奖励数据
var RegisterYearDictCode = "RR2024073001"
var RegisterYearDict = Dict{
	Code:        RegisterYearDictCode,
	ParamType:   1,
	Name:        "註冊年限",
	Desc:        "注册奖励",
	Fixed:       int64(Yes),
	Value:       "50000",
	CreatedTime: time.Now().Unix(),
}
var RegisterTGDictCode = "RR2024073002"
var RegisterTGDict = Dict{
	Code:        RegisterTGDictCode,
	ParamType:   1,
	Name:        "TG會員",
	Desc:        "注册奖励",
	Fixed:       int64(Yes),
	Value:       "100000",
	CreatedTime: time.Now().Unix(),
}
var RegisterRewardDictCode = "RR2024073003"
var RegisterRewardDict = Dict{
	Code:        RegisterRewardDictCode,
	ParamType:   1,
	Name:        "幸運獎勵",
	Desc:        "注册奖励",
	Fixed:       int64(Yes),
	Value:       "90000", // 需要加上基础的10000
	CreatedTime: time.Now().Unix(),
}
var ActivityLevelDictCode = "RR2024073004"
var ActivityLevelDict = Dict{
	Code:        ActivityLevelDictCode,
	ParamType:   1,
	Name:        "活動等級",
	Desc:        "注册奖励",
	Fixed:       int64(Yes),
	Value:       "100",
	CreatedTime: time.Now().Unix(),
}

// 引荐奖励
var RecommendDictCode = "RR2024073101"
var RecommendDict = Dict{
	Code:        RecommendDictCode,
	ParamType:   2,
	Name:        "引薦雙方獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是随机积分",
	Fixed:       int64(Yes),
	Value:       "100000",
	BackupValue: 50000,
	CreatedTime: time.Now().Unix(),
}

// 邀请tg会员奖励
var RecommendVipDictCode = "RR2024073112"
var RecommendVipDict = Dict{
	Code:        RecommendVipDictCode,
	ParamType:   2,
	Name:        "引薦tg会员獎勵",
	Desc:        "引荐奖励 参数1是积分参数",
	Fixed:       int64(Yes),
	Value:       "500000",
	CreatedTime: time.Now().Unix(),
}

// 引荐1人奖励
var RecommendOneDictCode = "RR2024073102"
var RecommendOneDict = Dict{
	Code:        RecommendOneDictCode,
	ParamType:   5,
	Name:        "引薦1人獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是人数",
	Fixed:       int64(Yes),
	Value:       "50000",
	BackupValue: 1,
	CreatedTime: time.Now().Unix(),
}

// 引荐3人奖励
var RecommendThreeDictCode = "RR2024073103"
var RecommendThreeDict = Dict{
	Code:        RecommendThreeDictCode,
	ParamType:   5,
	Name:        "引薦3人獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是人数",
	Fixed:       int64(Yes),
	Value:       "500000",
	BackupValue: 3,
	CreatedTime: time.Now().Unix(),
}

// 引荐10人奖励
var RecommendTenDictCode = "RR2024073104"
var RecommendTenDict = Dict{
	Code:        RecommendTenDictCode,
	ParamType:   5,
	Name:        "引薦10人獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是人数",
	Fixed:       int64(Yes),
	Value:       "2000000",
	BackupValue: 10,
	CreatedTime: time.Now().Unix(),
}

// 引荐25人奖励
var RecommendTwentyFiveDictCode = "RR2024073105"
var RecommendTwentyFiveDict = Dict{
	Code:        RecommendTwentyFiveDictCode,
	ParamType:   5,
	Name:        "引薦25人獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是人数",
	Fixed:       int64(Yes),
	Value:       "2500000",
	BackupValue: 25,
	CreatedTime: time.Now().Unix(),
}

// 引荐50人奖励
var RecommendFiftyDictCode = "RR2024073106"
var RecommendFiftyDict = Dict{
	Code:        RecommendFiftyDictCode,
	ParamType:   5,
	Name:        "引薦50人獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是人数",
	Fixed:       int64(Yes),
	Value:       "3000000",
	BackupValue: 50,
	CreatedTime: time.Now().Unix(),
}

// 引荐100人奖励
var RecommendHundredDictCode = "RR2024073107"
var RecommendHundredDict = Dict{
	Code:        RecommendHundredDictCode,
	ParamType:   5,
	Name:        "引薦100人獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是人数",
	Fixed:       int64(Yes),
	Value:       "5000000",
	BackupValue: 100,
	CreatedTime: time.Now().Unix(),
}

// 引荐500人奖励
var RecommendFiveHundredDictCode = "RR2024073108"
var RecommendFiveHundredDict = Dict{
	Code:        RecommendFiveHundredDictCode,
	ParamType:   5,
	Name:        "引薦500人獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是人数",
	Fixed:       int64(Yes),
	Value:       "15000000",
	BackupValue: 500,
	CreatedTime: time.Now().Unix(),
}

// 引荐1000人奖励
var RecommendThousandDictCode = "RR2024073109"
var RecommendThousandDict = Dict{
	Code:        RecommendThousandDictCode,
	ParamType:   5,
	Name:        "引薦1000人獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是人数",
	Fixed:       int64(Yes),
	Value:       "20000000",
	BackupValue: 1000,
	CreatedTime: time.Now().Unix(),
}

// 引荐10000人奖励
var RecommendTenThousandDictCode = "RR2024073110"
var RecommendTenThousandDict = Dict{
	Code:        RecommendTenThousandDictCode,
	ParamType:   5,
	Name:        "引薦10000人獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是人数",
	Fixed:       int64(Yes),
	Value:       "1000000000",
	BackupValue: 10000,
	CreatedTime: time.Now().Unix(),
}

// 引荐100000人奖励
var RecommendOneHundredThousandDictCode = "RR2024073111"
var RecommendOneHundredThousandDict = Dict{
	Code:        RecommendOneHundredThousandDictCode,
	ParamType:   5,
	Name:        "引薦100000人獎勵",
	Desc:        "引荐奖励 参数1是积分参数2是人数",
	Fixed:       int64(Yes),
	Value:       "10000000000",
	BackupValue: 100000,
	CreatedTime: time.Now().Unix(),
}

// 活动规则
var ActivityRuleDictCode = "RR2024073113"
var ActivityRuleDict = Dict{
	Code:        ActivityRuleDictCode,
	ParamType:   3,
	Name:        "活動規則",
	Desc:        "引荐规划对应文字",
	Fixed:       int64(Yes),
	Value:       "此次贈送空間為測試去中心化空間，更新正式版本會進行刪檔處理，請勿存放重要文件！",
	CreatedTime: time.Now().Unix(),
}

var SignDictCode = "RR2024080105"

// 签到奖励
var SignDict = Dict{
	Code:        SignDictCode,
	ParamType:   4,
	Name:        "簽到獎勵",
	Desc:        "其他设置",
	Fixed:       int64(Yes),
	Value:       "100",
	CreatedTime: time.Now().Unix(),
}
var SpaceMerchantExchangeDictCode = "RR2024080106"

// 空间商人兑换
var SpaceMerchantExchangeDict = Dict{
	Code:        SpaceMerchantExchangeDictCode,
	ParamType:   4,
	Name:        "空間商人兌換",
	Desc:        "空间商人兑换 参数1是MB 参数2是积分",
	Fixed:       int64(Yes),
	Value:       "1",
	BackupValue: 1000,
	CreatedTime: time.Now().Unix(),
}
