package grpc

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/betchi/tracer"
	"github.com/swagchat/chat-api/datastore"

	logger "github.com/betchi/zapper"
	"github.com/swagchat/chat-api/config"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		workspace := ""

		headers, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if v, ok := headers[strings.ToLower(config.HeaderWorkspace)]; ok {
				if len(v) > 0 {
					workspace = v[0]
				}
			}
		}

		if workspace == "" {
			workspace = config.Config().Datastore.Database
		}

		ctx = context.WithValue(ctx, config.CtxWorkspace, workspace)

		ctx, _ = tracer.StartTransaction(ctx, fmt.Sprintf("%s:%v", info.FullMethod, info.Server), "GRPC")
		defer tracer.CloseTransaction(ctx)

		reply, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return reply, nil
	}
}

// Run runs GRPC API server
func Run(ctx context.Context) {
	cfg := config.Config()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to serve %s server[GRPC]. %v", config.AppName, err))
	}

	ops := []grpc.ServerOption{grpc.UnaryInterceptor(unaryServerInterceptor())}
	s := grpc.NewServer(ops...)
	logger.Info(fmt.Sprintf("Starting %s server[GRPC] on listen tcp :%s", config.AppName, cfg.GRPCPort))

	scpb.RegisterBlockUserServiceServer(s, &blockUserServiceServer{})
	scpb.RegisterDeviceServiceServer(s, &deviceServiceServer{})
	scpb.RegisterMessageServiceServer(s, &messageServer{})
	scpb.RegisterRoomUserServiceServer(s, &roomUserServiceServer{})
	scpb.RegisterUserServiceServer(s, &userServiceServer{})
	scpb.RegisterUserRoleServiceServer(s, &userRoleServiceServer{})

	reflection.Register(s)

	errCh := make(chan error)
	go func() {
		errCh <- s.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		logger.Info(fmt.Sprintf("Stopping %s server[GRPC]", config.AppName))
		datastore.Provider(ctx).Close()
		s.GracefulStop()
	case err = <-errCh:
		logger.Error(fmt.Sprintf("Failed to serve %s server[GRPC]. %v", config.AppName, err))
		datastore.Provider(ctx).Close()
		s.GracefulStop()
	}
}
