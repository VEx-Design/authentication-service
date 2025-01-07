package router

import (
	handler "authentication-service/external/handler/adaptors/gin/api"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authHandler handler.AuthHandler) {
	router.GET("/auth/google", authHandler.GoogleLogin)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)
}
