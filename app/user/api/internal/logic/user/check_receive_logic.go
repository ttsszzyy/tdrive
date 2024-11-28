package user

import (
	"context"
	"fmt"
	"net"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"T-driver/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

// CheckReceiveLogic 校验是否有领取活动积分的资格
type CheckReceiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCheckReceiveLogic 新建 校验是否有领取活动积分的资格
func NewCheckReceiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckReceiveLogic {
	return &CheckReceiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CheckReceive 实现 校验是否有领取活动积分的资格
func (l *CheckReceiveLogic) CheckReceive(req *types.ReceivePointsReq) (resp *types.MsgResponse, err error) {
	var (
		lan = fmt.Sprintf("%s", l.ctx.Value("language"))
		ip  = fmt.Sprintf("%s", l.ctx.Value("ip"))
	)
	resp = &types.MsgResponse{Status: false}

	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}

	// 解密经纬度，并判断区间是否符合
	lon, lat, err := getLocation(req.Lon, req.Lat, l.svcCtx.PrvKey)
	if err != nil {
		logx.Error(err)
		return nil, errors.ParamsError(lan)
	}

	if ip == "" {
		resp.Msg = errors.GetError(errors.ErrGetRequestIP, lan).Msg()
		return resp, nil
	}
	ipv4 := net.ParseIP(ip)
	if ipv4.IsPrivate() || ipv4.IsLoopback() {
		resp.Msg = errors.GetError(errors.ErrGetRequestIP, lan).Msg()
		return resp, nil
	}

	distance := utils.HaversineDistance(lat, lon, l.svcCtx.Config.Action.Lat, l.svcCtx.Config.Action.Lon)
	logx.Errorf("lat:%v lon:%v distance:%v", lat, lon, distance)
	if distance > l.svcCtx.Config.Action.Distance {
		resp.Msg = errors.GetError(errors.ErrNotActionDistance, lan).Msg()
		return resp, nil
	}

	_, err = l.svcCtx.Rpc.CheckReceive(l.ctx, &pb.ReceiveActionPointsReq{Uid: userData.User.ID, Ip: ip, Name: l.svcCtx.Config.Action.Name, Lan: lan})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			resp.Msg = s.Message()
		} else {
			resp.Msg = errors.GetError(0, lan).Msg()
		}

		return resp, nil
	}

	resp.Status = true
	return resp, nil
}
