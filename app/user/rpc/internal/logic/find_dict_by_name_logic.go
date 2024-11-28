package logic

import (
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindDictByNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindDictByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindDictByNameLogic {
	return &FindDictByNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询任务池
func (l *FindDictByNameLogic) FindDictByName(in *pb.FindDictByNameReq) (*pb.FindDictByNameResp, error) {
	sb := squirrel.Select().Where("deleted_time = ?", 0)
	if len(in.ParamType) > 0 {
		sb = sb.Where(squirrel.Eq{"param_type": in.ParamType})
	}
	if len(in.Code) > 0 {
		sb = sb.Where(squirrel.Eq{"code": in.Code})
	}
	list, err := l.svcCtx.DictModel.List(l.ctx, sb)
	if err != nil {
		return nil, err
	}
	resp := &pb.FindDictByNameResp{DictList: make([]*pb.Dict, 0, len(list))}
	for _, v := range list {
		resp.DictList = append(resp.DictList, &pb.Dict{
			Id:          v.Id,
			ParamType:   v.ParamType,
			Name:        v.Name,
			Desc:        v.Desc,
			Value:       v.Value,
			BackupValue: v.BackupValue,
			Code:        v.Code,
		})
	}
	return resp, nil
}
