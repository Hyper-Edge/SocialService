package main

import (
	"context"
	"log"
	"net"
    "os"
    "time"

	"github.com/al3x3n0/socialservice/apigrpc"
	"github.com/al3x3n0/socialservice/server"
	"github.com/al3x3n0/socialservice/social"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	empty "github.com/golang/protobuf/ptypes/empty"
    codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type socialServer struct {
    apigrpc.UnimplementedSocialServer
    socialClient *social.Client
}

func NewServer() *socialServer {
    tmpLogger := server.NewJSONLogger(os.Stdout, zapcore.InfoLevel, server.JSONFormat)
    config := server.ParseArgs(tmpLogger, os.Args)
	logger, _ := server.SetupLogging(tmpLogger, config)
    socialClient := social.NewClient(logger, 5*time.Second)
    return &socialServer { apigrpc.UnimplementedSocialServer{}, socialClient }
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to liesten: %v", err)
	}

	s := grpc.NewServer()
	apigrpc.RegisterSocialServer(s, NewServer())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}


func (s *socialServer) CheckGoogleToken(ctx context.Context, req *apigrpc.CheckGoogleTokenRequest) (*empty.Empty, error) {
    s.socialClient.CheckGoogleToken(ctx, req.IdToken)
    return nil, status.Errorf(codes.Unimplemented, "method CheckGoogleToken not implemented")
}

func (s *socialServer) CheckAppleToken(ctx context.Context, req *apigrpc.CheckAppleTokenRequest) (*empty.Empty, error) {
    s.socialClient.CheckAppleToken(ctx, req.BundleId, req.IdToken)
	return nil, status.Errorf(codes.Unimplemented, "method CheckAppleToken not implemented")
}
