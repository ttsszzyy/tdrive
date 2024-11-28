package logic

import (
	"T-driver/app/user/model"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *pb.QueryUserReq) (*pb.UserList, error) {
	sb := squirrel.Select().Where("deleted_time =?", 0)
	if in.IsDisable > 0 {
		sb = sb.Where(squirrel.Eq{"is_disable": in.IsDisable})
	}
	if in.Id > 0 {
		sb = sb.Where(squirrel.Eq{"id": in.Id})
	}
	if in.Pid > 0 {
		sb = sb.Where(squirrel.Eq{"pid": in.Pid})
	}
	if in.StartTime > 0 {
		sb = sb.Where(squirrel.GtOrEq{"created_time": in.StartTime})
	}
	if in.EndTime > 0 {
		sb = sb.Where(squirrel.LtOrEq{"created_time": in.EndTime})
	}
	if in.Name != "" {
		sb = sb.Where(squirrel.Like{"name": fmt.Sprintf("%s%%", in.Name)})
	}
	if len(in.Uid) > 0 {
		sb = sb.Where(squirrel.Eq{"uid": in.Uid})
	}
	var total int64
	var list = make([]*model.User, 0)
	var err error
	if in.Page > 0 && in.Size > 0 {
		list, total, err = l.svcCtx.UserModel.ListPage(l.ctx, in.Page, in.Size, sb.OrderBy("created_time desc"))
		if err != nil {
			return nil, err
		}
	} else {
		list, err = l.svcCtx.UserModel.List(l.ctx, sb.OrderBy("created_time desc"))
		if err != nil {
			return nil, err
		}
		total = int64(len(list))
	}

	resp := &pb.UserList{
		Total: total,
		Users: make([]*pb.User, 0, len(list)),
	}
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

	return resp, nil
}
