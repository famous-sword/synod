package discovery

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd "go.etcd.io/etcd/client/v3"
	"synod/util/logx"
)

// Subscriber to subscribe service of service dependencies
type Subscriber struct {
	cli       *etcd.Client
	names     []string
	endpoints map[string]*container
	online    int
	offline   int
}

func NewSubscriber(names ...string) *Subscriber {
	subscriber := &Subscriber{
		cli:       Registry(),
		names:     names,
		endpoints: make(map[string]*container, len(names)),
		online:    0,
		offline:   0,
	}

	for _, name := range names {
		subscriber.endpoints[name] = newContainer()
	}

	return subscriber
}

func (s *Subscriber) Subscribe() (err error) {
	for _, endpoint := range s.endpoints {
		endpoint.listenUpdates()
	}

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
			logx.Infow("loading peers", "name", key, "addr", node)
		}
	}

	return nil
}

// PickPeer select one peer from multiple peers
func (s *Subscriber) PickPeer(name, key string) string {
	if _, ok := s.endpoints[name]; !ok {
		return ""
	}

	return s.endpoints[name].next(key)
}

func (s *Subscriber) Unsubscribe() error {
	return s.cli.Close()
}

func (s *Subscriber) Health() string {
	return ""
}

// listen to service online or offline
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
					logx.Infow("peer online", "key", key, "addr", addr)
				case mvccpb.DELETE:
					s.removeNode(name, key)
					logx.Infow("peer offline", "key", key)
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
