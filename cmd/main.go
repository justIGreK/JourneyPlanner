package main

import (
	"JourneyPlanner/cmd/config"
	"JourneyPlanner/cmd/handler"
	mongorepo "JourneyPlanner/internal/repository/mongo"
	"JourneyPlanner/internal/service"
	logger "JourneyPlanner/pkg/log"
	"context"
	"net/http"

	"go.uber.org/zap"
)
// @title Journer Planner
// @description Application for planning your journey

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadEnv()
	logger := logger.GetLogger()
	ctx := context.Background()
	dbclient := mongorepo.CreateMongoClient(ctx)
	userRepo := mongorepo.NewUserTaskRepo(dbclient)
	taskRepo := mongorepo.NewUserTaskRepo(dbclient)
	pollRepo := mongorepo.NewMongoPollRepo(dbclient)

	userSrv := service.NewUserSrv(userRepo)
	taskSrv := service.NewTaskSrv(taskRepo)
	pollSrv := service.NewPollSrv(pollRepo)

	handler := handler.NewHandler(pollSrv, taskSrv, userSrv)
	logger.Info("Server is now listening 8080...")
	err := http.ListenAndServe(":8080", handler.InitRoutes())
	if err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}
