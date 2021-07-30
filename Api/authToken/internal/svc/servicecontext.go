package svc

import (
	"ZeroProject/common/global"
	"github.com/tal-tech/go-zero/rest"
)

type ServiceContext struct {
	Config Config
}

type Config struct {
	rest.RestConf
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

func NewConfig(c Config) Config {
	// api 服务对应的（名称,ip,端口）
	c.Name = global.AuthToken.Name
	c.Host = global.AuthToken.Host
	c.Port = global.AuthToken.Port
	//开启日志路径
	c.Log.Mode = global.AuthToken.LogMode
	c.Log.Path = global.AuthToken.Path
	c.Log.Level = global.AuthToken.Level
	// 普罗米修斯
	c.Prometheus.Host = global.AuthToken.PrometheusHost
	c.Prometheus.Path = global.AuthToken.PrometheusPath
	c.Prometheus.Port = global.AuthToken.PrometheusPort
	return c
}
