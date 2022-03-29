package api

import (
	"github.com/gin-gonic/gin"
	"synod/render"
)

func (s *ObjectServer) Locate(ctx *gin.Context) {
	addr := s.subscriber.Next("storage")

	if addr == "" {
		render.Fail().WithMessage("no storage service available").To(ctx)
		return
	}

	render.Success().To(ctx)
}
