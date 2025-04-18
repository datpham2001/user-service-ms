package main

import (
	"fmt"
	"net/http"

	apiv1 "github.com/datpham/user-service-ms/api/v1"
	"github.com/datpham/user-service-ms/internal/delivery/http/auth"
	"github.com/datpham/user-service-ms/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (s *ServerManager) StartHttpServer(authHandler *auth.AuthHandler) {
	router := gin.New()

	// init middlewares
	loggerMiddleware := middleware.NewLoggerMiddleware(pkgLogger)
	commonMiddlewares := []middleware.CommonMiddleware{loggerMiddleware}
	middlewareManager := middleware.NewMiddlewareManager(commonMiddlewares...)

	router.Use(
		gin.Recovery(),
		middlewareManager.CommonHandle(),
	)

	setupRoutes(router, authHandler, loggerMiddleware.Handle())

	s.HTTPServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", appConfig.Server.Http.Port),
		Handler: router,
	}

	go func() {
		pkgLogger.Infof("Starting HTTP server on port %s", appConfig.Server.Http.Port)
		if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkgLogger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()
}

func setupRoutes(router *gin.Engine, authHandler *auth.AuthHandler, middlewares ...gin.HandlerFunc) {
	apiv1.SetupAPIRoutes(router, authHandler, middlewares...)
}
