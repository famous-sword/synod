package discovery

type Excludes map[int]string

func EmptyExcludes() Excludes {
	return Excludes{}
}

func LoadExcludes(peers map[int]string) Excludes {
	e := EmptyExcludes()

	for i, peer := range peers {
		e[i] = peer
	}

	return e
}

func (ex Excludes) In(peer string) bool {
	for _, e := range ex {
		if e == peer {
			return true
		}
	}

	return false
}

func (ex Excludes) Add(idx int, peer string) {
	ex[idx] = peer
}
