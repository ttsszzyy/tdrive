package config

import (
	"T-driver/common/asynq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MaxBytes         int64
	MaxUpload        int64
	UploadExpireTime int
	AsynqConf        asynq.AsynqConf
	Telegram         struct {
		Token           string
		Url             string
		PrefixFilePath  string
		TgUrl           string
		FailText        string //失败提示
		FailTextEn      string //失败提示
		FailBut         string //失败按钮提示
		FailButEn       string //失败按钮提示
		FailSpaceText   string //空间不足失败提示
		FailSpaceTextEn string //空间不足失败提示
		FailSpaceBut    string //空间不足失败按钮提示
		FailSpaceButEn  string //空间不足失败按钮提示
		FailTip         string //失败提示
		FailTipEn       string //失败提示
		SuccessText     string //成功提示
		SuccessTextEn   string //成功提示
		SuccessBut      string //成功按钮提示
		SuccessButEn    string //成功按钮提示
		ShareUrl        string //分享链接
		ShareText       string //分享文案
		ShareTextEn     string //分享文案
		FirstBut        string //首次提示按钮
		FirstButEn      string //首次提示按钮
		FirstText       string //首次提示文案
		FirstTextEn     string //首次提示文案
		LimitText       string //限制提示文案
		LimitTextEn     string //限制提示文案
	}
	RedisConf redis.RedisConf
	UserRpc   zrpc.RpcClientConf
	TiTan     struct {
		TitanURL string
		APIKey   string
	}
}
