package discovery

import (
	etcd "go.etcd.io/etcd/client/v3"
	"sync"
	"synod/conf"
	"synod/util/logx"
	"time"
)

var (
	client *etcd.Client
	once   sync.Once
)

// Registry is an etcd client factory
func Registry() *etcd.Client {
	once.Do(func() {
		cli, err := etcd.New(etcd.Config{
			Endpoints:   conf.StringSlice("etcd.endpoints"),
			DialTimeout: 2 * time.Second,
			TLS:         nil,
			Username:    conf.String("etcd.username"),
			Password:    conf.String("etcd.password"),
		})

		if err != nil {
			logx.Errorw("create etcd client", "error", err)
		}

		client = cli
	})

	return client
}

func Shutdown() {
	if client != nil {
		if err := client.Close(); err != nil {
			logx.Errorw("discovery shutdown error", "error", err)
		}
	}
}
