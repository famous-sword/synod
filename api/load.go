package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"synod/metadata"
	"synod/stream"
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
		render.Fail().WithStatus(http.StatusNotFound).To(ctx)
		return
	}

	peer := s.subscriber.PickPeer("storage", meta.Hash)

	if peer == "" {
		render.OfError(ErrNoPeer).To(ctx)
		return
	}

	copier, err := stream.NewCopier(peer, meta.Hash)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	_, err = io.Copy(ctx.Writer, copier)

	if err != nil {
		render.OfError(err).To(ctx)
	}
}
