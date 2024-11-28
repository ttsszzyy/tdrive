package share

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"math"
	"time"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareListLogic {
	return &ShareListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareListLogic) ShareList(req *types.QueryShareReq) (resp *types.QueryShareResp, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	if req.Size > 100 {
		req.Size = 100
	}
	list, err := l.svcCtx.Rpc.FindShare(l.ctx, &pb.FindShareReq{
		Uid:  userData.User.ID,
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.QueryShareResp{
		Total:  list.Total,
		Shares: make([]*types.Share, 0),
	}
	for _, v := range list.Shares {
		status := 0
		day := 0
		//获取到期时间
		switch {
		case v.EffectiveTime == 0:
			status = 0
		case time.Now().Unix() > v.EffectiveTime:
			status = 1
		case time.Now().Unix() < v.EffectiveTime:
			status = 2
			i := v.EffectiveTime - time.Now().Unix()
			day = int(math.Ceil(float64(i) / 86400))
		}
		isPass := false
		if v.Password != "" {
			isPass = true
		}
		resp.Shares = append(resp.Shares, &types.Share{
			Id:          v.Id,
			UId:         v.Uid,
			AssetIds:    v.AssetIds,
			AssetName:   v.AssetName,
			AssetSize:   v.AssetSize,
			AssetType:   v.AssetType,
			Link:        fmt.Sprintf("%s/api/v1/user/share/resource/%s", l.svcCtx.Config.BaseUrl, v.Uuid),
			ReadNum:     v.ReadNum,
			SaveNum:     v.SaveNum,
			Status:      int64(status),
			Day:         day,
			CreatedTime: v.CreatedTime,
			IsPass:      isPass,
		})
	}

	return
}
