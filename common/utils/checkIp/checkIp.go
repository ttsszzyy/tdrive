package checkIp

import (
	"fmt"
	"strconv"
	"strings"
)

type CheckIpRes struct {
	ThirdMinIp  int    `json:"thirdMinIp"`  //第三段ip 最大值
	ThirdMxaIp  int    `json:"thirdMxaIp"`  //第三段 最大值
	FourthMinIp int    `json:"fourthMinIp"` //第四段ip 最小值
	FourthMaxIp int    `json:"fourthMaxIp"` //第四段 最大值
	ThirdIp     int    `json:"thirdIp"`     //原ip 第三段
	FourthIp    int    `json:"fourthIp"`    //原ip 第四段
	MaxIp       string `json:"MaxIp"`
	MinIp       string `json:"MinIp"`
}

func GetCidrIpRange(cidr string) (res *CheckIpRes) {
	ip := strings.Split(cidr, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	ThirdMinIp, ThirdMxaIp := GetIpSeg3Range(ipSegs, maskLen)
	FourthMinIp, FourthMaxIp := GetIpSeg4Range(ipSegs, maskLen)
	res = &CheckIpRes{}
	res.ThirdMinIp = ThirdMinIp
	res.ThirdMxaIp = ThirdMxaIp
	res.FourthMinIp = FourthMinIp
	res.FourthMaxIp = FourthMaxIp
	ipSeg3, _ := strconv.Atoi(ipSegs[2])
	ipSeg4, _ := strconv.Atoi(ipSegs[3])
	res.ThirdIp = ipSeg3
	res.FourthIp = ipSeg4
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."
	res.MinIp = ipPrefix + strconv.Itoa(ThirdMinIp) + "." + strconv.Itoa(FourthMinIp+1)
	res.MaxIp = ipPrefix + strconv.Itoa(ThirdMxaIp) + "." + strconv.Itoa(FourthMaxIp-5)
	return res
}

//获取Cidr的掩码
func GetCidrIpMask(maskLen int) string {
	// ^uint32(0)二进制为32个比特1，通过向左位移，得到CIDR掩码的二进制
	cidrMask := ^uint32(0) << uint(32-maskLen)
	//计算CIDR掩码的四个片段，将想要得到的片段移动到内存最低8位后，将其强转为8位整型，从而得到
	cidrMaskSeg1 := uint8(cidrMask >> 24)
	cidrMaskSeg2 := uint8(cidrMask >> 16)
	cidrMaskSeg3 := uint8(cidrMask >> 8)
	cidrMaskSeg4 := uint8(cidrMask & uint32(255))
	return fmt.Sprint(cidrMaskSeg1) + "." + fmt.Sprint(cidrMaskSeg2) + "." + fmt.Sprint(cidrMaskSeg3) + "." + fmt.Sprint(cidrMaskSeg4)
}

//得到第三段IP的区间（第一片段.第二片段.第三片段.第四片段）
func GetIpSeg3Range(ipSegs []string, maskLen int) (int, int) {
	if maskLen > 24 {
		segIp, _ := strconv.Atoi(ipSegs[2])
		return segIp, segIp
	}
	ipSeg, _ := strconv.Atoi(ipSegs[2])
	return GetIpSegRange(uint8(ipSeg), uint8(24-maskLen))
}

//得到第四段IP的区间（第一片段.第二片段.第三片段.第四片段）
func GetIpSeg4Range(ipSegs []string, maskLen int) (int, int) {
	ipSeg, _ := strconv.Atoi(ipSegs[3])
	segMinIp, segMaxIp := GetIpSegRange(uint8(ipSeg), uint8(32-maskLen))
	return segMinIp + 1, segMaxIp
}

//根据用户输入的基础IP地址和CIDR掩码计算一个IP片段的区间
func GetIpSegRange(userSegIp, offset uint8) (int, int) {
	var ipSegMax uint8 = 255
	netSegIp := ipSegMax << offset
	segMinIp := netSegIp & userSegIp
	segMaxIp := userSegIp&(255<<offset) | ^(255 << offset)
	return int(segMinIp), int(segMaxIp)
}
