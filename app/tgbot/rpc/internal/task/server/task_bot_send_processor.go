package server

import (
	"T-driver/app/tgbot/rpc/internal/svc"
	"T-driver/app/tgbot/rpc/types"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/utils"
	"T-driver/common/utils/base64"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type TaskBotSendProcessor struct {
	svcCtx  *svc.ServiceContext
	ctx     context.Context
	limiter *limit.PeriodLimit
}

func NewTaskBotSendProcessor(svcCtx *svc.ServiceContext, ctx context.Context) *TaskBotSendProcessor {
	//使用计数器限制10秒内最多返回一次提示
	limiter := limit.NewPeriodLimit(10, 1, svcCtx.Redis, model.UserUploadLimit)
	return &TaskBotSendProcessor{
		svcCtx:  svcCtx,
		ctx:     ctx,
		limiter: limiter,
	}
}

func (p *TaskBotSendProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	update := &models.Update{}
	err := json.Unmarshal(t.Payload(), update)
	if err != nil {
		return err
	}
	//logx.Info("Processing task:", update.ID, update.Message.ID)
	if update == nil || (update.Message == nil && update.CallbackQuery == nil) {
		return nil
	}
	languageCode := model.LANGUAGE_EN

	//按钮回调解析
	if update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, "/") {
		p.CallbackQuery(update)
		return nil
	}
	u, err := p.svcCtx.Rpc.FindOneByUid(ctx, &pb.UidReq{Uid: update.Message.From.ID})
	if err != nil {
		logx.Error("find user by uid error", err)
		return err
	}
	if u.Id > 0 {
		//更新用户头像
		if u.LanguageCode == model.LANGUAGE_EN || u.LanguageCode == model.LANGUAGE_TW {
			languageCode = u.LanguageCode
		}
	}
	//解析命令
	if update.Message != nil && strings.HasPrefix(update.Message.Text, "/") || (update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, "/")) {
		p.SendCommand(update, languageCode)
		return nil
	}

	if u.Id == 0 {
		if languageCode == "en" {
			p.sendMsg(0, update.Message.Chat.ID, "", p.svcCtx.Config.Telegram.FirstTextEn, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FirstButEn, URL: p.svcCtx.Config.Telegram.TgUrl}}})
		} else {
			p.sendMsg(0, update.Message.Chat.ID, "", p.svcCtx.Config.Telegram.FirstText, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FirstBut, URL: p.svcCtx.Config.Telegram.TgUrl}}})
		}
		return nil
	}

	var fileSize int64
	var fileName string
	var fileId string
	isFile := 0
	//获取文件
	fileId, fileName, fileSize, isFile = isFileType(update)

	if isFile > 0 {
		//检查用户当前的上传数量
		maxUpload, _ := p.svcCtx.Redis.ScardCtx(ctx, fmt.Sprintf("%s%d", model.UserUpload, update.Message.From.ID))
		if maxUpload >= p.svcCtx.Config.MaxUpload {
			take, err := p.limiter.Take(update.Message.MediaGroupID)
			if err != nil {
				logx.Error("Error:", err)
				return nil
			}
			if take == limit.HitQuota {
				if languageCode == "en" {
					_, err = p.svcCtx.TgBot.SendMessage(context.Background(), &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   p.svcCtx.Config.Telegram.LimitTextEn,
					})
				} else {
					_, err = p.svcCtx.TgBot.SendMessage(context.Background(), &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   p.svcCtx.Config.Telegram.LimitText,
					})
				}
			}
			return nil
		}
		// 校验文件大小不能超过限制
		/*if fileSize > p.svcCtx.Config.MaxBytes {
			if languageCode == "en" {
				p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.FailTipEn, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FailSpaceButEn, URL: p.svcCtx.Config.Telegram.TgUrl}}})
			} else {
				p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.FailTip, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FailSpaceBut, URL: p.svcCtx.Config.Telegram.TgUrl}}})
			}
			return nil
		}*/
		//校验用户剩余存储空间
		s, err := p.svcCtx.Rpc.FindOneUserStorage(ctx, &pb.FindOneUserStorageReq{Uid: update.Message.From.ID})
		if err != nil {
			logx.Error("find user by uid error", err)
			return err
		}
		if fileSize > s.SurStorage {
			if languageCode == "en" {
				p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.FailSpaceTextEn, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FailSpaceButEn, URL: p.svcCtx.Config.Telegram.TgUrl + "?startapp=space"}}})
			} else {
				p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.FailSpaceText, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FailSpaceBut, URL: p.svcCtx.Config.Telegram.TgUrl + "?startapp=space"}}})
			}
			return nil
		}

		//查询用户文件夹
		pid, err := p.AssetsFolder(model.MyTelegram, update.Message.From.ID)
		if err != nil {
			return fmt.Errorf("error finding default assets: %w", err)
		}

		// 下载照片
		getFile, err := p.svcCtx.TgBot.GetFile(ctx, &bot.GetFileParams{FileID: fileId})
		if err != nil {
			logx.Error("Error:", err)
			if languageCode == "en" {
				p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.FailTextEn, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FailButEn, URL: p.svcCtx.Config.Telegram.TgUrl + "?startapp=folder"}}})
			} else {
				p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.FailText, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FailBut, URL: p.svcCtx.Config.Telegram.TgUrl + "?startapp=folder"}}})
			}
			return nil
		}
		if fileName == "" {
			prefix := strings.TrimPrefix(getFile.FilePath, p.svcCtx.Config.Telegram.PrefixFilePath)
			fileName = prefix
		}

		// 保存资产信息
		asset, err := p.svcCtx.Rpc.SaveAssetFile(ctx, &pb.SaveAssetFileReq{
			Uid:       update.Message.From.ID,
			AssetName: fileName,
			AssetSize: getFile.FileSize,
			AssetType: utils.IsFileType(fileName),
			Pid:       pid,
			Source:    3,
			Path:      fileId,
		})
		if err != nil {
			if languageCode == "en" {
				p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.FailTextEn, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FailButEn, URL: p.svcCtx.Config.Telegram.TgUrl + "?startapp=folder"}}})
			} else {
				p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.FailText, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FailBut, URL: p.svcCtx.Config.Telegram.TgUrl + "?startapp=folder"}}})
			}
			return nil
		}

		// 将文件ID添加到上传集合中
		p.svcCtx.Redis.SaddCtx(ctx, fmt.Sprintf("%s%d", model.UserUpload, update.Message.From.ID), asset.Id)
		//设置过期时间
		p.svcCtx.Redis.Expire(fmt.Sprintf("%s%d", model.UserUpload, update.Message.From.ID), p.svcCtx.Config.UploadExpireTime)
		filePath := getFile.FilePath
		if strings.HasPrefix(getFile.FilePath, p.svcCtx.Config.Telegram.PrefixFilePath) {
			filePath = strings.TrimPrefix(getFile.FilePath, p.svcCtx.Config.Telegram.PrefixFilePath)
		}
		link := p.svcCtx.TgBot.FileDownloadLink(filePath)
		logx.Info("download link:", link)
		//刪除上传集合中的文件ID
		defer func(Redis *redis.Redis, key string, values ...any) {
			_, err := Redis.Srem(key, values)
			if err != nil {
				return
			}
		}(p.svcCtx.Redis, fmt.Sprintf("%s%d", model.UserUpload, update.Message.From.ID), asset.Id)

		assetID, err := p.upload(link, asset.Id, fileSize, u.Uid)
		if err != nil {
			//上传失败
			if languageCode == "en" {
				p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.FailTextEn, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FailButEn, URL: p.svcCtx.Config.Telegram.TgUrl + "?startapp=folder"}}})
			} else {
				p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.FailText, [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FailBut, URL: p.svcCtx.Config.Telegram.TgUrl + "?startapp=folder"}}})
			}
			return nil
		}
		//上传成功
		if languageCode == "en" {
			p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.SuccessTextEn, [][]models.InlineKeyboardButton{
				{{Text: p.svcCtx.Config.Telegram.SuccessButEn, URL: p.svcCtx.Config.Telegram.TgUrl + fmt.Sprintf("?startapp=transmission-id=%v", assetID)}},
			})
		} else {
			p.sendMsg(isFile, update.Message.Chat.ID, fileId, p.svcCtx.Config.Telegram.SuccessText, [][]models.InlineKeyboardButton{
				{{Text: p.svcCtx.Config.Telegram.SuccessBut, URL: p.svcCtx.Config.Telegram.TgUrl + fmt.Sprintf("?startapp=transmission-id=%v", assetID)}},
			})
		}
		return nil
	}

	return nil
}

