package user

import (
	"T-driver/app/user/rpc/pb"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserLogic {
	return &QueryUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryUserLogic) QueryUser(req *types.QueryUserReq) (resp *types.QueryUserRes, err error) {
	time, err := l.svcCtx.Rpc.FindUserByNameIsDisable(l.ctx, &pb.QueryUserReq{
		Name:      req.Name,
		IsDisable: req.IsDisable,
		Page:      req.Page,
		Size:      req.Size,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.QueryUserRes{
		Total: time.Total,
		Users: make([]*types.User, 0, len(time.Users)),
	}
	for _, user := range time.Users {
		resp.Users = append(resp.Users, &types.User{
			Id:            user.Id,
			Name:          user.Name,
			Avatar:        user.Avatar,
			Pid:           user.Pid,
			IsDisable:     user.IsDisable,
			Mail:          user.Mail,
			WalletAddress: user.WalletAddress,
			CreatedTime:   user.CreatedTime,
			UpdatedTime:   user.UpdatedTime,
			Pname:         user.Name,
		})
	}

	return resp, nil
}
