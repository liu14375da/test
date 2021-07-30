package main

import (
	"ZeroProject/Rpc/qyWxBinDing/internal/server"
	"ZeroProject/Rpc/qyWxBinDing/internal/svc"
	pb "ZeroProject/Rpc/qyWxBinDing/pb"
	"context"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"testing"
	"time"
)

// 统一登录 rpc 单元测试
func TestDepositServer_Deposit(t *testing.T) {
	var c zrpc.RpcServerConf
	c = svc.RpcServer(c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewQyWxBinDingServer(ctx)
	s := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		pb.RegisterQyWxBinDingServer(grpcServer, srv)
	})
	s.AddOptions(grpc.ConnectionTimeout(time.Hour))
	s.AddUnaryInterceptors(UnaryCrashInterceptor())
	s.AddStreamInterceptors(StreamCrashInterceptor)
	go s.Start()
	s.Stop()
}

func StreamCrashInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})

	return handler(srv, stream)
}

func UnaryCrashInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})

		return handler(ctx, req)
	}
}

func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(r interface{}) error {
	logx.Errorf("%+v %s", r, debug.Stack())
	return status.Errorf(codes.Internal, "panic: %v", r)
}
