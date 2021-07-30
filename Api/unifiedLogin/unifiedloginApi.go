package main

import (
	"ZeroProject/Api/unifiedLogin/internal/handler"
	"ZeroProject/Api/unifiedLogin/internal/svc"
	"ZeroProject/common/middleware"
	"flag"
	"fmt"
	"github.com/tal-tech/go-zero/rest"
)

/*
	统一登录
*/
func main() {
	flag.Parse()


	var c svc.Config
	c = svc.ClientConfig(c)
	//初始化rpc连接请求
	ctx := svc.NewServiceContext(c)


	server := rest.MustNewServer(c.RestConf,
		rest.WithNotAllowedHandler(middleware.NewCorsMiddleware().Handler()))
	defer server.Stop()
	//跨域
	server.Use(middleware.NewCorsMiddleware().Handle)

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
