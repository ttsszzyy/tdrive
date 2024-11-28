package user

import (
	"T-driver/common/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	storage "github.com/utopiosphe/titan-storage-sdk"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/zeromicro/go-zero/core/logx"
)

var storageCli storage.Storage

func init() {
	titanURL := "https://api-test1.container1.titannet.io"
	apiKey := "YkrCMaGmz+Gu1sBnfKr7UzKnMq2u5Lu5Czq5l0xEuV5vCec09BCGhTCEn+CcF/uQ"
	//secret := "ts-AdvyBwcOZbAvrUPfzcnTSkWtQIQqkIhQRaSeQYBeOyFeBxOO"
	if len(titanURL) == 0 {
		fmt.Println("please set environment variable TITAN_URL, example: export TITAN_URL=Your_titan_url")
		return
	}

	if len(apiKey) == 0 {
		fmt.Println("please set environment variable API_KEY, example: export API_KEY=Your_API_KEY")
		return
	}
	storageCli, _ = storage.Initialize(&storage.Config{
		TitanURL:    titanURL,
		APIKey:      apiKey,
		GroupID:     0,
		UseFastNode: false,
	})
}
func TestName(t *testing.T) {
	titanURL := "https://api-test1.container1.titannet.io"
	apiKey := "YkrCMaGmz+Gu1sBnfKr7UzKnMq2u5Lu5Czq5l0xEuV5vCec09BCGhTCEn+CcF/uQ"
	//secret := "ts-AdvyBwcOZbAvrUPfzcnTSkWtQIQqkIhQRaSeQYBeOyFeBxOO"
	if len(titanURL) == 0 {
		fmt.Println("please set environment variable TITAN_URL, example: export TITAN_URL=Your_titan_url")
		return
	}

	if len(apiKey) == 0 {
		fmt.Println("please set environment variable API_KEY, example: export API_KEY=Your_API_KEY")
		return
	}
	storageCli, err := storage.Initialize(&storage.Config{
		TitanURL:    titanURL,
		APIKey:      apiKey,
		GroupID:     0,
		UseFastNode: false,
	})
	fmt.Println(err)

	filePath := "C:\\Users\\19623\\GolandProjects\\T-driver\\driver项目功能.docx"
	//filePath := "C:\\Users\\19623\\GolandProjects\\T-driver\\a.txt"
	//filePath := "C:\\Users\\19623\\GolandProjects\\T-driver\\u=1453051168,3320119945&fm=253&fmt=auto.webp"
	open, err := os.Open(filePath)
	defer open.Close()
	open.Name()
	myReader := utils.NewMyReader(open, 0)

	//show upload progress
	progress := func(doneSize int64, totalSize int64) {
		fmt.Printf("upload %d, total %d\n", doneSize, totalSize)
		if doneSize == totalSize {
			fmt.Printf("upload success\n")
		}
		if atomic.LoadInt32(&myReader.Paused) == 1 {
			fmt.Println("暂停")
			<-myReader.BlockCh
		}
	}
	go func(my *utils.MyReader) {
		time.Sleep(time.Second * 15)
		fmt.Println("继续")
		myReader.Resume()
	}(myReader)
	go func(my *utils.MyReader) {
		time.Sleep(time.Second * 5)
		fmt.Println("暂停")
		myReader.Pause()
	}(myReader)
	stream, err := storageCli.UploadStream(context.Background(), myReader, open.Name(), progress)
	if err != nil {
		fmt.Println("UploadFile error ", err.Error())
		return
	}
	fmt.Printf("upload %s success\n", stream.String())

	/*fi, err := open.Stat()
	if err != nil {
		fmt.Println("Error stating file:", err)
		return
	}

	fileSize := fi.Size()
	var blockSize int64 = 1024 * 1024
	// 计算需要多少个请求
	numBlocks := int((fileSize + blockSize - 1) / blockSize)
	// 上传每个块
	for blockNum := 0; blockNum < numBlocks; blockNum++ {
		// 计算块的起始位置和长度
		start := int64(blockNum) * blockSize
		length := blockSize
		if start+length > fileSize {
			length = fileSize - start
		}

		// 读取块数据
		block := make([]byte, length)
		_, err := open.ReadAt(block, start)
		if err != nil {
			fmt.Println("Error reading block:", err)
			return
		}
		newReader := bytes.NewReader(block)
		// 发送块的索引和数据
		stream, err := storageCli.UploadStream(context.Background(), myReader, open.Name(), progress)
		if err != nil {
			fmt.Println("UploadFile error ", err.Error())
			return
		}
		fmt.Printf("upload %s success\n", stream.String())
	}*/

	/*shareAssetResult, err := storageCli.GetURL(context.Background(), stream.String())
	urLs := shareAssetResult.URLs
	fmt.Println(urLs)*/

	time.Sleep(5 * time.Second)
}

