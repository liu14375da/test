package main

import (
	nacos "ZeroProject/nacos/server"
	"flag"
	"fmt"

	"ZeroProject/Rpc/qyWxBinDing/internal/server"
	"ZeroProject/Rpc/qyWxBinDing/internal/svc"
	pb "ZeroProject/Rpc/qyWxBinDing/pb"

	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()

	var c zrpc.RpcServerConf
	c = svc.RpcServer(c)

	ctx := svc.NewServiceContext(c)
	srv := server.NewQyWxBinDingServer(ctx)

	s := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		pb.RegisterQyWxBinDingServer(grpcServer, srv)
	})
	defer s.Stop()

	nacos.RegisterService("WechatBinding",nil)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

