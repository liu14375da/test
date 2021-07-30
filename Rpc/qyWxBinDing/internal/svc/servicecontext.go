package svc

import (
	"ZeroProject/Rpc/qyWxBinDing/internal/model/bandingWx"
	"ZeroProject/common/global"
	_ "database/sql"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/zrpc"
)

type ServiceContext struct {
	zr        zrpc.RpcServerConf
	BandingWx bandingWx.BinDing
}

func NewServiceContext(c zrpc.RpcServerConf) *ServiceContext {
	// redis 参数连接
	cluster := cache.CacheConf{
		cache.NodeConf{
			RedisConf: redis.RedisConf{
				Host: global.Redis.Host,
				Pass: global.Redis.Pass,
			},
			Weight: 100,
		},
	}
	//获取数据库连接
	conn := sqlx.NewSqlConn("mssql", global.UnifiedLogin.SqlConn)
	return &ServiceContext{
		zr:        c,
		BandingWx: bandingWx.NewBandingWx(conn, cluster),
	}
}

// 添加统一登录对应的配置文件（etcd的ip地址,key值和rpc服务名称,监听地址,日志）
func RpcServer(zr zrpc.RpcServerConf) zrpc.RpcServerConf {
	//zr.Etcd.Hosts = append(zr.Etcd.Hosts, global.QyWxRpc.BingDingEtcdHost)
	//zr.Etcd.Key = global.QyWxRpc.BingDingEtcdKey

	zr.Name = global.QyWxRpc.BingDingName
	zr.ListenOn = global.QyWxRpc.BingDingListenOn
	// 日志相关
	zr.ServiceConf.Log.Mode = global.QyWxRpc.LogMode
	zr.ServiceConf.Log.Path = global.QyWxRpc.Path
	zr.ServiceConf.Log.Level = global.QyWxRpc.Level
	// Prometheus（普罗米修斯报警配置信息）
	zr.ServiceConf.Prometheus.Host = global.QyWxRpc.PrometheusHost
	zr.ServiceConf.Prometheus.Port = global.QyWxRpc.PrometheusPort
	zr.ServiceConf.Prometheus.Path = global.QyWxRpc.PrometheusPath
	return zr
}
