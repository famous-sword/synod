package lb

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type (
	// Hash a hash function
	Hash func(data []byte) uint32

	// Map is a loading balancer
	Map struct {
		hash     Hash
		replicas int

		// hashes stores the hashed value of the peer
		hashes []int
		// peers `hash` => `peer`
		peers map[int]string
	}
)

// New create a new consistent hash load balancer
func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		peers:    make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add many peers to replicas
func (m *Map) Add(peers ...string) {
	for _, peer := range peers {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + peer)))
			m.hashes = append(m.hashes, hash)
			m.peers[hash] = peer
		}
	}

	sort.Ints(m.hashes)
}

// Get a peers by a key hash
func (m *Map) Get(key string) string {
	if len(m.hashes) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	// Binary search for appropriate replica.
	index := sort.Search(len(m.hashes), func(i int) bool {
		return m.hashes[i] >= hash
	})

	return m.peers[m.hashes[index%len(m.hashes)]]
}