// upload 负责将数据从reader中上传，并记录上传进度。
// 参数:
// reader: 数据源，通常是一个文件或数据流。
// name: 上传数据的名称，通常用于存储时的文件名。
// id: 上传数据的唯一标识符，用于跟踪和更新上传状态。
func (m *TaskBotSendProcessor) upload(link string, id string, fileSize, uid int64) (int64, error) {
	ctx := context.Background()

	logx.Errorf("uid:%v", uid)
	storageCli, err := m.svcCtx.GetStorage(uid)
	if err != nil {
		logx.Error("获取存储失败", err)
		return 0, err
	}

	// progress 是一个回调函数，用于在上传过程中更新上传进度。
	progress := func(doneSize int64, totalSize int64) {
		// 如果上传完成，从Redis中删除上传进度记录，并打印成功信息。
		if doneSize == totalSize {
			_, err := m.svcCtx.Redis.Del(fmt.Sprintf(model.UploadId+"%s", id))
			if err != nil {
				logx.Error("redis删除上传进度失败", err)
			}
		}
		// 更新Redis中的上传进度。
		err := m.svcCtx.Redis.Set(fmt.Sprintf(model.UploadId+"%s", id), strconv.Itoa(int(float64(doneSize)/float64(totalSize)*100)))
		if err != nil {
			logx.Error("redis更新上传进度失败", err)
		}
	}

	// 使用存储服务进行上传。
	cid, url, err := storageCli.UploadFileWithURL(ctx, link, progress)
	if err != nil || url == "" || cid == "" {
		logx.Error("上传失败", err)
		//保存错误到redis
		m.svcCtx.Redis.Set(fmt.Sprintf(model.UploadErrId+"%s", id), err.Error())
		//上传失败
		m.svcCtx.Rpc.UpdateAssetFile(ctx, &pb.UpdateAssetFileReq{
			Id:     id,
			Status: model.AssetStatusError,
		})
		return 0, fmt.Errorf("上傳失敗%s", err)
	}
	// 更新资产状态为可用，并记录上传的CID。
	resp, err := m.svcCtx.Rpc.UpdateAssetFile(ctx, &pb.UpdateAssetFileReq{
		Id:        id,
		Status:    model.AssetStatusEnable,
		Cid:       cid,
		AssetSize: fileSize,
		Link:      url,
	})
	if err != nil {
		// 如果更新失败，记录错误。
		logx.Error("更新失败", err)
		return 0, fmt.Errorf("上傳失敗%s", err)
	}
	return resp.Id, nil
}

