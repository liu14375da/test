package svc

import (
	"ZeroProject/Rpc/unifiedLogin/unifiedloginclient"
	"ZeroProject/common/global"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/zrpc"
)

type ServiceContext struct {
	Config       Config
	UnifiedLogin unifiedloginclient.UnifiedLogin //调用proto request 请求
}

type Config struct {
	rest.RestConf
	UserRpc zrpc.RpcClientConf
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		UnifiedLogin: unifiedloginclient.NewUnifiedLogin(zrpc.MustNewClient(c.UserRpc)),
	}
}


func ClientConfig(c Config) Config {
	// api 服务对应的（名称,ip,端口）
	c.Name = global.UnifiedLoginApiClient.Name
	c.Host = global.UnifiedLoginApiClient.Host
	c.Port = global.UnifiedLoginApiClient.Port
	// etcd 对应的 ip 和 key
	//c.UserRpc.Etcd.Hosts = append(c.UserRpc.Etcd.Hosts, global.UnifiedLoginApiClient.Hosts)
	//c.UserRpc.Etcd.Key = global.UnifiedLoginApiClient.Key
	c.UserRpc.Endpoints = append(c.UserRpc.Endpoints, global.UnifiedLoginRpcServer.ListenOn)

	// 日志
	c.Log.Mode = global.UnifiedLoginApiClient.LogMode
	c.Log.Path = global.UnifiedLoginApiClient.Path
	c.Log.Level = global.UnifiedLoginApiClient.Level
	// 普罗米修斯
	c.Prometheus.Host = global.UnifiedLoginApiClient.PrometheusHost
	c.Prometheus.Path = global.UnifiedLoginApiClient.PrometheusPath
	c.Prometheus.Port = global.UnifiedLoginApiClient.PrometheusPort
	return c
}
