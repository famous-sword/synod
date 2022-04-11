package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"synod/stream"
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

	peer := s.subscriber.PickPeer("storage", hash)

	if peer == "" {
		return ErrNoPeer
	}

	tmp, err := stream.NewTemp(peer, hash, size)

	if err != nil {
		return err
	}

	c := util.SumHash(io.TeeReader(reader, tmp))
	logx.Debugw("sum hash in put", hash, c)

	if hash != c {
		tmp.Commit(false)
		return ErrHashCheckFail
	}

	tmp.Commit(true)

	return nil
}
