package svc

import (
	"T-driver/app/resource/api/internal/config"
	"T-driver/app/resource/api/internal/middleware"
	"T-driver/app/tgbot/rpc/tgbot"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/app/user/rpc/user"
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/hibiken/asynq"
	storage "github.com/utopiosphe/titan-storage-sdk"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
	Redis          *redis.Redis
	Rpc            user.UserZrpcClient
	BotRpc         tgbot.Tgbot
	Storage        storage.Storage
	AsynqClient    *asynq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCli := redis.MustNewRedis(c.RedisConf)
	storageCli, err := storage.Initialize(&storage.Config{
		TitanURL:    c.TiTan.TitanURL,
		APIKey:      c.TiTan.APIKey,
		GroupID:     0,
		UseFastNode: false,
	})
	if err != nil {
		logx.Must(err)
	}
	storageCli.SetAreas(context.Background(), []string{c.TiTan.AreaID})
	return &ServiceContext{
		Storage:        storageCli,
		Config:         c,
		Redis:          redisCli,
		AuthMiddleware: middleware.NewAuthMiddleware(c.Telegram.Token).Handle,
		Rpc:            user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpc)),
		AsynqClient:    c.AsynqConf.NewClient(),
	}
}

// Upload 负责将数据从reader中上传，并记录上传进度。
// 参数:
// reader: 数据源，通常是一个文件或数据流。
// name: 上传数据的名称，通常用于存储时的文件名。
// id: 上传数据的唯一标识符，用于跟踪和更新上传状态。
func (m *ServiceContext) Upload(reader io.Reader, name string, id int64, assetSize int64) {
	ctx := context.Background()

	// progress 是一个回调函数，用于在上传过程中更新上传进度。
	progress := func(doneSize int64, totalSize int64) {
		// 如果上传完成，从Redis中删除上传进度记录，并打印成功信息。
		if doneSize == totalSize {
			_, err := m.Redis.Del(fmt.Sprintf(model.UploadId+"%v", id))
			if err != nil {
				logx.Error("redis删除上传进度失败", err)
			}
		}
		// 更新Redis中的上传进度。
		err := m.Redis.Set(fmt.Sprintf(model.UploadId+"%v", id), strconv.Itoa(int(float64(doneSize)/float64(totalSize)*100)))
		if err != nil {
			logx.Error("redis更新上传进度失败", err)
		}
	}
	// 使用存储服务进行上传。
	upload, err := m.Storage.UploadStreamV2(ctx, reader, name, progress)
	if err != nil {
		logx.Error("上传失败", err)
		//保存错误到redis
		m.Redis.Set(fmt.Sprintf(model.UploadErrId+"%v", id), err.Error())
		//上传失败
		m.Rpc.UpdateAssetsName(ctx, &pb.UpdateAssetsNameReq{
			Id:     id,
			Status: model.AssetStatusError,
		})
		return
	}
	url := ""
	shareAssetResult, err := m.Storage.GetURL(ctx, upload.String())
	if err != nil {
		logx.Error("获取url失败：", err)
		// todo 放入队列获取
		/*go func(id int64, cid string, m *ServiceContext) {
			payload, err := json.Marshal(&model.StorageGetUrl{
				Id:  id,
				Cid: cid,
			})
			if err != nil {
				logx.WithContext(ctx).Error(err)
			}
			t := asynq.NewTask(model.StorageProcess, payload)
			_, err = m.AsynqClient.EnqueueContext(
				ctx, t,
				asynq.TaskID(fmt.Sprintf("%s:%v", model.StorageProcess, id)),
				asynq.Queue(model.QueueLow),
			)
			if err != nil {
				logx.WithContext(ctx).Error(err)
			}
		}(id, upload.String(), m)*/
	}
	if shareAssetResult != nil && len(shareAssetResult.URLs) > 0 {
		url = shareAssetResult.URLs[0]
	}
	// 更新资产状态为可用，并记录上传的CID。
	_, err = m.Rpc.UpdateAssetsName(ctx, &pb.UpdateAssetsNameReq{
		Id:        id,
		Status:    model.AssetStatusEnable,
		Cid:       upload.String(),
		AssetSize: assetSize,
		Link:      url,
	})
	if err != nil {
		// 如果更新失败，记录错误。
		logx.Error("更新失败", err)
	}
}

// Upload 负责将数据从reader中上传，并记录上传进度。
// 参数:
// reader: 数据源，通常是一个文件或数据流。
// name: 上传数据的名称，通常用于存储时的文件名。
// id: 上传数据的唯一标识符，用于跟踪和更新上传状态。
func (m *ServiceContext) LinkUpload(link string, id int64, assetSize int64) {
	ctx := context.Background()

	// progress 是一个回调函数，用于在上传过程中更新上传进度。
	progress := func(doneSize int64, totalSize int64) {
		// 如果上传完成，从Redis中删除上传进度记录，并打印成功信息。
		if doneSize == totalSize {
			_, err := m.Redis.Del(fmt.Sprintf(model.UploadId+"%v", id))
			if err != nil {
				logx.Error("redis删除上传进度失败", err)
			}
		}
		// 更新Redis中的上传进度。
		err := m.Redis.Set(fmt.Sprintf(model.UploadId+"%v", id), strconv.Itoa(int(float64(doneSize)/float64(totalSize)*100)))
		if err != nil {
			logx.Error("redis更新上传进度失败", err)
		}
	}
	// 使用存储服务进行上传。
	cid, url, err := m.Storage.UploadFileWithURL(ctx, link, progress)
	if err != nil {
		logx.Error("上传失败", err)
		//保存错误到redis
		m.Redis.Set(fmt.Sprintf(model.UploadErrId+"%v", id), err.Error())
		//上传失败
		m.Rpc.UpdateAssetsName(ctx, &pb.UpdateAssetsNameReq{
			Id:     id,
			Status: model.AssetStatusError,
		})
		return
	}
	/*url := ""
	shareAssetResult, err := m.Storage.GetURL(ctx, upload.String())
	if err != nil {
		logx.Error("获取url失败：", err)
	}
	if shareAssetResult != nil && len(shareAssetResult.URLs) > 0 {
		url = shareAssetResult.URLs[0]
	}*/
	// 更新资产状态为可用，并记录上传的CID。
	_, err = m.Rpc.UpdateAssetsName(ctx, &pb.UpdateAssetsNameReq{
		Id:        id,
		Status:    model.AssetStatusEnable,
		Cid:       cid,
		AssetSize: assetSize,
		Link:      url,
	})
	if err != nil {
		// 如果更新失败，记录错误。
		logx.Error("更新失败", err)
	}
}

// 查询用户文件夹
func (m *ServiceContext) AssetsFolder(assetsName string, uid int64) (pid int64, err error) {
	assets, err := m.Rpc.FindOneAssets(context.TODO(), &pb.FindOneAssetsReq{Uid: uid, AssetName: assetsName, AssetType: 1})
	if err != nil {
		return 1, err
	}
	if assets.Id == 0 {
		//保存文件夹
		assetsResp, err := m.Rpc.SaveAssets(context.TODO(), &pb.SaveAssetsReq{
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
