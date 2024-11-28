package logic

import (
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"
	"github.com/Masterminds/squirrel"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserByPidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserByPidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserByPidLogic {
	return &FindUserByPidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询推荐用户
func (l *FindUserByPidLogic) FindUserByPid(in *pb.FindUserByPidReq) (*pb.UserList, error) {
	sb := squirrel.Select().Where("deleted_time = ?", 0)
	if in.Pid > 0 {
		sb = sb.Where("pid = ?", in.Pid)
	}
	if in.Puid > 0 {
		sb = sb.Where("puid = ?", in.Puid)
	}
	list, err := l.svcCtx.UserModel.List(l.ctx, sb)
	if err != nil {
		return nil, err
	}
	resp := &pb.UserList{Users: make([]*pb.User, 0, len(list))}
	for _, user := range list {
		resp.Users = append(resp.Users, &pb.User{
			Id:            user.Id,
			Uid:           user.Uid,
			Name:          user.Name,
			Avatar:        user.Avatar,
			Mail:          user.Mail,
			WalletAddress: user.WalletAddress,
			Source:        user.Source,
			RecommendCode: user.RecommendCode,
			Distribution:  user.Distribution,
			Pid:           user.Pid,
			IsDisable:     user.IsDisable,
			CreatedTime:   user.CreatedTime,
			UpdatedTime:   user.UpdatedTime,
			Puid:          user.Puid,
			IsReceive:     user.IsReceive,
			LanguageCode:  user.LanguageCode,
		})
	}
	resp.Total = int64(len(list))
	return resp, nil
}
