package logic

import (
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDictLogic {
	return &SaveDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存任务池
func (l *SaveDictLogic) SaveDict(in *pb.SaveDictReq) (*pb.Response, error) {
	for _, dict := range in.Dict {
		one, err := l.svcCtx.DictModel.FindOne(l.ctx, dict.Id)
		if err != nil {
			return nil, err
		}
		one.ParamType = dict.ParamType
		one.Name = dict.Name
		one.Desc = dict.Desc
		one.Value = dict.Value
		one.BackupValue = dict.BackupValue
		err = l.svcCtx.DictModel.Update(l.ctx, one)
		if err != nil {
			return nil, err
		}
	}

	return &pb.Response{}, nil
}
