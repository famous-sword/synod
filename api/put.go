package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"synod/discovery"
	"synod/rs"
	"synod/util"
	"synod/util/logx"
	"synod/util/render"
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

	size := getSize(ctx)

	if err := s.doPut(ctx.Request.Body, hash, size); err != nil {
		render.OfError(err).To(ctx)
		return
	}

	if err := s.metaManager.AddVersion(name, hash, size); err != nil {
		render.OfError(err).To(ctx)
		return
	}

	render.Success().To(ctx)
}

func (s *Service) doPut(reader io.Reader, hash string, size int64) error {
	exists, err := s.exists(hash)

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	peers := s.subscriber.ChoosePeers("data", rs.TotalShards, discovery.EmptyExcludes())

	if len(peers) != rs.TotalShards {
		return ErrNotEnoughPeers
	}

	uploader, err := rs.NewUploader(peers, hash, size)

	if err != nil {
		return err
	}

	calculated := util.SumHash(io.TeeReader(reader, uploader))
	logx.Debugw("check hash", "expected", hash, "calculated", calculated)

	if hash != calculated {
		uploader.Commit(false)
		return ErrHashCheckFail
	}

	uploader.Commit(true)

	return nil
}
