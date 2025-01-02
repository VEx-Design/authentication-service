package router

import (
	handler "project-management-service/external/handler/adaptors/gin/api"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authHandler handler.AuthHandler) {
	router.GET("/auth", authHandler.GoogleLogin)
	router.GET("/auth/callback/google", authHandler.GoogleCallback)
}
