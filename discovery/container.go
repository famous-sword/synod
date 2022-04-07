package discovery

import (
	"log"
	"sync"
	"synod/discovery/lb"
)

var DefaultReplicas = 11

// container store service information of a service
// every server has each container
type container struct {
	// service names
	keys []string
	// service addr
	values []string
	// index of service name
	indexes map[string]int
	// current index
	cursor int
	mux    sync.RWMutex

	// event of service online or offline
	changed chan int
	// loading balancer for pick peer of service
	balancer *lb.Map
}

func newContainer() *container {
	return &container{
		keys:     make([]string, 0),
		values:   make([]string, 0),
		indexes:  make(map[string]int, 0),
		balancer: lb.New(DefaultReplicas, nil),

		changed: make(chan int),
	}
}

func (c *container) add(name, value string) {
	c.mux.Lock()
	c.keys = append(c.keys, name)
	c.values = append(c.values, value)
	c.indexes[name] = c.cursor
	c.cursor++
	c.mux.Unlock()

	c.changed <- c.cursor
}

func (c *container) get(name string) string {
	if index, ok := c.indexes[name]; ok {
		return c.values[index]
	}

	return ""
}

func (c *container) remove(name string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if index, ok := c.indexes[name]; ok {
		c.keys = append(c.keys[:index], c.keys[index+1:]...)
		c.values = append(c.values[:index], c.values[index+1:]...)
		delete(c.indexes, name)
		c.cursor--
	}

	c.changed <- c.cursor
}

func (c *container) next(key string) string {
	return c.balancer.Get(key)
}

// listenUpdates when service changed, load balancer reload
func (c *container) listenUpdates() {
	go func() {
		for {
			select {
			case <-c.changed:
				c.balancer.Add(c.values...)
				log.Println("refresh load balancer peers...")
			}
		}
	}()
}