func Test5(t *testing.T) {
	notifyContext, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(handler1),
	}
	b, err := bot.New("https://26e60eb4cf4e44379d3c81c523759c50.cassini-l1.titannet.io", "7441175021:AAGYI1bXK-j-NTyYh03JRy2349sn8cxSatE", opts...)
	if nil != err {
		logx.Error("链接bot失败！")
		return
	}
	go b.Start(notifyContext)
	http.ListenAndServe(":8443", b.WebhookHandler())
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}

func handler1(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Printf("%d say: %s", update.Message.From.ID, update.Message.Text)
	//图片
	if len(update.Message.Photo) > 0 {
		logx.Info("图片")
		// 获取最新的照片
		photo := update.Message.Photo[len(update.Message.Photo)-1]
		// 下载照片
		b.GetUserProfilePhotos(ctx, &bot.GetUserProfilePhotosParams{UserID: update.Message.From.ID})
		getFile, err := b.GetFile(ctx, &bot.GetFileParams{FileID: photo.FileID})
		if err != nil {
			logx.Error("Error getting file:", err)
		}
		filePath := getFile.FilePath
		if strings.HasPrefix(getFile.FilePath, fmt.Sprintf("/var/lib/telegram-bot-api/%s/", "7441175021:AAGYI1bXK-j-NTyYh03JRy2349sn8cxSatE")) {
			filePath = strings.TrimPrefix(getFile.FilePath, fmt.Sprintf("/var/lib/telegram-bot-api/%s/", "7441175021:AAGYI1bXK-j-NTyYh03JRy2349sn8cxSatE"))
		}
		resp, err := http.Get(fmt.Sprintf("https://26e60eb4cf4e44379d3c81c523759c50.cassini-l1.titannet.io/file/bot%s/%s", "7441175021:AAGYI1bXK-j-NTyYh03JRy2349sn8cxSatE", filePath))
		if err != nil {
			logx.Error("Error downloading file:", err)
		}
		defer resp.Body.Close()
		progress := func(doneSize int64, totalSize int64) {
			fmt.Printf("upload %d, total %d\n", doneSize, totalSize)
			if doneSize == totalSize {
				fmt.Printf("upload success\n")
			}
		}
		stream, err := storageCli.UploadStream(context.Background(), resp.Body, getFile.FilePath, progress)
		if err != nil {
			fmt.Println("UploadFile error ", err.Error())
			return
		}
		fmt.Printf("upload %s success\n", stream.String())
	}

	if update.Message.Audio != nil {
		audio := update.Message.Audio
		logx.Info("语音")
		logx.Info(audio.FileID)
		logx.Info(audio.FileUniqueID)
		logx.Info(audio.FileSize)
		logx.Info(audio.FileName)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "上传成功",
	})
	return
}

func Test8(t *testing.T) {
	// 定义两个int64数字
	dividend := int64(23456)
	divisor := int64(100)

	// 使用操作符获取除法的商和余数
	quotient := dividend / divisor
	remainder := dividend % divisor

	// 打印结果
	fmt.Printf("商: %v\n", quotient)
	fmt.Printf("余数: %v\n", remainder)
	fmt.Printf("总数: %v\n", dividend-dividend%divisor)
	// 原始文件名
	filename := "新建 文本文档 (2)(1).txt"
	newFilename := replaceLastOccurrence(filename, "(1)", "(3)")
	fmt.Println(newFilename)
}

func replaceLastOccurrence(str, old, new string) string {
	// 找到最后一个旧字符串的索引
	pos := strings.LastIndex(str, old)
	if pos == -1 {
		// 如果没有找到旧字符串，直接返回原字符串
		return str
	}
	// 使用新字符串替换旧字符串
	return str[:pos] + new + str[pos+len(old):]
}
