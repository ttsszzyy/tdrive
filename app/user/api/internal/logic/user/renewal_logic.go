package user

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	storage "github.com/utopiosphe/titan-storage-sdk"
	"github.com/zeromicro/go-zero/core/logx"
)

type RenewalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRenewalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RenewalLogic {
	return &RenewalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RenewalLogic) Renewal(req *types.RenewalReq) (resp *types.RenewalResp, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}

	// res, err := l.svcCtx.Tenant.RefreshToken(l.ctx, req.Token)
	// if err != nil {
	// 	return nil, err
	// }
	u, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		logx.Error("find user by uid error", err)
		return nil, err
	}
	res, err := l.svcCtx.Tenant.SSOLogin(l.ctx, storage.SubUserInfo{
		EntryUUID: strconv.FormatInt(u.Uid, 10),
		Username:  u.Name,
		Avatar:    u.Avatar,
	})
	if err != nil {
		logx.Errorf("login titan error:%v name:%v uid:%v", err, u.Name, u.Uid)
		return nil, err
	}
	_, err = l.svcCtx.Rpc.UpdateUserTitanToken(l.ctx, &pb.UpdateUserTitanTokenReq{
		Uid:    userData.User.ID,
		Token:  res.Token,
		Expire: res.Exp,
	})
	if err != nil {
		return nil, err
	}
	one := &pb.UserTitanToken{}
	one.Token = res.Token
	one.Expire = res.Exp
	one.Uid = userData.User.ID
	marshal, _ := json.Marshal(one)
	//更新redis
	err = l.svcCtx.Redis.SetCtx(l.ctx, fmt.Sprintf(model.GetStorageUser+"%v", one.Uid), string(marshal))
	if err != nil {
		return nil, err
	}
	return &types.RenewalResp{TitanToken: res.Token, Exp: res.Exp}, nil
}
