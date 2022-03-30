package discovery

import (
	"context"
	etcd "go.etcd.io/etcd/client/v3"
	"log"
)

type Publisher struct {
	cli       *etcd.Client
	name      string
	addr      string
	leaseId   etcd.LeaseID
	heartbeat <-chan *etcd.LeaseKeepAliveResponse
}

func NewPublisher(name, addr string) *Publisher {
	return &Publisher{
		cli:  Registry(),
		name: name,
		addr: addr,
	}
}

func (p *Publisher) Publish() error {
	lease, err := p.cli.Grant(context.TODO(), 5)

	if err != nil {
		return err
	}

	p.leaseId = lease.ID

	_, err = p.cli.Put(context.TODO(), forPublish(p.name), p.addr, etcd.WithLease(p.leaseId))

	if err != nil {
		return err
	}

	return p.keepalive()
}

func (p *Publisher) keepalive() (err error) {
	p.heartbeat, err = p.cli.KeepAlive(context.TODO(), p.leaseId)

	if err != nil {
		return err
	}

	go func() {
		for _ = range p.heartbeat {
			// log.Printf("%s leased renew: %v\n", p.name, heartbeat.ID)
		}
	}()

	return nil
}

func (p *Publisher) Unpublished() (err error) {
	_, err = p.cli.Delete(context.TODO(), forPublish(p.name))

	if err != nil {
		return err
	}

	log.Printf("%s unpublished\n", p.name)

	_, err = p.cli.Revoke(context.TODO(), p.leaseId)

	if err != nil {
		return err
	}

	log.Printf("%s revoke lease: %v\n", p.name, p.leaseId)

	return nil
}