func (m *TaskBotSendProcessor) sendMsg(isFile int, chatId int64, fileID string, caption string, inlineKey [][]models.InlineKeyboardButton) (err error) {
	switch isFile {
	case 1:
		_, err = m.svcCtx.TgBot.SendPhoto(context.Background(), &bot.SendPhotoParams{
			ChatID:      chatId,
			Photo:       &models.InputFileString{Data: fileID},
			Caption:     caption,
			ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: inlineKey},
		})
	case 2:
		_, err = m.svcCtx.TgBot.SendDocument(context.Background(), &bot.SendDocumentParams{
			ChatID:      chatId,
			Document:    &models.InputFileString{Data: fileID},
			Caption:     caption,
			ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: inlineKey},
		})
	case 3:
		_, err = m.svcCtx.TgBot.SendVideo(context.Background(), &bot.SendVideoParams{
			ChatID:      chatId,
			Video:       &models.InputFileString{Data: fileID},
			Caption:     caption,
			ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: inlineKey},
		})
	case 4:
		_, err = m.svcCtx.TgBot.SendVoice(context.Background(), &bot.SendVoiceParams{
			ChatID:      chatId,
			Voice:       &models.InputFileString{Data: fileID},
			Caption:     caption,
			ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: inlineKey},
		})
	case 5:
		_, err = m.svcCtx.TgBot.SendAudio(context.Background(), &bot.SendAudioParams{
			ChatID:      chatId,
			Audio:       &models.InputFileString{Data: fileID},
			Caption:     caption,
			ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: inlineKey},
		})
	default:
		_, err = m.svcCtx.TgBot.SendMessage(context.Background(), &bot.SendMessageParams{
			ChatID:      chatId,
			Text:        caption,
			ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: inlineKey},
		})
	}
	return nil
}

