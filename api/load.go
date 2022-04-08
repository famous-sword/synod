package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
	"synod/metadata"
	"synod/streams"
	"synod/util/render"
)

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
		return
	}

	if meta.Hash == "" {
		// object not found
		render.NotFound().To(ctx)
		return
	}

	peer := s.subscriber.PickPeer("storage", meta.Hash)

	if peer == "" {
		render.OfError(ErrNoPeer).To(ctx)
		return
	}

	from := fmt.Sprintf("http://%s/objects/%s", peer, meta.Hash)

	stream, err := streams.NewCopyStream(from)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	_, err = io.Copy(ctx.Writer, stream)

	if err != nil {
		render.OfError(err).To(ctx)
	}
}
