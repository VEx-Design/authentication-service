package main

import (
	"log"
	handler "project-management-service/external/handler/adaptors/gin/api"
	"project-management-service/external/handler/adaptors/gin/router"
	gorm "project-management-service/external/repository/adaptors/postgres"
	repository "project-management-service/external/repository/adaptors/postgres/controller"
	"project-management-service/internal/core/service"
	"project-management-service/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	postgresDB := db.ConnectToPG()
	client := postgresDB.GetClient()

	gorm.SyncDB(client)

	userRepo := repository.NewUserRepositoryPQ(client)
	userSrv := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userSrv)

	r := gin.Default()
	router.AuthRoutes(r, *authHandler)
	r.Run(":3000")
}
