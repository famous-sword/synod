package discovery

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
)

type Subscriber struct {
	cli     *etcd.Client
	nodes   sync.Map
	name    string
	online  int
	offline int
}

func NewSubscriber(name string) *Subscriber {
	return &Subscriber{
		cli:     Registry(),
		name:    name,
		online:  0,
		offline: 0,
	}
}

func (s *Subscriber) Subscribe() error {
	loaded, err := s.cli.Get(context.TODO(), forSubscribe(s.name), etcd.WithPrefix())

	if err != nil {
		return err
	}

	for _, kv := range loaded.Kvs {
		key := string(kv.Key)
		node := string(kv.Value)
		if node != "" {
			s.addNode(key, node)
			log.Printf("add node from loading: %s => %s\n", key, node)
		}
	}

	go func() {
		s.listen()
	}()

	return nil
}

func (s *Subscriber) Next() string {
	// todo
	return ""
}

func (s *Subscriber) Unsubscribe() error {
	return s.cli.Close()
}

func (s *Subscriber) listen() {
	watcher := s.cli.Watch(context.TODO(), forSubscribe(s.name), etcd.WithPrefix())

	for action := range watcher {
		for _, event := range action.Events {
			key := string(event.Kv.Key)

			if key != "" {
				switch event.Type {
				case mvccpb.PUT:
					addr := string(event.Kv.Value)
					s.addNode(key, addr)
					log.Printf("add node from listening: %s => %s\n", key, addr)
				case mvccpb.DELETE:
					s.removeNode(key)
					log.Printf("remove node from listening: %s\n", key)
				}
			}
		}
	}
}

func (s *Subscriber) addNode(key, addr string) {
	s.nodes.Store(key, addr)
	s.online++
}

func (s *Subscriber) removeNode(key string) {
	s.nodes.Delete(key)
	s.offline++
}
