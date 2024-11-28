package user

import (
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

// UpdateReadStatusLogic 修改用户为已阅读注意事项
type UpdateReadStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateReadStatusLogic 新建 修改用户为已阅读注意事项
func NewUpdateReadStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReadStatusLogic {
	return &UpdateReadStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateReadStatus 实现 修改用户为已阅读注意事项
func (l *UpdateReadStatusLogic) UpdateReadStatus() error {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}

	user, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		return errors.DbError()
	}
	_, err = l.svcCtx.Rpc.SaveUser(l.ctx, &pb.User{Id: user.Id, IsRead: 1})
	if err != nil {
		logx.Error(err)
		return errors.DbError()
	}

	return nil
}
