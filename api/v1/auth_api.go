package apiv1

import (
	// "net/http" // No longer needed here

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Signup(c *gin.Context)
	GoogleCallback(c *gin.Context)
	GoogleLogin(c *gin.Context)
}

func SetupAuthRoutes(router *gin.RouterGroup, authHandler AuthHandler, middlewares ...gin.HandlerFunc) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.Signup)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/google/login", authHandler.GoogleLogin)
		authGroup.GET("/google/callback", authHandler.GoogleCallback)
	}
}
