package lb

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type LoadBalancer struct {
	hash     Hash
	replicas int
	keys     []int // sorted
	hashes   map[int]string
}

func New(replicas int, fn Hash) *LoadBalancer {
	ch := &LoadBalancer{
		hash:     fn,
		replicas: replicas,
		hashes:   make(map[int]string),
	}

	if ch.hash == nil {
		ch.hash = crc32.ChecksumIEEE
	}

	return ch
}

func (b *LoadBalancer) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < b.replicas; i++ {
			hash := int(b.hash([]byte(strconv.Itoa(i) + key)))
			b.keys = append(b.keys, hash)
			b.hashes[hash] = key
		}
	}
	sort.Ints(b.keys)
}

func (b *LoadBalancer) Get(key string) string {
	if len(b.keys) == 0 {
		return ""
	}

	hash := int(b.hash([]byte(key)))

	// Binary search for appropriate replica.
	index := sort.Search(len(b.keys), func(i int) bool {
		return b.keys[i] >= hash
	})

	return b.hashes[b.keys[index%len(b.keys)]]
}
