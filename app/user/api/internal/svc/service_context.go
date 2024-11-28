package svc

import (
	"T-driver/app/tgbot/rpc/tgbot"
	"T-driver/app/user/api/internal/config"
	"T-driver/app/user/api/internal/middleware"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/app/user/rpc/user"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

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
	Tenant         storage.Tenant
	AsynqClient    *asynq.Client

	PrvKey *rsa.PrivateKey
	PubKey string
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCli := redis.MustNewRedis(c.RedisConf)
	tenant, err := storage.NewTenant(c.TiTan.TitanURL, c.TiTan.APIKey)
	if err != nil {
		logx.Must(err)
	}
	prk, pbk := generateKeys()
	svcCtx := &ServiceContext{
		Tenant:         tenant,
		Config:         c,
		Redis:          redisCli,
		AuthMiddleware: middleware.NewAuthMiddleware(c.Telegram.Token).Handle,
		Rpc:            user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpc)),
		BotRpc:         tgbot.NewTgbot(zrpc.MustNewClient(c.BotRpc)),
		AsynqClient:    c.AsynqConf.NewClient(),
		PrvKey:         prk,
		PubKey:         pbk,
	}

	return svcCtx
}

// Upload 负责将数据从reader中上传，并记录上传进度。
// 参数:
// reader: 数据源，通常是一个文件或数据流。
// name: 上传数据的名称，通常用于存储时的文件名。
// id: 上传数据的唯一标识符，用于跟踪和更新上传状态。
/*func (m *ServiceContext) Upload(reader io.Reader, name string, id int64, assetSize int64) {
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
		logx.Error("：获取url失败", err)
		//todo 放入到队列中去
		go func(id int64, cid string, m *ServiceContext) {
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
		}(id, upload.String(), m)


		m.Rpc.UpdateAssetsName(ctx, &pb.UpdateAssetsNameReq{
			Id:     id,
			Status: model.AssetStatusError,
		})
		return*
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
}*/

// Upload 负责将数据从reader中上传，并记录上传进度。
// 参数:
// reader: 数据源，通常是一个文件或数据流。
// name: 上传数据的名称，通常用于存储时的文件名。
// id: 上传数据的唯一标识符，用于跟踪和更新上传状态。
/*func (m *ServiceContext) LinkUpload(link string, id int64, assetSize int64) {
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
	url := ""
	shareAssetResult, err := m.Storage.GetURL(ctx, upload.String())
	if err != nil {
		logx.Error("获取url失败：", err)
	}
	if shareAssetResult != nil && len(shareAssetResult.URLs) > 0 {
		url = shareAssetResult.URLs[0]
	}
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
}*/

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

func (m *ServiceContext) GetStorage(uid int64) (st storage.Storage, err error) {
	ctx := context.TODO()
	getCtx, err := m.Redis.GetCtx(ctx, fmt.Sprintf(model.GetStorageUser+"%v", uid))
	if err != nil && err != redis.Nil {
		return nil, err
	}
	one := &pb.UserTitanToken{}
	if err == redis.Nil || getCtx == "" || getCtx == "{}" {
		one, err = m.Rpc.FindOneUserTitanToken(ctx, &pb.FindOneUserTitanTokenReq{Uid: uid})
		if err != nil {
			return nil, err
		}
		marshal, _ := json.Marshal(one)
		//更新redis
		err = m.Redis.SetCtx(ctx, fmt.Sprintf(model.GetStorageUser+"%v", uid), string(marshal))
		if err != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal([]byte(getCtx), one)
		if err != nil {
			return nil, err
		}
	}
	//续期
	if time.Now().Unix() >= one.Expire {
		// token, err := m.Tenant.RefreshToken(ctx, one.Token)
		// if err != nil {
		// 	return nil, err
		// }
		u, err := m.Rpc.FindOneByUid(ctx, &pb.UidReq{Uid: uid})
		if err != nil {
			logx.Error("find user by uid error", err)
			return nil, err
		}
		token, err := m.Tenant.SSOLogin(ctx, storage.SubUserInfo{
			EntryUUID: strconv.FormatInt(u.Uid, 10),
			Username:  u.Name,
			Avatar:    u.Avatar,
		})
		if err != nil {
			logx.Errorf("login titan error:%v name:%v uid:%v", err, u.Name, u.Uid)
			return nil, err
		}
		_, err = m.Rpc.UpdateUserTitanToken(ctx, &pb.UpdateUserTitanTokenReq{
			Uid:    uid,
			Token:  token.Token,
			Expire: token.Exp,
		})
		if err != nil {
			return nil, err
		}
		one.Token = token.Token
		one.Expire = token.Exp
		marshal, _ := json.Marshal(one)
		//更新redis
		err = m.Redis.SetCtx(ctx, fmt.Sprintf(model.GetStorageUser+"%v", uid), string(marshal))
		if err != nil {
			return nil, err
		}
	}
	storageCli, err := storage.Initialize(&storage.Config{
		TitanURL: m.Config.TiTan.TitanURL,
		Token:    one.Token,
	})
	if err != nil {
		return nil, err
	}
	return storageCli, nil
}

// 生成公钥和私钥
func generateKeys() (*rsa.PrivateKey, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}
	// 获取公钥部分并序列化为 PKIX 格式
	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}
	// 将公钥编码为 PEM 格式
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	return privateKey, string(pubPEM)
}

func (m *ServiceContext) GetImgBase64(ctx context.Context, storageCli storage.Storage, assetType int64, cid string) (string, []string) {
	var (
		done      = make(chan bool)
		imgBase64 string
		links     []string
	)

	go func() {
		defer func() {
			done <- true
		}()
		res, err := storageCli.GetURL(ctx, cid)
		if err != nil {
			logx.Error("根据cid获取文件url列表失败", err)
			return
		}
		links = res.URLs
		// 如果是图片则转换为base64进行存储
		if assetType == 4 {
			return
		}
		f, _, err := storageCli.GetFileWithCid(ctx, cid)
		if err != nil {
			logx.Error("根据cid获取文件内容失败", err)
			return
		}
		body, err := io.ReadAll(f)
		if err != nil {
			logx.Error("获取文件内容失败", err)
			return
		}

		imgBase64 = base64.StdEncoding.EncodeToString(body)
	}()

	select {
	case <-time.After(500 * time.Millisecond):
		return "", links
	case <-done:
		return imgBase64, links
	}
}
