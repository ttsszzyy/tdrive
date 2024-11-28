package parameter

import (
	"T-driver/app/user/rpc/pb"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryParameterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryParameterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryParameterLogic {
	return &QueryParameterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryParameterLogic) QueryParameter(req *types.QueryParameterReq) (resp *types.QueryParameterResp, err error) {
	list, err := l.svcCtx.Rpc.FindDictByName(l.ctx, &pb.FindDictByNameReq{ParamType: []int64{req.ParamType}})
	if err != nil {
		return nil, err
	}
	resp = &types.QueryParameterResp{DictList: make([]*types.Dict, 0, len(list.DictList))}
	for _, v := range list.DictList {
		resp.DictList = append(resp.DictList, &types.Dict{
			Id:          v.Id,
			ParamType:   v.ParamType,
			Name:        v.Name,
			Desc:        v.Desc,
			Value:       v.Value,
			BackupValue: v.BackupValue,
		})
	}
	return resp, nil
}
