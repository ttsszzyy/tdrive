package svc

import (
	"T-driver/app/tgbot/rpc/internal/config"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/app/user/rpc/user"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	storage "github.com/utopiosphe/titan-storage-sdk"

	"github.com/go-telegram/bot"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	TgBot       *bot.Bot
	AsynqServer *asynq.Server
	Rpc         user.UserZrpcClient
	Redis       *redis.Redis
	Tenant      storage.Tenant
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCli := redis.MustNewRedis(c.RedisConf)
	tenant, err := storage.NewTenant(c.TiTan.TitanURL, c.TiTan.APIKey)
	if err != nil {
		logx.Must(err)
	}
	serverQueue := map[string]int{
		model.QueueCritical: 5,
		model.QueueNormal:   3,
		model.QueueLow:      1,
	}
	svcCtx := &ServiceContext{
		Config:      c,
		Tenant:      tenant,
		Redis:       redisCli,
		AsynqServer: c.AsynqConf.NewServer(serverQueue, c.AsynqConf.Concurrency),
		Rpc:         user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpc)),
	}
	b, err := bot.New(svcCtx.Config.Telegram.Url, svcCtx.Config.Telegram.Token)
	if nil != err {
		logx.Must(fmt.Errorf("链接bot失败！%s", err))
	}
	svcCtx.TgBot = b
	return svcCtx
}

func (m *ServiceContext) GetStorage(uid int64) (st storage.Storage, err error) {
	ctx := context.TODO()
	getCtx, err := m.Redis.GetCtx(ctx, fmt.Sprintf(model.GetStorageUser+"%v", uid))
	if err != nil && err != redis.Nil {
		return nil, err
	}
	one := &pb.UserTitanToken{}
	if err == redis.Nil || getCtx == "" {
		one, err = m.Rpc.FindOneUserTitanToken(ctx, &pb.FindOneUserTitanTokenReq{Uid: uid})
		if err != nil {
			return nil, err
		}
		// if one.Uid == 0 {
		// 	u, err := m.Rpc.FindOneByUid(ctx, &pb.UidReq{Uid: uid})
		// 	if err != nil {
		// 		logx.Error("find user by uid error", err)
		// 		return nil, err
		// 	}
		// 	login, err := m.Tenant.SSOLogin(ctx, storage.SubUserInfo{
		// 		EntryUUID: strconv.FormatInt(u.Uid, 10),
		// 		Username:  u.Name,
		// 		Avatar:    u.Avatar,
		// 	})
		// 	if err != nil {
		// 		logx.Errorf("login titan error:%v name:%v uid:%v", err, u.Name, u.Uid)
		// 		return nil, err
		// 	}
		// 	one.Token = login.Token
		// 	one.Expire = login.Exp
		// 	one.Uid = u.Uid
		// 	m.Rpc.UpdateUserTitanToken(ctx, &pb.UpdateUserTitanTokenReq{
		// 		Uid:    u.Uid,
		// 		Token:  login.Token,
		// 		Expire: login.Exp,
		// 	})
		// }
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
		logx.Info("token过期，正在续期")
		user, err := m.Rpc.FindOneByUid(ctx, &pb.UidReq{Uid: uid})
		if err != nil {
			return nil, err
		}
		token, err := m.Tenant.SSOLogin(ctx, storage.SubUserInfo{
			EntryUUID: strconv.FormatInt(one.Uid, 10),
			Username:  user.Name,
			Avatar:    user.Avatar,
		})
		if err != nil {
			logx.Error("token续期失败", err)
			return nil, err
		}
		_, err = m.Rpc.UpdateUserTitanToken(ctx, &pb.UpdateUserTitanTokenReq{
			Uid:    uid,
			Token:  token.Token,
			Expire: token.Exp,
		})
		if err != nil {
			logx.Error("token更新失败", err)
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
	logx.Error("token", one.Token)
	storageCli, err := storage.Initialize(&storage.Config{
		TitanURL: m.Config.TiTan.TitanURL,
		Token:    one.Token,
	})
	if err != nil {
		return nil, err
	}
	return storageCli, nil
}
