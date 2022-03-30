package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"synod/render"
	"synod/streams"
)

func (s *RESTServer) putObject(ctx *gin.Context) {
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

	to := fmt.Sprintf("http://%s/objects%s", peer, path)

	stream := streams.NewPutStream(to)

	io.Copy(stream, ctx.Request.Body)

	if err := stream.Close(); err != nil {
		render.OfError(err).To(ctx)
		return
	}

	r := gin.H{
		"put": to,
	}

	render.Success().With(r).To(ctx)
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
