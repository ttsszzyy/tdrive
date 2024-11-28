package parameter

import (
	"T-driver/app/user/rpc/pb"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateParameterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateParameterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateParameterLogic {
	return &UpdateParameterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateParameterLogic) UpdateParameter(req *types.UpdateParameterReq) (resp *types.Response, err error) {
	resp = &types.Response{}
	p := &pb.SaveDictReq{Dict: make([]*pb.Dict, 0, len(req.Parms))}
	for _, v := range req.Parms {
		p.Dict = append(p.Dict, &pb.Dict{
			Id:          v.Id,
			ParamType:   v.ParamType,
			Name:        v.Name,
			Desc:        v.Desc,
			Value:       v.Value,
			BackupValue: v.BackupValue,
		})
	}
	//修改数字字典
	_, err = l.svcCtx.Rpc.SaveDict(l.ctx, p)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
