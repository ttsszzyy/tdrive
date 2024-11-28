package task

import (
	"T-driver/app/video/api/internal/svc"
	"T-driver/app/video/model"
	"T-driver/app/video/rpc/videoPb"
	"T-driver/common/locker"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

var C *cron.Cron

func InitTask(svcCtx *svc.ServiceContext) {
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	C = cron.New(cron.WithParser(parser))
	//短视频标签数量
	C.AddJob("0 * */1 * * *", cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger)).Then(&VideoLabelCount{svcCtx: svcCtx}))
	C.Start()
}

type (
	VideoLabelCount struct {
		svcCtx *svc.ServiceContext
	}
)

func (this *VideoLabelCount) Run() {
	defer func() {
		if err := recover(); err != nil {
			logx.Error(fmt.Errorf("UserRanking err = %s", err))
		}
	}()
	locker := locker.NewRedisLocker(this.svcCtx.Redis, model.VideoLabelCount)
	defer locker.Unlock()
	if err := locker.Lock(); err != nil {
		logx.Error(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	list, err := this.svcCtx.VideoRpc.GetLabelList(ctx, &videoPb.GetLabelListReq{})
	if err != nil {
		logx.Error("获取所有短视频失败err:", err)
		return
	}
	ids := make([]string, 0, len(list.List))
	for _, v := range list.List {
		ids = append(ids, v.Id)
	}
	//获取标签总数
	labelUsers, err := this.svcCtx.VideoRpc.CountLabelUser(ctx, &videoPb.CountLabelUserReq{
		Lid: ids,
	})
	if err != nil {
		logx.Error("获取标签总数失败err:", err)
		return
	}
	for _, v := range labelUsers.List {
		//更新标签数量
		err = this.svcCtx.Redis.SetCtx(ctx, fmt.Sprintf("%s%v", model.LabelLikeKey, v.Lid), strconv.FormatInt(v.LikesNum, 10))
		if err != nil {
			logx.Error("更新标签点赞数量失败err:" + err.Error())
			return
		}
		err = this.svcCtx.Redis.SetCtx(ctx, fmt.Sprintf("%s%v", model.LabelNoLikeKey, v.Lid), strconv.FormatInt(v.NoLikesNum, 10))
		if err != nil {
			logx.Error("更新标签点踩数量失败err:" + err.Error())
			return
		}
	}
	return
}