func isFileType(update *models.Update) (fileId string, fileName string, fileSize int64, isFile int) {
	switch {
	case len(update.Message.Photo) > 0: //图片
		// 获取最新的照片
		photo := update.Message.Photo[len(update.Message.Photo)-1]
		fileSize = int64(photo.FileSize)
		fileId = photo.FileID
		isFile = 1
	case update.Message.Document != nil: //文档
		document := update.Message.Document
		fileSize = document.FileSize
		fileId = document.FileID
		fileName = document.FileName
		isFile = 2
	case update.Message.Video != nil: //视频
		video := update.Message.Video
		fileSize = video.FileSize
		fileId = video.FileID
		fileName = video.FileName
		isFile = 3
	case update.Message.Voice != nil: //音频
		voice := update.Message.Voice
		fileSize = voice.FileSize
		fileId = voice.FileID
		isFile = 4
	case update.Message.Audio != nil: //语音
		audio := update.Message.Audio
		fileSize = audio.FileSize
		fileId = audio.FileID
		fileName = audio.FileName
		isFile = 5
	}
	return
}

// 查询用户文件夹
func (m *TaskBotSendProcessor) AssetsFolder(assetsName string, uid int64) (pid int64, err error) {
	assets, err := m.svcCtx.Rpc.FindOneAssets(context.TODO(), &pb.FindOneAssetsReq{Uid: uid, AssetName: assetsName, AssetType: 1})
	if err != nil {
		return 1, err
	}
	if assets.Id == 0 {
		//保存文件夹
		assetsResp, err := m.svcCtx.Rpc.SaveAssets(context.TODO(), &pb.SaveAssetsReq{
			Uid:         uid,
			AssetName:   assetsName,
			AssetType:   1,
			TransitType: 1,
			IsTag:       2,
			Pid:         1,
			Source:      1,
			Status:      model.AssetStatusEnable,
		})
		if err != nil {
			return 1, err
		}
		assets.Id = assetsResp.Id
	}
	return assets.Id, nil
}

