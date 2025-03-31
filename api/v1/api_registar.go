package apiv1

import "github.com/gin-gonic/gin"

func SetupAPIRoutes(router *gin.Engine, authHandler AuthHandler, middlewares ...gin.HandlerFunc) {
	// public routes
	apiV1 := router.Group("/api/v1")
	{
		SetupAuthRoutes(apiV1, authHandler, middlewares...)
	}
}
