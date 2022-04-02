package discovery

import (
	etcd "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"synod/conf"
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
			log.Fatalln(err)
		}

		client = cli
	})

	return client
}

func Shutdown() {
	if client != nil {
		client.Close()
	}
}