// 回复命令
func (p *TaskBotSendProcessor) SendCommand(update *models.Update, languageCode string) (err error) {

	command, err := p.svcCtx.Rpc.FindOneBotCommand(p.ctx, &pb.FindOneBotCommandReq{BotCommand: update.Message.Text, LanguageCode: languageCode})
	if err != nil {
		logx.Error("find command error", err)
		return nil
	}
	if command.Id == 0 {
		return nil
	}
	text := command.Text
	//配置模版消息
	switch {
	case strings.Contains(command.Text, "{{.username}}"):
		parse, err := template.New(strconv.FormatInt(command.Id, 10)).Parse(command.Text)
		if err != nil {
			logx.Error("parse error", err)
			return nil
		}
		data := map[string]string{
			"username": update.Message.From.Username,
		}
		var out bytes.Buffer
		err = parse.Execute(&out, data)
		if err != nil {
			fmt.Println("Error executing template:", err)
		}
		text = out.String()
	case strings.Contains(command.Text, "{{.userLink}}"):
		parse, err := template.New(strconv.FormatInt(command.Id, 10)).Parse(command.Text)
		if err != nil {
			logx.Error("parse error", err)
			return nil
		}
		user, err := p.svcCtx.Rpc.FindOneByUid(p.ctx, &pb.UidReq{Uid: update.Message.From.ID})
		if err != nil {
			logx.Error("find user error", err)
			return nil
		}
		if user.Id == 0 {
			if languageCode == "en" {
				_, err = p.svcCtx.TgBot.SendMessage(context.Background(), &bot.SendMessageParams{
					ChatID:      update.Message.Chat.ID,
					Text:        p.svcCtx.Config.Telegram.FirstTextEn,
					ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FirstButEn, URL: p.svcCtx.Config.Telegram.TgUrl}}}},
				})
			} else {
				_, err = p.svcCtx.TgBot.SendMessage(context.Background(), &bot.SendMessageParams{
					ChatID:      update.Message.Chat.ID,
					Text:        p.svcCtx.Config.Telegram.FirstText,
					ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FirstBut, URL: p.svcCtx.Config.Telegram.TgUrl}}}},
				})
			}
			return nil
		}
		data := map[string]string{
			"userLink": p.svcCtx.Config.Telegram.TgUrl + "?" + "startapp=" + user.RecommendCode,
		}
		var out bytes.Buffer
		err = parse.Execute(&out, data)
		if err != nil {
			fmt.Println("Error executing template:", err)
		}
		text = out.String()
	}

	//配置回复按钮
	var InlineKeyboardMarkup [][]models.InlineKeyboardButton
	if command.ButtonArray != "" {
		replyMarkups := make([][]types.Markup, 0)
		err = json.Unmarshal([]byte(command.ButtonArray), &replyMarkups)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return nil
		}
		InlineKeyboardMarkup = make([][]models.InlineKeyboardButton, 0, len(replyMarkups))
		for _, markup := range replyMarkups {
			inlineKeyboardMarkup := make([]models.InlineKeyboardButton, 0, len(markup))
			for _, item := range markup {
				//分享链接特殊处理
				if item.Button == "分享给朋友" || item.Button == "Share with friends" {
					user, err := p.svcCtx.Rpc.FindOneByUid(p.ctx, &pb.UidReq{Uid: update.Message.From.ID})
					if err != nil {
						logx.Error("find user error", err)
						return nil
					}
					if user.Id == 0 {
						if languageCode == "en" {
							_, err = p.svcCtx.TgBot.SendMessage(context.Background(), &bot.SendMessageParams{
								ChatID:      update.Message.Chat.ID,
								Text:        p.svcCtx.Config.Telegram.FirstTextEn,
								ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FirstButEn, URL: p.svcCtx.Config.Telegram.TgUrl}}}},
							})
						} else {
							_, err = p.svcCtx.TgBot.SendMessage(context.Background(), &bot.SendMessageParams{
								ChatID:      update.Message.Chat.ID,
								Text:        p.svcCtx.Config.Telegram.FirstText,
								ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: p.svcCtx.Config.Telegram.FirstBut, URL: p.svcCtx.Config.Telegram.TgUrl}}}},
							})
						}
						return nil
					}
					if languageCode == "en" {
						item.Url = fmt.Sprintf("%s?text=%s&url=%s", p.svcCtx.Config.Telegram.ShareUrl, p.svcCtx.Config.Telegram.ShareTextEn, p.svcCtx.Config.Telegram.TgUrl+"?startapp="+user.RecommendCode)
					} else {
						item.Url = fmt.Sprintf("%s?text=%s&url=%s", p.svcCtx.Config.Telegram.ShareUrl, p.svcCtx.Config.Telegram.ShareText, p.svcCtx.Config.Telegram.TgUrl+"?startapp="+user.RecommendCode)
					}
				}
				inlineKeyboardMarkup = append(inlineKeyboardMarkup, models.InlineKeyboardButton{
					Text:         item.Button,
					URL:          item.Url,
					CallbackData: item.CallbackData,
				})
			}
			InlineKeyboardMarkup = append(InlineKeyboardMarkup, inlineKeyboardMarkup)
		}
	}
	message := &models.Message{}
	switch command.SendType {
	case 2:
		if command.ButtonArray == "" {
			message, err = p.svcCtx.TgBot.SendMessage(p.ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   text,
			})
			if err != nil {
				logx.Error(err)
				return nil
			}
		} else {
			message, err = p.svcCtx.TgBot.SendMessage(p.ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        text,
				ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: InlineKeyboardMarkup},
			})
			if err != nil {
				logx.Error(err)
				return nil
			}
		}
		//图片
	case 4:
		decode, err := base64.Decode([]byte(command.Photo))
		if err != nil {
			logx.Error("decode error", err)
			return nil
		}
		if command.ButtonArray != "" {
			update.Message.ReplyMarkup = models.InlineKeyboardMarkup{InlineKeyboard: InlineKeyboardMarkup}
			message, err = p.svcCtx.TgBot.SendPhoto(p.ctx, &bot.SendPhotoParams{
				ChatID:      update.Message.Chat.ID,
				Photo:       &models.InputFileUpload{Data: bytes.NewReader(decode)},
				Caption:     text,
				ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: InlineKeyboardMarkup},
			})
			if err != nil {
				logx.Error(err)
				return nil
			}
		} else {
			message, err = p.svcCtx.TgBot.SendPhoto(p.ctx, &bot.SendPhotoParams{
				ChatID:  update.Message.Chat.ID,
				Photo:   &models.InputFileUpload{Data: bytes.NewReader(decode)},
				Caption: text,
			})
			if err != nil {
				logx.Error(err)
				return nil
			}
		}
		//视频
	case 3:
		decode, err := base64.Decode([]byte(command.Photo))
		if err != nil {
			logx.Error("decode error", err)
			return nil
		}
		if command.ButtonArray != "" {
			update.Message.ReplyMarkup = models.InlineKeyboardMarkup{InlineKeyboard: InlineKeyboardMarkup}
			message, err = p.svcCtx.TgBot.SendAnimation(p.ctx, &bot.SendAnimationParams{
				ChatID:      update.Message.Chat.ID,
				Animation:   &models.InputFileUpload{Data: bytes.NewReader(decode), Filename: "animation.gif"},
				Caption:     text,
				ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: InlineKeyboardMarkup},
			})
			if err != nil {
				logx.Error(err)
				return nil
			}
		} else {
			message, err = p.svcCtx.TgBot.SendAnimation(p.ctx, &bot.SendAnimationParams{
				ChatID:    update.Message.Chat.ID,
				Animation: &models.InputFileUpload{Data: bytes.NewReader(decode), Filename: "animation.gif"},
				Caption:   text,
			})
			if err != nil {
				logx.Error(err)
				return nil
			}
		}
	}
	if command.Command == "/settings" {
		marshal, err := json.Marshal(message)
		if err != nil {
			logx.Error("marshal error", err)
			return nil
		}
		p.svcCtx.Redis.Set(fmt.Sprintf("%s%v", model.UserSettings, message.From.ID), string(marshal))
		p.svcCtx.Redis.Expire(fmt.Sprintf("%s%v", model.UserSettings, message.From.ID), 120)
	}
	return nil
}

