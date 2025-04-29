package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/datpham/user-service-ms/config"
	"github.com/datpham/user-service-ms/internal/client/oauth"
	authHandler "github.com/datpham/user-service-ms/internal/delivery/http/auth"
	"github.com/datpham/user-service-ms/internal/infra/cache"
	"github.com/datpham/user-service-ms/internal/infra/database"
	"github.com/datpham/user-service-ms/internal/infra/rabbitmq"
	"github.com/datpham/user-service-ms/internal/pkg/httpclient"
	"github.com/datpham/user-service-ms/internal/pkg/logger"
	authRepo "github.com/datpham/user-service-ms/internal/repository/auth"
	authSvc "github.com/datpham/user-service-ms/internal/service/auth"
	tokensvc "github.com/datpham/user-service-ms/internal/service/token"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GrpcServerRegistar interface {
	RegisterGRPCHandlers(server *grpc.Server)
}

type ServerManager struct {
	GRPCServer *grpc.Server
	HTTPServer *http.Server
}

var (
	appConfig   *config.Config = &config.Config{}
	pkgLogger   *logger.Logger
	pkgCache    *cache.Cache
	pkgDatabase *database.Database
)

func init() {
	// init configuration
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	if err := config.LoadConfig(workDir+"/config", appConfig); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// init logger
	loggerConfig := logger.LoggerConfig{
		Env:         "development",
		Level:       logrus.InfoLevel,
		ServiceName: "user-ms-service",
		EnableJSON:  false,
		Fields: map[string]any{
			"version": "1.0.0",
		},
	}
	pkgLogger = logger.SetupLogger(loggerConfig)

	// init cache client
	pkgCache = cache.NewCacheClient(pkgLogger, appConfig)
	if err := pkgCache.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping cache: %v", err)
	}

	// init database
	pkgDatabase, err = database.NewDatabase(appConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func main() {
	ctx := context.Background()

	// init infra
	serverManager := &ServerManager{}
	defer func() {
		pkgDatabase.Close()
		serverManager.ShutdownGRPCServer()
		serverManager.ShutdownHTTPServer(ctx)
	}()
	dbConn := pkgDatabase.GetDB()

	rabbitMQ, err := rabbitmq.NewRabbitMQClient(pkgLogger, appConfig)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}

	if err := rabbitMQ.Setup("user-events", "user-events.*"); err != nil {
		log.Fatalf("Failed to setup RabbitMQ: %v", err)
	}

	// init client
	httpClient := httpclient.NewClient(10 * time.Second)
	oauthClient := oauth.NewOauthClient(httpClient)

	// init repositories
	authRepo := authRepo.New(dbConn)

	// init services
	tokenSvc := tokensvc.NewJwtToken(appConfig.Jwt.Secret)
	oauthSvc := tokensvc.NewOAuthService(appConfig, oauthClient)
	authSvc := authSvc.New(pkgLogger, authRepo, tokenSvc, oauthSvc, pkgCache, rabbitMQ)

	// init handlers
	authHandler := authHandler.New(authSvc)

	go func() {
		//serverManager.StartGrpcServer(grpcServerRegistry)
		serverManager.StartHttpServer(authHandler)
	}()

	<-ctx.Done()
}

func (s *ServerManager) ShutdownGRPCServer() error {
	s.GRPCServer.GracefulStop()
	return nil
}

func (s *ServerManager) ShutdownHTTPServer(ctx context.Context) error {
	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
