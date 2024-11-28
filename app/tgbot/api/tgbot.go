package main

import (
	"T-driver/app/tgbot/api/internal/config"
	"T-driver/app/tgbot/api/internal/handler"
	tgbot "T-driver/app/tgbot/api/internal/handler/bot"
	"T-driver/app/tgbot/api/internal/svc"
	"context"
	"flag"
	"fmt"
	"github.com/go-telegram/bot"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/tgbot-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	//启动机器人
	opts := []bot.Option{
		bot.WithDefaultHandler(tgbot.RegisterBotProcessHandler(ctx)),
	}
	b, err := bot.New(ctx.Config.Telegram.Url, ctx.Config.Telegram.Token, opts...)
	if nil != err {
		logx.Must(fmt.Errorf("链接bot失败！%s", err))
	}
	ctx.TgBot = b
	ctx.TgBot.SetWebhook(context.Background(), &bot.SetWebhookParams{
		URL: c.Telegram.WebhookUrl,
	})
	go ctx.TgBot.StartWebhook(context.Background())
	handler.RegisterHandlers(server, ctx)

	/*notifyContext, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(tgbot.RegisterBotProcessHandler(ctx)),
	}
	b, err := bot.New(ctx.Config.Telegram.Token, opts...)
	if nil != err {
		logx.Error("链接bot失败！")
		return
	}
	go b.Start(notifyContext)*/

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
