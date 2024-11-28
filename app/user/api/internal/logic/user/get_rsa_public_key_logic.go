package user

import (
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/common/db"
	"T-driver/common/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetRsaPublicKeyLogic 获取rsa公钥
type GetRsaPublicKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetRsaPublicKeyLogic 新建 获取rsa公钥
func NewGetRsaPublicKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRsaPublicKeyLogic {
	return &GetRsaPublicKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetRsaPublicKey 实现 获取rsa公钥
func (l *GetRsaPublicKeyLogic) GetRsaPublicKey() (resp string, err error) {
	var (
		lan = fmt.Sprintf("%s", l.ctx.Value("language"))
	)

	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return "", errors.UnauthError(lan)
	}

	if l.svcCtx.PubKey == "" {
		return "", errors.GetError(errors.ErrGetRSAPublickey, lan)
	}

	return l.svcCtx.PubKey, nil
}
