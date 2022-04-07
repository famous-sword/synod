package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"synod/metadata"
	"synod/render"
	"synod/streams"
)

var (
	ErrHashRequired = errors.New("required object hash in digest header")
	ErrInvalidName  = errors.New("invalid name")
	ErrNoPeer       = errors.New("no peer available")
)

func (s *Service) put(ctx *gin.Context) {
	hash := getHash(ctx)

	if hash == "" {
		render.OfError(ErrHashRequired).To(ctx)
		return
	}

	name := ctx.Param("name")

	if name == "" {
		render.OfError(ErrInvalidName).To(ctx)
		return
	}

	peer := s.subscriber.PickPeer("storage", hash)

	if peer == "" {
		render.OfError(ErrNoPeer).To(ctx)
		return
	}

	to := fmt.Sprintf("http://%s/objects/%s", peer, hash)

	stream := streams.NewSender(to)

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

func (s *Service) load(ctx *gin.Context) {
	name := ctx.Param("name")

	if name == "" {
		render.OfError(ErrInvalidName).To(ctx)
		return
	}

	versionId := ctx.Query("version")

	var (
		version int
		meta    metadata.Meta
		err     error
	)

	if len(versionId) != 0 {
		version, err = strconv.Atoi(versionId)

		if err != nil {
			render.OfError(err).To(ctx)
			return
		}
	}

	meta, err = s.metaManager.Get(name, version)

	if err != nil {
		render.OfError(err).To(ctx)
	}

	if meta.Hash == "" {
		render.NotFound().To(ctx)
		return
	}

	peer := s.subscriber.PickPeer("storage", meta.Hash)

	if peer == "" {
		render.OfError(ErrNoPeer).To(ctx)
		return
	}

	from := fmt.Sprintf("http://%s/objects/%s", peer, meta.Hash)

	stream, err := streams.NewPuller(from)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	io.Copy(ctx.Writer, stream)
}
