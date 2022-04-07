package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"synod/render"
	"synod/streams"
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

func (s *Service) doPut(reader io.Reader, hash string, _ int64) error {
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

	to := fmt.Sprintf("http://%s/objects/%s", peer, hash)

	stream := streams.NewSender(to)

	_, err = io.Copy(stream, reader)

	if err != nil {
		return err
	}

	if err = stream.Close(); err != nil {
		return err
	}

	return nil
}
