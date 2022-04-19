package api

import (
	"fmt"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"synod/discovery"
	"synod/rs"
	"synod/util/logx"
	"synod/util/urlbuilder"
)

// confirm object in which storage service
func (s *Service) locate(hash string) rs.Locates {
	peers := s.subscriber.ChoosePeers("data", rs.TotalShards, discovery.EmptyExcludes())
	locates := make(map[int]string)

	for _, peer := range peers {
		id := getIdFromStorage(peer, hash)

		if id != -1 {
			locates[id] = peer
		}
	}

	return locates
}

func (s *Service) exists(hash string) (bool, error) {
	addr := s.subscriber.PickPeer("data", hash)

	if addr == "" {
		return false, ErrNoPeer
	}

	url := urlbuilder.Join(addr, "locates", hash).Build()

	response, err := http.Get(url)

	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%s responsed %d", url, response.StatusCode)
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return false, err
	}

	return fastjson.GetBool(bytes, "data", "exists"), nil
}

func getIdFromStorage(peer, hash string) int {
	b := urlbuilder.Join(peer, "locates", hash)
	response, err := http.Get(b.Build())

	if err != nil {
		logx.Errorw("get object id", "error", err)
		return -1
	}

	bytes, _ := ioutil.ReadAll(response.Body)

	return fastjson.GetInt(bytes, "data")
}
