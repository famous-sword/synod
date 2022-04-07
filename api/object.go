package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"synod/render"
	"synod/streams"
)

var (
	ErrHashRequired = errors.New("required object hash in digest header")
)

func (s *RESTServer) putObject(ctx *gin.Context) {
	hash := getHash(ctx)

	if hash == "" {
		render.OfError(ErrHashRequired).To(ctx)
		return
	}

	name := ctx.Param("name")

	if name == "" {
		render.Fail().WithMessage("invalid name").To(ctx)
		return
	}

	peer := s.subscriber.PickPeer("storage", hash)

	if peer == "" {
		render.Fail().WithMessage("no peer available").To(ctx)
		return
	}

	to := fmt.Sprintf("http://%s/objects/%s", peer, hash)

	stream := streams.NewPutStream(to)

	io.Copy(stream, ctx.Request.Body)

	if err := stream.Close(); err != nil {
		render.OfError(err).To(ctx)
		return
	}

	size := getSize(ctx)

	err := s.metaManager.AddVersion(name, hash, size)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	render.Success().To(ctx)
}

func (s *RESTServer) loadObject(ctx *gin.Context) {
	path := ctx.Param("path")

	if path == "" {
		render.Fail().WithMessage("invalid path").To(ctx)
		return
	}

	peer := s.subscriber.PickPeer("storage", path)

	if peer == "" {
		render.Fail().WithMessage("no peer available").To(ctx)
		return
	}

	from := fmt.Sprintf("http://%s/objects%s", peer, path)

	stream, err := streams.NewFetchStream(from)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	io.Copy(ctx.Writer, stream)
}
