package task

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/locker"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var C *cron.Cron

func InitTask(svcCtx *svc.ServiceContext) {
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	C = cron.New(cron.WithParser(parser))
	//清理过期资源
	C.AddJob("0 * * */1 * *", cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger)).Then(&ClearAsset{svcCtx: svcCtx}))
	//积分排行 --废弃
	//C.AddJob("0 * */1 * * *", cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger)).Then(&UserRanking{svcCtx: svcCtx}))
	C.Start()
}

type (
	ClearAsset struct {
		svcCtx *svc.ServiceContext
	}
	/*UserRanking struct {
		svcCtx *svc.ServiceContext
	}*/
)

func (this *ClearAsset) Run() {
	defer func() {
		if err := recover(); err != nil {
			logx.Error(fmt.Errorf("ClearAsset err = %s", err))
		}
	}()
	locker := locker.NewRedisLocker(this.svcCtx.Redis, model.ClearAsset)
	defer locker.Unlock()
	if err := locker.Lock(); err != nil {
		logx.Error(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	//查询删除的资源
	assets, err := this.svcCtx.Rpc.FindAssets(ctx, &pb.FindAssetsReq{
		IsDel: true,
	})
	if err != nil {
		logx.Error(err)
		return
	}
	ids := make([]int64, 0)
	for _, asset := range assets.Assets {

		clearTime := time.Unix(asset.DeletedTime, 0).AddDate(0, 0, this.svcCtx.Config.FastReward.RecoveryDate).Unix()
		if clearTime < time.Now().Unix() {
			//删除资源
			if asset.Cid != "" {
				//查询资源是否有其他用户引用，如果没有则删除资源
				oneAssets, err := this.svcCtx.Rpc.CountAssets(ctx, &pb.CountAssetsReq{
					Cid: asset.Cid,
				})
				if err != nil {
					logx.Error("获取资源失败err:", err)
					return
				}
				if oneAssets.Total == 1 {
					storageCli, err := this.svcCtx.GetStorage(asset.Uid)
					if err != nil {
						logx.Error(err)
					}
					//删除资源释放空间
					storageCli.Delete(ctx, asset.Cid)
				}
			}
			ids = append(ids, asset.Id)
		}
	}
	if len(ids) > 0 {
		_, err = this.svcCtx.Rpc.ClearAssets(ctx, &pb.DelAssetsReq{Ids: ids})
		if err != nil {
			logx.Error("删除资源失败err:", err)
			return
		}
	}
	return
}

/*func (this *UserRanking) Run() {
	defer func() {
		if err := recover(); err != nil {
			logx.Error(fmt.Errorf("UserRanking err = %s", err))
		}
	}()
	locker := locker.NewRedisLocker(this.svcCtx.Redis, model.UserRanking)
	defer locker.Unlock()
	if err := locker.Lock(); err != nil {
		logx.Error(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	userList, err := this.svcCtx.Rpc.FindUser(ctx, &pb.QueryUserReq{
		IsDisable: 2,
	})
	if err != nil {
		logx.Error("获取所有用户失败err:", err)
		return
	}
	ids := make([]int64, 0, len(userList.Users))
	for _, v := range userList.Users {
		ids = append(ids, v.Uid)
	}
	points, err := this.svcCtx.Rpc.FindUserPoints(ctx, &pb.FindUserPointsReq{Uid: ids})
	if err != nil {
		logx.Error("获取用户积分失败err:", err)
		return
	}
	for _, user := range points.List {
		this.svcCtx.Redis.ZaddCtx(ctx, model.UserIntegral, user.Points, strconv.FormatInt(user.Uid, 10))
	}
	return
}*/
