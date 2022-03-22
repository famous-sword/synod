package discovery

import (
	"context"
	"fmt"
	cli "go.etcd.io/etcd/client/v3"
	"log"
	"synod/conf"
	"time"
)

var etcd *cli.Client

func Startup() error {
	client, err := cli.New(cli.Config{
		Endpoints:   conf.StringSlice("etcd.endpoints"),
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		return err
	}

	etcd = client

	return nil
}

func Close() {
	fmt.Println()
	log.Println("unregister service...")
	_, err := etcd.Delete(context.TODO(), conf.String("app.id"))

	if err != nil {
		log.Println(err)
	}

	fmt.Println("synod will to close etcd")

	err = etcd.Close()

	if err != nil {
		log.Println(err)
	}
}
