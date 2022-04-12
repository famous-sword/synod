package discovery

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd "go.etcd.io/etcd/client/v3"
	"math/rand"
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

// ChoosePeers select many peers by random
// name: choose from subscribed service
// n: how many you want to choose
// excludes: what can't be selected
func (s *Subscriber) ChoosePeers(name string, n int, excludes Excludes) (list []string) {
	candidates := make([]string, 0)
	peers := s.mustGetEndpoint(name).all()

	for _, peer := range peers {
		if !excludes.In(peer) {
			candidates = append(candidates, peer)
		}
	}

	length := len(candidates)

	if length < n {
		return list
	}

	p := rand.Perm(length)

	for i := 0; i < n; i++ {
		list = append(list, candidates[p[i]])
	}

	return list
}

// PickPeer select one peer by load balancer
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

func (s *Subscriber) mustGetEndpoint(name string) *container {
	if _, ok := s.endpoints[name]; !ok {
		s.endpoints[name] = newContainer()
	}

	return s.endpoints[name]
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
	s.mustGetEndpoint(name).add(key, addr)
	s.online++
}

func (s *Subscriber) removeNode(name, key string) {
	s.mustGetEndpoint(name).remove(key)
	s.offline++
}
