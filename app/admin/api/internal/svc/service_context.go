package svc

import (
	"T-driver/app/admin/api/internal/config"
	"T-driver/app/admin/api/internal/middleware"
	"T-driver/app/tgbot/rpc/tgbot"
	"T-driver/app/user/rpc/user"
	"T-driver/app/video/rpc/video"
	"T-driver/common/lib/jwt"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	storage "github.com/utopiosphe/titan-storage-sdk"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	Jwt                 jwt.JWT
	AdminAuthMiddleware rest.Middleware
	Redis               *redis.Redis
	Storage             storage.Storage
	Rpc                 user.UserZrpcClient
	BotRpc              tgbot.Tgbot
	VideoRpc            video.VideoZrpcClient
	Sess                *session.Session
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCli := redis.MustNewRedis(c.RedisConf)
	jwtCli := jwt.NewJWT(c.AuthCfg.AccessSecret, c.AuthCfg.AccessExpire, jwt.SetRedis(redisCli), jwt.SetBlackListOpt(true))
	// storageCli, err := storage.NewStorage(&storage.Config{
	// 	TitanURL:    c.TiTan.TitanURL,
	// 	APIKey:      c.TiTan.APIKey,
	// 	GroupID:     0,
	// 	UseFastNode: false,
	// 	AreaID:      c.TiTan.AreaID,
	// })
	// if err != nil {
	// 	logx.Must(err)
	// }
	Sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(c.AwsS3.AccessKey, c.AwsS3.SecretKey, ""),
		Endpoint:         aws.String(c.AwsS3.EndPoint),
		Region:           aws.String(c.AwsS3.Region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false), //virtual-host style方式，不要修改
	})
	if err != nil {
		logx.Must(err)
	}

	return &ServiceContext{
		Config: c,
		// Storage:             storageCli,
		Redis:               redisCli,
		Jwt:                 jwtCli,
		Sess:                Sess,
		AdminAuthMiddleware: middleware.NewAdminAuthMiddleware(jwtCli).Handle,
		Rpc:                 user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpc)),
		BotRpc:              tgbot.NewTgbot(zrpc.MustNewClient(c.BotRpc)),
		VideoRpc:            video.NewVideoZrpcClient(zrpc.MustNewClient(c.VideoRpc)),
	}
}

func (s *ServiceContext) BucketUpload(filename, name string) (location string, err error) {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Unable to open file , fmt.Println", err)
		return "", err
	}
	defer file.Close()

	bucket := fmt.Sprintf("/%v", s.Config.AwsS3.Bucket)
	uploader := s3manager.NewUploader(s.Sess)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
		Body:   file,
	})
	if err != nil {
		// Print the error and exit.
		fmt.Println("BucketUpload err ", err)
		return "", err
	}
	/*if res.Location != "" {
		//修改资源路径
		location = strings.Replace(res.Location, end_point+""+bucket, "/resources", 1)
	}*/
	return res.Location, nil
}
