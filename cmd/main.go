package main

import (
	handler "authentication-service/external/handler/adaptors/gin/api"
	"authentication-service/external/handler/adaptors/gin/router"
	gorm "authentication-service/external/repository/adaptors/postgres"
	repository "authentication-service/external/repository/adaptors/postgres/controller"
	"authentication-service/internal/core/service"
	"authentication-service/pkg/db"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, relying on environment variables")
	}

	postgresDB := db.ConnectToPG()
	client := postgresDB.GetClient()

	gorm.SyncDB(client)

	userRepo := repository.NewUserRepositoryPQ(client)
	userSrv := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userSrv)

	r := gin.Default()
	router.AuthRoutes(r, *authHandler)
	r.Run(":6740")
}
