package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"synod/render"
)

// confirm object in which storage service
func (s *Service) locate(ctx *gin.Context) {
	hash := ctx.Param("hash")

	if hash == "" {
		render.Fail().WithMessage("hash is invalid").To(ctx)
		return
	}

	addr := s.subscriber.PickPeer("storage", hash)

	if addr == "" {
		render.Fail().WithMessage("no storage service available").To(ctx)
		return
	}

	response, err := http.Get(fmt.Sprintf("http://%s/locates/%s", addr, hash))

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	if response.StatusCode != http.StatusOK {
		render.Fail().To(ctx)
		return
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	fmt.Printf("%s\n", string(bytes))

	r := gin.H{
		"peer":   addr,
		"exists": fastjson.GetBool(bytes, "data", "exists"),
	}

	render.Success().With(r).To(ctx)
}
