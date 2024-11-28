package logic

import (
	"context"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateShareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateShareLogic {
	return &UpdateShareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改分享任务密码和有效期
func (l *UpdateShareLogic) UpdateShare(in *pb.UpdateShareReq) (*pb.Response, error) {
	one, err := l.svcCtx.ShareModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Password != "" {
		one.Password = in.Password
	}
	if in.SaveNum > 0 {
		one.SaveNum = one.SaveNum + in.SaveNum
	}
	if in.ReadNum > 0 {
		one.ReadNum = one.ReadNum + in.ReadNum
	}
	if in.EffectiveTime > 0 {
		one.EffectiveTime = in.EffectiveTime
	}
	one.UpdatedTime = time.Now().Unix()
	err = l.svcCtx.ShareModel.Update(l.ctx, one)
	if err != nil {
		return nil, err
	}
	return &pb.Response{}, nil
}
