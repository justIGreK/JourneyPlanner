package main

import (
	"JourneyPlanner/cmd/config"
	"JourneyPlanner/cmd/handler"
	"JourneyPlanner/cmd/handler/ws"
	mongorepo "JourneyPlanner/internal/repository/mongo"
	"JourneyPlanner/internal/service"
	"JourneyPlanner/internal/service/chat"
	logger "JourneyPlanner/pkg/log"
	"context"
	"net/http"

	"go.uber.org/zap"
)

var logs *zap.Logger

// @title Journer Planner
// @description Application for planning your journey

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()
	logs := logger.GetLogger()
	setUpProjectLogger(logs)
	config.LoadEnv()
	defer logs.Sync()
	dbclient := mongorepo.CreateMongoClient(ctx)
	userRepo := mongorepo.NewMongoUserRepo(dbclient)
	taskRepo := mongorepo.NewMongoTaskRepo(dbclient)
	pollRepo := mongorepo.NewMongoPollRepo(dbclient)
	groupRepo := mongorepo.NewMongoGroupRepo(dbclient)
	inviteRepo := mongorepo.NewMongoInviteRepo(dbclient)
	blacklistRepo := mongorepo.NewMongoBlacklistRepo(dbclient)
	chatRepo := mongorepo.NewChatRepository(dbclient)
	chatService := chat.NewChatService(chatRepo)

	userSrv := service.NewUserSrv(userRepo)
	pollSrv := service.NewPollSrv(pollRepo, groupRepo)
	taskSrv := service.NewTaskSrv(taskRepo, groupRepo)
	groupSrv := service.NewGroupSrv(groupRepo, userRepo, inviteRepo, blacklistRepo)
	wsHandler := ws.NewWebSocketHandler(chatService, groupSrv)
	groupSrv.NotifyUserDisconnect = wsHandler.NotifyUserDisconnect

	handler := handler.NewHandler(pollSrv, taskSrv, userSrv, groupSrv)
	logs.Sugar().Info("Server is now listening 8080...")
	err := http.ListenAndServe(":8080", handler.InitRoutes(wsHandler))
	if err != nil {
		logs.Sugar().Fatal("Server error", zap.Error(err))
	}
}

func setUpProjectLogger(logger *zap.Logger) {
	config.SetLogger(logger)
	handler.SetLogger(logger)
	service.SetLogger(logger)
	mongorepo.SetLogger(logger)
	ws.SetLogger(logger)
}
