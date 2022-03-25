package discovery

import (
	"math/rand"
	"sync"
)

type container struct {
	keys    []string
	values  []string
	indexes map[string]int
	cursor  int
	mux     sync.RWMutex
}

func newContainer() *container {
	return &container{
		keys:    make([]string, 0),
		values:  make([]string, 0),
		indexes: make(map[string]int, 0),
	}
}

func (c *container) add(key, value string) {
	c.mux.Lock()
	c.keys = append(c.keys, key)
	c.values = append(c.values, value)
	c.indexes[key] = c.cursor
	c.cursor++
	c.mux.Unlock()
}

func (c *container) get(key string) string {
	if index, ok := c.indexes[key]; ok {
		return c.values[index]
	}

	return ""
}

func (c *container) remove(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if index, ok := c.indexes[key]; ok {
		c.keys = append(c.keys[:index], c.keys[index+1:]...)
		c.values = append(c.values[:index], c.values[index+1:]...)
		delete(c.indexes, key)
		c.cursor--
	}
}

func (c *container) random() string {
	if c.cursor <= 0 {
		return ""
	}

	return c.values[rand.Intn(c.cursor)]
}
