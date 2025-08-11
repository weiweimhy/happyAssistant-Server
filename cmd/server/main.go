package main

import (
	"happyAssistant/internal/config"
	"happyAssistant/internal/controller"
	"happyAssistant/internal/initialize"
	"happyAssistant/internal/logger"
	"happyAssistant/pkg/wshub"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartWebsocketServer() {
	hub := wshub.GetInstance()
	hub.SetClientOptions(
		wshub.WithReadDeadline(45*time.Second),
		wshub.WithSupportPing(20*time.Second),
	)

	// 创建协议控制器
	protocolController := controller.NewProtocolController()

	hub.OnOpen = func(client wshub.IClient) {
		log.Info("Client connected")
	}

	hub.OnClose = func(client wshub.IClient) {
		log.Info("Client disconnected")
	}

	hub.OnMessage = func(client wshub.IClient, msg []byte) {
		// 将消息路由到协议控制器处理
		protocolController.HandleMessage(client, msg)
	}

	hub.OnError = func(client wshub.IClient, err error) {
		log.Errorf("WebSocket error: %v", err)
	}

	err := hub.Start(config.Cfg.Server.Route, config.Cfg.Server.Port)

	if err != nil {
		log.Fatalln("Websocket server start error:", err)
	}
}

func main() {
	config.LoadConfig("configs/config_debug.yaml")
	logger.InitLogger(config.Cfg.Log)
	initialize.InitMongoDBClient(config.Cfg.MongoDB)
	StartWebsocketServer()
}
