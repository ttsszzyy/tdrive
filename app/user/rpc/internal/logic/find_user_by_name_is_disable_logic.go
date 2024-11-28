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

type FindUserByNameIsDisableLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserByNameIsDisableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserByNameIsDisableLogic {
	return &FindUserByNameIsDisableLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户信息
func (l *FindUserByNameIsDisableLogic) FindUserByNameIsDisable(in *pb.QueryUserReq) (*pb.QueryUserRes, error) {
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

	res := &pb.QueryUserRes{
		Total: total,
		Users: make([]*pb.UserInfo, 0, len(list)),
	}
	for _, user := range list {
		/*//剩余积分
		integral, err := l.svcCtx.UserPointsModel.CountPoints(l.ctx, squirrel.Select().Where("uid = ?", user.Uid))
		if err != nil {
			return nil, err
		}
		//总空间
		where := squirrel.Select().Where(squirrel.Gt{"storage": 0}).Where("uid = ?", user.Uid)
		storage, err := l.svcCtx.UserStorageModel.CountStorage(l.ctx, where)
		if err != nil {
			return nil, err
		}
		//已使用空间
		sbw := squirrel.Select().Where(squirrel.Lt{"storage": 0}).Where("uid = ?", user.Uid)
		storageUse, err := l.svcCtx.UserStorageModel.CountStorage(l.ctx, sbw)
		if err != nil {
			return nil, err
		}*/

		info := &pb.UserInfo{
			Id:            user.Id,
			Name:          user.Name,
			Avatar:        user.Avatar,
			Pid:           user.Pid,
			IsDisable:     user.IsDisable,
			Mail:          user.Mail,
			WalletAddress: user.WalletAddress,
			CreatedTime:   user.CreatedTime,
			UpdatedTime:   user.UpdatedTime,
			Source:        user.Source,
			Uid:           user.Uid,
		}
		if user.Pid != 0 {
			one, err := l.svcCtx.UserModel.FindOne(l.ctx, user.Pid)
			if err != nil {
				if err != model.ErrNotFound {
					return nil, err
				}
			} else {
				info.Pname = one.Name
			}
		}
		res.Users = append(res.Users, info)
	}

	return res, nil
}
