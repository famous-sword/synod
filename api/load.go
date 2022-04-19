package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"synod/discovery"
	"synod/metadata"
	"synod/rs"
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

	locates := s.locate(meta.Hash)

	if len(locates) < rs.NumDataShard {
		render.Fail().WithMessage("locate fail").To(ctx)
		return
	}

	var servers []string

	if len(servers) != rs.TotalShards {
		servers = s.subscriber.ChoosePeers(
			"data",
			rs.TotalShards-len(locates),
			discovery.LoadExcludes(locates),
		)
	}

	stream, err := rs.NewDownloader(locates, servers, meta.Hash, meta.Size)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	_, err = io.Copy(ctx.Writer, stream)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	stream.Close()
}
