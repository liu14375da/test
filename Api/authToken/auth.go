package main

import (
	"ZeroProject/Api/authToken/internal/handler"
	"ZeroProject/Api/authToken/internal/svc"
	"ZeroProject/common/middleware"
	nacos "ZeroProject/nacos/server"
	"flag"
	"fmt"
	"github.com/tal-tech/go-zero/rest"
	"strconv"
)


func main() {
	flag.Parse()

	var c svc.Config

	//获取naCos配置中心api的信息
	c = svc.NewConfig(c)
	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf,
		rest.WithNotAllowedHandler(middleware.NewCorsMiddleware().Handler()))
	defer server.Stop()
	//跨域
	server.Use(middleware.NewCorsMiddleware().Handle)

	//路由
	handler.RegisterHandlers(server, ctx)

	metadata := RegisterMetadata(c.Port)


	//注册naCos配置中心
	nacos.RegisterService("TokenVerification", metadata)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}


func RegisterMetadata(port int) map[string]string {
	metadata := make(map[string]string)
	metadata["path"] = "/to_ken/dataByTime"
	metadata["apiType"] = "GET"
	metadata["port"] = strconv.Itoa(port)
	return metadata
}
