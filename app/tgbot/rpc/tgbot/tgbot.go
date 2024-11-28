// Code generated by goctl. DO NOT EDIT.
// Source: tgbot.proto

package tgbot

import (
	"context"

	"T-driver/app/tgbot/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	BotCommand         = pb.BotCommand
	DelAllPinMsgReq    = pb.DelAllPinMsgReq
	DelPinMsgReq       = pb.DelPinMsgReq
	DelPinMsgResp      = pb.DelPinMsgResp
	Item               = pb.Item
	Markup             = pb.Markup
	SendBotCommandReq  = pb.SendBotCommandReq
	SendBotMsgRequest  = pb.SendBotMsgRequest
	SendBotMsgResponse = pb.SendBotMsgResponse
	SendPhotoRequest   = pb.SendPhotoRequest

	Tgbot interface {
		// 文本发送
		SendBotMsg(ctx context.Context, in *SendBotMsgRequest, opts ...grpc.CallOption) (*SendBotMsgResponse, error)
		// 发送pin消息
		SendPhoto(ctx context.Context, in *SendPhotoRequest, opts ...grpc.CallOption) (*SendBotMsgResponse, error)
		DelPinMsg(ctx context.Context, in *DelPinMsgReq, opts ...grpc.CallOption) (*DelPinMsgResp, error)
		DelAllPinMsg(ctx context.Context, in *DelAllPinMsgReq, opts ...grpc.CallOption) (*DelPinMsgResp, error)
		SendBotCommand(ctx context.Context, in *SendBotCommandReq, opts ...grpc.CallOption) (*DelPinMsgResp, error)
	}

	defaultTgbot struct {
		cli zrpc.Client
	}
)

func NewTgbot(cli zrpc.Client) Tgbot {
	return &defaultTgbot{
		cli: cli,
	}
}

// 文本发送
func (m *defaultTgbot) SendBotMsg(ctx context.Context, in *SendBotMsgRequest, opts ...grpc.CallOption) (*SendBotMsgResponse, error) {
	client := pb.NewTgbotClient(m.cli.Conn())
	return client.SendBotMsg(ctx, in, opts...)
}

// 发送pin消息
func (m *defaultTgbot) SendPhoto(ctx context.Context, in *SendPhotoRequest, opts ...grpc.CallOption) (*SendBotMsgResponse, error) {
	client := pb.NewTgbotClient(m.cli.Conn())
	return client.SendPhoto(ctx, in, opts...)
}

func (m *defaultTgbot) DelPinMsg(ctx context.Context, in *DelPinMsgReq, opts ...grpc.CallOption) (*DelPinMsgResp, error) {
	client := pb.NewTgbotClient(m.cli.Conn())
	return client.DelPinMsg(ctx, in, opts...)
}

func (m *defaultTgbot) DelAllPinMsg(ctx context.Context, in *DelAllPinMsgReq, opts ...grpc.CallOption) (*DelPinMsgResp, error) {
	client := pb.NewTgbotClient(m.cli.Conn())
	return client.DelAllPinMsg(ctx, in, opts...)
}

func (m *defaultTgbot) SendBotCommand(ctx context.Context, in *SendBotCommandReq, opts ...grpc.CallOption) (*DelPinMsgResp, error) {
	client := pb.NewTgbotClient(m.cli.Conn())
	return client.SendBotCommand(ctx, in, opts...)
}