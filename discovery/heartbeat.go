package discovery

import (
	"context"
	"fmt"
	cli "go.etcd.io/etcd/client/v3"
	"synod/conf"
)

// StartHeartbeat start heartbeat to keep service alive
func StartHeartbeat() error {
	lease, err := etcd.Grant(context.TODO(), 5)

	if err != nil {
		return err
	}

	// register self to etcd
	_, err = etcd.Put(context.TODO(), conf.String("app.id"), conf.String("api.address"), cli.WithLease(lease.ID))

	if err != nil {
		return err
	}

	aliveChan, err := etcd.KeepAlive(context.TODO(), lease.ID)

	if err != nil {
		return err
	}

	for alive := range aliveChan {
		fmt.Printf("leased renewed: %v\n", alive.ID)
	}

	return err
}
