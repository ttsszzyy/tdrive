package logic

import (
	"context"
	"time"

	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SaveUserTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveUserTokenLogic {
	return &SaveUserTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveUserTokenLogic) SaveUserToken(in *pb.SaveUserTokenReq) (*pb.Response, error) {
	err := l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, s sqlx.Session) error {
		storage, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(ctx, in.Uid, 0)
		if err != nil {
			return err
		}
		token, err := l.svcCtx.UserTokenModel.FindOneByUidDeletedTime(ctx, in.Uid, 0)
		if err != nil {
			return err
		}
		if in.Token > storage.Storage-token.Token {
			in.Token = storage.Storage - token.Token
		}

		if in.Token <= 0 {
			return nil
		}

		_, err = l.svcCtx.UserTokenExchangeModel.Insert(ctx, &model.UserTokenExchange{
			Uid:           in.Uid,
			ExchangeToken: in.Token,
			CreatedTime:   time.Now().Unix(),
		}, s)
		if err != nil {
			return err
		}
		token.Token += in.Token
		err = l.svcCtx.UserTokenModel.Update(ctx, token, s)
		if err != nil {
			return err
		}

		return nil
	})

	return &pb.Response{}, err
}
