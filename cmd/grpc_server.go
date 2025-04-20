package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MAX_RETRIES = 3
)

func (s *ServerManager) StartGrpcServer(sr GrpcServerRegistar) {
	if !checkGrpcConfig() {
		pkgLogger.Fatal("Grpc server is not configured")
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		LoggingUnaryInterceptor,
	}

	s.GRPCServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(16*1024*1024), // 16MB
		grpc.MaxSendMsgSize(16*1024*1024), // 16MB
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)
	sr.RegisterGRPCHandlers(s.GRPCServer)

	grpcAddress := net.JoinHostPort(appConfig.Server.Grpc.Host, appConfig.Server.Grpc.Port)
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen to tcp address: %w", err)
	}

	go func() {
		pkgLogger.Infof("Starting GRPC server on port %s", appConfig.Server.Grpc.Port)
		if err := s.GRPCServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %w", err)
		}
	}()
}

func checkGrpcConfig() bool {
	return appConfig.Server.Grpc.Port != "" && appConfig.Server.Grpc.Host != ""
}

func LoggingUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	pkgLogger.Printf("GRPC request: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		pkgLogger.Printf("GRPC response: %v, error: %w, duration: %s", resp, err, time.Since(start))
	} else {
		pkgLogger.Printf("GRPC response: %v, duration: %s", resp, time.Since(start))
	}

	return resp, err
}

func RetryUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	for attempt := 0; attempt < MAX_RETRIES; attempt++ {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil || !isRetryable(err) {
			return err
		}

		pkgLogger.Infof("Retrying %s, attempt %d after error: %v", method, attempt+1, err)
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

func isRetryable(err error) bool {
	code := status.Code(err)
	return code == codes.Unavailable || code == codes.DeadlineExceeded
}
