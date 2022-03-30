package api

import (
	"github.com/gin-gonic/gin"
	"synod/render"
)

func (s *ObjectServer) locate(ctx *gin.Context) {
	path := ctx.Param("path")

	if path == "" {
		render.Fail().WithMessage("path is invalid").To(ctx)
		return
	}

	addr := s.subscriber.PickPeer("storage", path)

	if addr == "" {
		render.Fail().WithMessage("no storage service available").To(ctx)
		return
	}

	render.Success().With(addr).To(ctx)
}
