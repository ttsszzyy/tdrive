package auth

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"
	"T-driver/common/lib/jwt"
	"T-driver/common/utils"
	"context"
	"fmt"
	"time"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLoginLogic {
	return &AdminLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLoginLogic) AdminLogin(req *types.AdminLoginReq) (resp *types.LoginRes, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	admin, err := l.svcCtx.Rpc.FindOneByAccountDeletedTime(l.ctx, &pb.AccountReq{Account: req.Username})
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.CustomError("用户名不存在")
		}
		l.Error(err)
		return nil, errors.DbError(lan)
	}
	if admin.IsDisable != model.State_Enable {
		return nil, errors.CustomError("用户角色已被禁用, 无法登陆")
	}
	if admin.Password != "" {
		if !utils.CheckPassword(req.Password, admin.Password, model.PassSalt) {
			l.Error(err)
			return nil, errors.CustomError("用户名或密码错误")
		}
	}
	payload := jwt.Payload{
		"id":      admin.Id,
		"account": admin.Account,
	}
	token, err := l.svcCtx.Jwt.Token(payload)
	if err != nil {
		l.Error(err)
		return nil, errors.SystemError(lan)
	}

	/*if err := l.svcCtx.Jwt.Store(l.ctx, strconv.FormatInt(admin.Id, 10), token); err != nil {
		l.Error(err)
		return nil, errors.SystemError()
	}*/
	//更新用户登陆时间
	l.svcCtx.Rpc.SaveAdmin(l.ctx, &pb.Admin{
		Id:        admin.Id,
		Account:   admin.Account,
		Password:  admin.Password,
		Avatar:    admin.Avatar,
		IsDisable: admin.IsDisable,
		LastTime:  time.Now().Unix(),
		Remark:    admin.Remark,
	})

	return &types.LoginRes{
		Token: token,
	}, nil
}