// 按钮回调
func (p *TaskBotSendProcessor) CallbackQuery(update *models.Update) (err error) {

	u, err := p.svcCtx.Rpc.FindOneByUid(p.ctx, &pb.UidReq{Uid: update.CallbackQuery.From.ID})
	if err != nil {
		logx.Error("find user by uid error", err)
		return err
	}
	languageCode := model.LANGUAGE_EN
	if u.Id > 0 {
		if u.LanguageCode == model.LANGUAGE_EN || u.LanguageCode == model.LANGUAGE_TW {
			languageCode = u.LanguageCode
		}
	}

	command, err := p.svcCtx.Rpc.FindOneBotCommand(p.ctx, &pb.FindOneBotCommandReq{BotCommand: update.CallbackQuery.Data, LanguageCode: languageCode})
	if err != nil {
		logx.Error("find command error", err)
		return nil
	}
	if command.Id == 0 {
		return nil
	}
	text := command.Text

	message := &models.Message{}
	//配置回复按钮
	var InlineKeyboardMarkup [][]models.InlineKeyboardButton
	if command.ButtonArray != "" {
		replyMarkups := make([][]types.Markup, 0)
		err = json.Unmarshal([]byte(command.ButtonArray), &replyMarkups)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return nil
		}
		InlineKeyboardMarkup = make([][]models.InlineKeyboardButton, 0, len(replyMarkups))
		for _, markup := range replyMarkups {
			inlineKeyboardMarkup := make([]models.InlineKeyboardButton, 0, len(markup))
			for _, item := range markup {
				inlineKeyboardMarkup = append(inlineKeyboardMarkup, models.InlineKeyboardButton{
					Text:         item.Button,
					URL:          item.Url,
					CallbackData: item.CallbackData,
				})
			}
			InlineKeyboardMarkup = append(InlineKeyboardMarkup, inlineKeyboardMarkup)
		}
	}
	if command.Photo == "" {
		if command.ButtonArray == "" {
			message, err = p.svcCtx.TgBot.SendMessage(p.ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.Message.Message.Chat.ID,
				Text:   text,
			})
			if err != nil {
				logx.Error(err)
				return nil
			}
		} else {
			message, err = p.svcCtx.TgBot.SendMessage(p.ctx, &bot.SendMessageParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				Text:        text,
				ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: InlineKeyboardMarkup},
			})
			if err != nil {
				logx.Error(err)
				return nil
			}
		}
	} else {
		decode, err := base64.Decode([]byte(command.Photo))
		if err != nil {
			logx.Error("decode error", err)
			return nil
		}
		if command.ButtonArray != "" {
			message, err = p.svcCtx.TgBot.SendPhoto(p.ctx, &bot.SendPhotoParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				Photo:       &models.InputFileUpload{Data: bytes.NewReader(decode)},
				Caption:     text,
				ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: InlineKeyboardMarkup},
			})
			if err != nil {
				logx.Error(err)
				return nil
			}
		} else {
			message, err = p.svcCtx.TgBot.SendPhoto(p.ctx, &bot.SendPhotoParams{
				ChatID:  update.CallbackQuery.Message.Message.Chat.ID,
				Photo:   &models.InputFileUpload{Data: bytes.NewReader(decode)},
				Caption: text,
			})
			if err != nil {
				logx.Error(err)
				return nil
			}
		}
	}
	//设置语言
	switch update.CallbackQuery.Data {
	case "/language":
		p.DelMsg(update.CallbackQuery.Message.Message.From.ID, model.UserSettings)
		//缓存语言消息
		marshal, err := json.Marshal(message)
		if err != nil {
			logx.Error("marshal error", err)
			return nil
		}
		p.svcCtx.Redis.Set(fmt.Sprintf("%s%v", model.UserLanguage, message.From.ID), string(marshal))
		p.svcCtx.Redis.Expire(fmt.Sprintf("%s%v", model.UserLanguage, message.From.ID), 120)
		return nil
	case "/language_zh":
		if u.Id > 0 {
			p.svcCtx.Rpc.SaveUser(p.ctx, &pb.User{Id: u.Id, LanguageCode: model.LANGUAGE_TW})
		}

		p.DelMsg(update.CallbackQuery.Message.Message.From.ID, model.UserLanguage)
		p.svcCtx.TgBot.SendMessage(p.ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "資訊已更新",
		})

		return nil
	case "/language_en":
		if u.Id > 0 {
			p.svcCtx.Rpc.SaveUser(p.ctx, &pb.User{Id: u.Id, LanguageCode: model.LANGUAGE_EN})
		}

		p.DelMsg(update.CallbackQuery.Message.Message.From.ID, model.UserLanguage)
		p.svcCtx.TgBot.SendMessage(p.ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Information has been updated",
		})
		return nil
	}
	return nil
}

// 删除消息
func (p *TaskBotSendProcessor) DelMsg(uid int64, code string) {
	language, err := p.svcCtx.Redis.Get(fmt.Sprintf("%s%v", code, uid))
	if err != nil || err == redis.Nil {
		logx.Error("get message error", err)
	}
	message := &models.Message{}
	json.Unmarshal([]byte(language), message)
	p.svcCtx.TgBot.DeleteMessage(p.ctx, &bot.DeleteMessageParams{
		ChatID:    message.Chat.ID,
		MessageID: message.ID,
	})
	p.svcCtx.Redis.Del(fmt.Sprintf("%s%v", code, uid))
}
