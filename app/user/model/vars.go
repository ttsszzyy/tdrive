package model

import (
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound
var ErrInvalidObjectId = errors.New("invalid objectId")

var PassSalt = "T.drive"                 //管理端用户密码加密盐
var UserIntegral = "user_integral"       //用户积分排行
var UploadId = "upload_id:"              //用户上传文件进度
var UploadVideoId = "upload_video_id:"   //短视频上传文件进度
var ClearAsset = "clearAsser"            //清理过期资源
var UserUpload = "userUpload:"           //用戶上傳次數
var BotProcess = "bot:process"           //上传队列
var TitanProcess = "titan:process"       //上传队列
var UploadErrId = "uploadErrId:"         //用户上传失败原因
var UserReward = "userReward:"           //用户幸运奖励
var UserSettings = "userSettings:"       //用户设置命令删除
var UserLanguage = "userLanguage:"       //用户语言删除
var UserUploadLimit = "userUploadLimit:" //用户语言删除
var GetStorageUser = "GetStorageUser:"   //用户Titan Token

// 上传资源的文件夹
var (
	MyTondriver    = "My Tondrive"       //用户默认文件夹
	MyUpload       = "My Upload"         //本地上传
	MyTelegram     = "My Telegram "      //tg上传
	PackFromShared = "Pack From Shared " //分享上传
)

// 通用禁用启用状态
var State_Disable int64 = 1
var State_Enable int64 = 2

type CommonBit int8

var (
	Yes CommonBit = 1
	No  CommonBit = 2

	//队列优先级
	QueueCritical = "critical" // 紧急 5*consumer
	QueueNormal   = "normal"   // 正常 3*consumer
	QueueLow      = "low"      // 慢 1*consumer
)

const (
	SourceInvitation int64 = iota + 1 //引荐
	SourceShare                       //分享
)

const (
	AssetStatusDisable    int64 = iota + 1 //禁用
	AssetStatusAfoot                       //进行中
	AssetStatusEnable                      //完成
	AssetStatusError                       //上传Titan失败
	AssetStatusSpaceError                  //空间不足
)
const (
	LANGUAGE_EN = "en" //英语
	LANGUAGE_TW = "tw" //繁体中文
	LANGUAGE_ES = "es" //西班牙语
	LANGUAGE_FR = "fr" //法语
	LANGUAGE_ID = "id" //印尼语
	LANGUAGE_VI = "vi" //越南语
	LANGUAGE_RU = "ru" //俄语
)

type StorageGetUrl struct {
	Id  int64
	Cid string
}

type GetTGUrl struct {
	Id     string
	FileId string
}
