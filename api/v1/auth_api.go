package apiv1

import (
	// "net/http" // No longer needed here

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Signup(c *gin.Context)
}

func SetupAuthRoutes(router *gin.RouterGroup, authHandler AuthHandler) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.Signup)
		authGroup.POST("/login", authHandler.Login)
	}
}
