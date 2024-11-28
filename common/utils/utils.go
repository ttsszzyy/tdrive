package utils

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/mssola/useragent"
)

// EarthRadiusKm 地球半径（单位：公里）
const EarthRadiusKm = 6371.0

var (
	rndVerifyCode = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func HmacEncryptHex(cryptoHash crypto.Hash, key, content string) string {
	hash := hmac.New(cryptoHash.New, []byte(key))
	hash.Write([]byte(content))
	return hex.EncodeToString(hash.Sum(nil))
}

func SHA1(src ...string) string {
	var data []byte
	buf := bytes.NewBuffer(data)
	for _, v := range src {
		buf.WriteString(v)
	}
	h := sha1.New()
	h.Write(buf.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(src string) string {
	m := sha256.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

func GenWxchaPayAuthCode() string {
	return strings.ToUpper(MD5(uuid.NewString()))
}

// MD5 求值MD5码
func MD5(src string) string {
	h := md5.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5ToLower(src string) string {
	h := md5.New()
	h.Write([]byte(src))
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

func GenPassword(pwd string, salt string) string {
	return SHA1(pwd, salt)
}

func GenStoreSn() uint64 {
	source := rand.NewSource(time.Now().UnixNano())
	var max int64 = 999999999
	var min int64 = 100000000
	return uint64(rand.New(source).Int63n(max-min) + min)
}

func CheckPassword(pwd string, src string, salt string) bool {
	check := SHA1(pwd, salt)
	return strings.Compare(check, src) == 0
}

func GenVerifyCode() string {
	return fmt.Sprintf("%06v", rndVerifyCode.Int31n(1000000))
}

func VerifyPhone(phone string) bool {
	rgx := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return rgx.MatchString(phone)
}

func RandomPhone() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return "170" + string(b)
}

func ValidateNameAccount(s string, leng ...int) bool {
	var length int = 20
	if len(leng) > 0 {
		length = leng[0]
	}
	if utf8.RuneCountInString(s) > length {
		return false
	}
	reg := regexp.MustCompile(`^[A-Za-z0-9_]{1,` + strconv.Itoa(length) + `}$`)
	return reg.MatchString(s)
}

func IsSameDay(t1 time.Time, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func CalcAmount(unitPrice float64, count int64) int64 {
	unit := int64(unitPrice * 100.0)
	return unit * count
}

func CutUTF8(src []byte, retain int) []byte {
	var count, idx int
	for i := 0; i < len(src) && idx < retain; idx++ {
		if int(src[i]) > 127 {
			i += 3
			count += 3
		} else {
			i++
			count += 1
		}
	}

	return src[:count]
}

func IntnSliceContains[T int8 | uint8 | int16 | uint16 | int | uint | int32 | uint32 | int64 | uint64](arr []T, elem T) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

// 删除数组中某个元素, true删除全部 false仅删除第一个
func TrimSlice[T comparable](arr []T, t T, f bool) []T {
	var b []T
	for i := 0; i < len(arr); i++ {
		if arr[i] == t {
			if !f {
				if len(b) != i+1 {
					b = append(b, arr[i:]...)
				}
				break
			}
			continue
		}
		b = append(b, arr[i])
	}
	return b
}

func DistinctSlice[T comparable](arr []T) []T {
	var m = make(map[T]bool)
	var b []T
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v] = true
			b = append(b, v)
		}
	}
	return b
}

type HeaderInfo struct {
	OS      string
	Broswer string
	IP      string
}

func GetHeaderFromRequest(r *http.Request) HeaderInfo {
	val := r.Header.Get("User-Agent")
	ua := useragent.New(val)
	// reqCtx := r.Context()
	// ctx := context.WithValue(reqCtx, "os", ua.OS())
	browser, _ := ua.Browser()
	// ctx = context.WithValue(ctx, "broswer", browser)
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	if strings.Contains(ip, ",") {
		ip = strings.Split(ip, ",")[0]
	}
	// fmt.Println("X-Forwarded-For:", r.Header.Get("X-Forwarded-For"))
	// fmt.Println("X-Real-Ip:", r.Header.Get("X-Real-Ip"))
	// fmt.Println("Remote-Addr:", r.RemoteAddr)
	if strings.Contains(ip, ":") {
		ip = ip[0:strings.Index(ip, ":")]
	}
	// fmt.Println("ip is ", ip)
	// ctx = context.WithValue(ctx, "ip", ip)
	// ctx = context.WithValue(ctx, "location", utils.GetIPLocation(ip))
	return HeaderInfo{
		OS:      ua.OS(),
		Broswer: browser,
		IP:      ip,
	}
}

// 生成随机字符串
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i, v := range b {
		b[i] = letters[v%byte(len(letters))]
	}

	return string(b), nil
}

// 生成用户推荐码
func GenerateReferralCode() (string, error) {
	// 获取当前时间戳（纳秒）
	timestamp := time.Now().UnixNano()

	// 生成随机字符串
	randomString, err := generateRandomString(6)
	if err != nil {
		return "", err
	}

	// 将时间戳和随机字符串拼接成推荐码
	referralCode := fmt.Sprintf("%d%s", timestamp, randomString)

	return referralCode, nil
}

// 生成随机字符串
func GenerateRandomFileName() (string, error) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	timestamp := time.Now().UnixNano()
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i, v := range b {
		b[i] = letters[v%byte(len(letters))]
	}
	// 将时间戳和随机字符串拼接成推荐码
	referralCode := fmt.Sprintf("%d%s", timestamp, b)

	return referralCode, nil
}

// 2文件3视频4图片
func IsFileType(filename string) int64 {
	if !strings.Contains(filename, ".") {
		return 1
	}
	// 支持的视频格式
	videoExtensions := []string{".mp4", ".mov", ".avi", ".mkv", ".flv"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, videoExt := range videoExtensions {
		if ext == videoExt {
			return 3
		}
	}
	// 支持的图片格式
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return 4
		}
	}
	return 2
}

