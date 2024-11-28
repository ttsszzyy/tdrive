package validator

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// CheckMobile 检验11位手机号
func CheckMobile(phone string) bool {
	regRuler := "^1[345789]{1}\\d{9}$"
	reg := regexp.MustCompile(regRuler)
	return reg.MatchString(phone)
}

// CheckMail 检验邮箱地址
func CheckMail(mail string) bool {
	regRuler := "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	reg := regexp.MustCompile(regRuler)
	return reg.MatchString(mail)
}

// CheckIdCard 检验身份证
func CheckIdCard(card string) bool {
	regRuler := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"
	// 正则调用规则
	reg := regexp.MustCompile(regRuler)
	// 返回 MatchString 是否匹配
	return reg.MatchString(card)
}

// CheckDomain 判断是否域名合法
func CheckDomain(domain string) bool {
	pattern := `^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(domain)
}

// CheckIp 校验ipv4段
func CheckIPV4(IP string) bool {
	// 字符串这样切割
	strs := strings.Split(IP, ".")
	if len(strs) != 4 {
		return false
	}
	for _, s := range strs {
		if len(s) == 0 || (len(s) > 1 && s[0] == '0') {
			return false
		}
		// 直接访问字符串的值
		if s[0] < '0' || s[0] > '9' {
			return false
		}
		// 字符串转数字
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		if n < 0 || n > 255 {
			return false
		}
	}
	return true
}

// CheckIp 校验ipv6段
func CheckIPV6(IP string) bool {
	strs := strings.Split(IP, ":")
	if len(strs) != 8 {
		return false
	}
	for _, s := range strs {
		if len(s) <= 0 || len(s) > 4 {
			return false
		}
		for i := 0; i < len(s); i++ {
			if s[i] >= '0' && s[i] <= '9' {
				continue
			}
			if s[i] >= 'A' && s[i] <= 'F' {
				continue
			}
			if s[i] >= 'a' && s[i] <= 'f' {
				continue
			}
			return false
		}
	}
	return true
}

// 判断字符串里是否包含中午呢
func CheckContainsZHChar(str string) bool {
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

// 判断端口是否合法
func CheckPort(p int64) bool {
	return p > 0x0000 && p <= 0xFFFF
}
