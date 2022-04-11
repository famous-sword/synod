package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"synod/util/render"
	"synod/util/urlbuilder"
)

// confirm object in which storage service
func (s *Service) locate(ctx *gin.Context) {
	hash := ctx.Param("hash")

	if hash == "" {
		render.Fail().WithMessage("hash is invalid").To(ctx)
		return
	}

	exists, err := s.exists(hash)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	r := gin.H{
		"exists": exists,
	}

	render.Success().With(r).To(ctx)
}

func (s *Service) exists(hash string) (bool, error) {
	addr := s.subscriber.PickPeer("storage", hash)

	if addr == "" {
		return false, ErrNoPeer
	}

	url := urlbuilder.Join(addr, "locates", hash).Build()

	response, err := http.Get(url)

	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("service: %s responsed %d", url, response.StatusCode)
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return false, err
	}

	return fastjson.GetBool(bytes, "data", "exists"), nil
}
