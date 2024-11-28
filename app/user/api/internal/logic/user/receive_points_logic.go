package user

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"T-driver/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

// ReceivePointsLogic 领取活动积分
type ReceivePointsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewReceivePointsLogic 新建 领取活动积分
func NewReceivePointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceivePointsLogic {
	return &ReceivePointsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ReceivePoints 实现 领取活动积分
func (l *ReceivePointsLogic) ReceivePoints(req *types.ReceivePointsReq) (resp *types.MsgResponse, err error) {
	lan, ok := l.ctx.Value("language").(string)
	ip, ok := l.ctx.Value("ip").(string)
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
	distance := utils.HaversineDistance(lat, lon, l.svcCtx.Config.Action.Lat, l.svcCtx.Config.Action.Lon)
	logx.Errorf("lon:%v lat:%v distance:%v", lon, lat, distance)
	if distance > l.svcCtx.Config.Action.Distance {
		resp.Msg = errors.GetError(errors.ErrNotActionDistance, lan).Msg()
		return resp, nil
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

	_, err = l.svcCtx.Rpc.ReceiveActionPoints(l.ctx, &pb.ReceiveActionPointsReq{Uid: userData.User.ID, Location: fmt.Sprintf("%v %v", lat, lon),
		Ip: ip, Name: l.svcCtx.Config.Action.Name, Points: l.svcCtx.Config.Action.Points, Lan: lan})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			resp.Msg = s.Message()
		} else {
			resp.Msg = errors.GetError(0, lan).Msg()
		}

		return resp, err
	}

	resp.Status = true
	return resp, nil
}

// 使用 RSA 私钥解密经纬度
func getLocation(lon, lat string, privateKey *rsa.PrivateKey) (float64, float64, error) {
	var lonf, latf float64

	lonbytes, err := base64.StdEncoding.DecodeString(lon)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode base64(lon) data:%w", err)
	}
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, lonbytes)
	if err != nil {
		return 0, 0, fmt.Errorf("decryption lon failed: %v", err)
	}
	lonf, err = strconv.ParseFloat(string(decryptedData), 10)
	if err != nil {
		return 0, 0, fmt.Errorf("parse lon failed: %v", err)
	}
	latbytes, err := base64.StdEncoding.DecodeString(lat)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode base64(lat) data:%w", err)
	}
	decryptedData, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, latbytes)
	if err != nil {
		return 0, 0, fmt.Errorf("decryption lat failed: %v", err)
	}
	latf, err = strconv.ParseFloat(string(decryptedData), 10)
	if err != nil {
		return 0, 0, fmt.Errorf("parse lat failed: %v", err)
	}
	return lonf, latf, nil
}
