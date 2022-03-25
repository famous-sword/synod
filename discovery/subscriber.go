package discovery

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd "go.etcd.io/etcd/client/v3"
	"log"
)

type Subscriber struct {
	cli       *etcd.Client
	names     []string
	endpoints map[string]*container
	online    int
	offline   int
}

func NewSubscriber(names ...string) *Subscriber {
	return &Subscriber{
		cli:       Registry(),
		names:     names,
		endpoints: make(map[string]*container, len(names)),
		online:    0,
		offline:   0,
	}
}

func (s *Subscriber) Subscribe() (err error) {
	for _, name := range s.names {
		if err = s.subscribe(name); err != nil {
			return err
		}

		go func(n string) {
			s.listen(n)
		}(name)
	}

	return err
}

func (s *Subscriber) subscribe(name string) error {
	loaded, err := s.cli.Get(context.TODO(), forSubscribe(name), etcd.WithPrefix())

	if err != nil {
		return err
	}

	for _, kv := range loaded.Kvs {
		key := string(kv.Key)
		node := string(kv.Value)
		if node != "" {
			s.addNode(name, key, node)
			log.Printf("add node from loading: %s => %s\n", key, node)
		}
	}

	return nil
}

func (s *Subscriber) Next(name string) string {
	if _, ok := s.endpoints[name]; !ok {
		return ""
	}

	return s.endpoints[name].random()
}

func (s *Subscriber) Unsubscribe() error {
	return s.cli.Close()
}

func (s *Subscriber) listen(name string) {
	watcher := s.cli.Watch(context.TODO(), forSubscribe(name), etcd.WithPrefix())

	for action := range watcher {
		for _, event := range action.Events {
			key := string(event.Kv.Key)

			if key != "" {
				switch event.Type {
				case mvccpb.PUT:
					addr := string(event.Kv.Value)
					s.addNode(name, key, addr)
					log.Printf("add node from listening: %s => %s\n", key, addr)
				case mvccpb.DELETE:
					s.removeNode(name, key)
					log.Printf("remove node from listening: %s\n", key)
				}
			}
		}
	}
}

func (s *Subscriber) addNode(name, key, addr string) {
	if _, ok := s.endpoints[name]; !ok {
		s.endpoints[name] = newContainer()
	}

	s.endpoints[name].add(key, addr)
	s.online++
}

func (s *Subscriber) removeNode(name, key string) {
	if _, ok := s.endpoints[name]; !ok {
		s.endpoints[name] = newContainer()
	}

	s.endpoints[name].remove(key)
	s.offline++
}
