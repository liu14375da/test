package svc

import (
	"ZeroProject/Rpc/qyWxBinDing/qywxbindingclient"
	"ZeroProject/common/global"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/zrpc"
)

type ServiceContext struct {
	Config           Config
	QyWxBinDingLogin qywxbindingclient.QyWxBinDing //调用proto request 请求
}

type Config struct {
	rest.RestConf
	QyWxRpc zrpc.RpcClientConf
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		QyWxBinDingLogin: qywxbindingclient.NewQyWxBinDing(zrpc.MustNewClient(c.QyWxRpc)),
	}
}

func ClientConfig(c Config) Config {
	// api 服务对应的（名称,ip,端口）
	c.Name = global.QyWxApi.BingDingName
	c.Host = global.QyWxApi.BingDingHost
	c.Port = global.QyWxApi.BingDingPort

	// etcd 对应的 ip 和 key
	//c.QyWxRpc.Etcd.Hosts = append(c.QyWxRpc.Etcd.Hosts, global.QyWxApi.BingDingEtcdHost)
	//c.QyWxRpc.Etcd.Key = global.QyWxApi.BingDingEtcdKey

	c.QyWxRpc.Endpoints = append(c.QyWxRpc.Endpoints, global.QyWxRpc.BingDingListenOn)

	// 日志
	c.Log.Mode = global.QyWxApi.LogMode
	c.Log.Path = global.QyWxApi.Path
	c.Log.Level = global.QyWxApi.Level
	// 普罗米修斯
	c.Prometheus.Host = global.QyWxApi.PrometheusHost
	c.Prometheus.Path = global.QyWxApi.PrometheusPath
	c.Prometheus.Port = global.QyWxApi.PrometheusPort
	return c
}
