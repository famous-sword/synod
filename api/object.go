package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"synod/conf"
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

	binary, err := os.Open(diskPath(path))

	if err != nil {
		render.Fail().WithError(err).To(ctx)
		return
	}

	defer binary.Close()

	io.Copy(ctx.Writer, binary)
}

func diskPath(path string) string {
	disk := conf.String("storage.local")
	return filepath.Join(disk, path)
}
