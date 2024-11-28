package logic

import (
	"T-driver/app/user/model"
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveMessageLogic {
	return &SaveMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送消息
func (l *SaveMessageLogic) SaveMessage(in *pb.SaveMessageReq) (*pb.Response, error) {
	if in.Id > 0 {
		one, err := l.svcCtx.MessageModel.FindOne(l.ctx, in.Id)
		if err != nil {
			return nil, err
		}
		if in.Status > 0 {
			one.Status = in.Status
		}
		if in.Remark != "" {
			one.Remark = in.Remark
		}
		one.UpdatedTime = time.Now().Unix()
		err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			err = l.svcCtx.MessageModel.Update(l.ctx, one, session)
			if err != nil {
				return err
			}
			if in.Status == 2 {
				user, err := l.svcCtx.UserModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
				if err != nil {
					return err
				}
				user.Puid = one.Puid
				err = l.svcCtx.UserModel.Update(l.ctx, user)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		_, err := l.svcCtx.MessageModel.Insert(l.ctx, &model.Message{
			Uid:         in.Uid,
			Name:        in.Name,
			Puid:        in.Puid,
			Status:      in.Status,
			Remark:      in.Remark,
			CreatedTime: time.Now().Unix(),
		})
		if err != nil {
			return nil, err
		}
	}

	return &pb.Response{}, nil
}