type MyReader struct {
	Id       int64
	file     *os.File
	Paused   int32
	pauseCh  chan struct{}
	resumeCh chan struct{}
	BlockCh  chan struct{}
}

func NewMyReader(file *os.File, id int64) *MyReader {
	return &MyReader{
		Id:       id,
		file:     file,
		pauseCh:  make(chan struct{}),
		resumeCh: make(chan struct{}),
	}
}
func (r *MyReader) Read(p []byte) (n int, err error) {
	if atomic.LoadInt32(&r.Paused) == 1 {
		<-r.resumeCh
	}

	select {
	case <-r.pauseCh:
		atomic.StoreInt32(&r.Paused, 1)
	case <-r.resumeCh:
		atomic.StoreInt32(&r.Paused, 0)
		r.BlockCh <- struct{}{}
	default:
	}

	return r.file.Read(p)
}

func (r *MyReader) Pause() {
	r.pauseCh <- struct{}{}
}

func (r *MyReader) Resume() {
	r.resumeCh <- struct{}{}
}

func ReplaceLastOccurrence(str, old, new string) string {
	// 找到最后一个旧字符串的索引
	pos := strings.LastIndex(str, old)
	if pos == -1 {
		// 如果没有找到旧字符串，直接返回原字符串
		return str
	}
	// 使用新字符串替换旧字符串
	return str[:pos] + new + str[pos+len(old):]
}

// GetFileName 函数返回给定路径的文件名部分，不包含扩展名。
// 主要用于从完整文件路径中提取文件名，用于需要区分文件扩展名的情况。
// 参数:
//
//	filename: 完整文件路径的字符串。
//
// 返回值:
//
//	包含扩展名和后缀。
func GetFileName(filename string) (string, string) {
	if filename == "" {
		return filename, ""
	}
	// 获取文件的扩展名
	ext := filepath.Ext(filename)
	// 去掉文件名中的扩展名
	filenameWithoutExtension := filename[:len(filename)-len(ext)]
	return filenameWithoutExtension, ext
}

// toRadians 将角度转换为弧度
func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

// HaversineDistance 使用哈弗赛公式计算两个坐标之间的距离
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// 将纬度和经度转换为弧度
	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	// 将起点和终点的纬度也转换为弧度
	lat1 = toRadians(lat1)
	lat2 = toRadians(lat2)

	// 哈弗赛公式
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 计算距离
	distance := EarthRadiusKm * c
	return distance
}
