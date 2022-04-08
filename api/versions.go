package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"synod/util/render"
)

func (s *Service) versions(ctx *gin.Context) {
	name := ctx.Param("name")

	if name == "" {
		render.OfError(ErrInvalidName).To(ctx)
		return
	}

	from := 0
	size := 1000

	for {
		metas, err := s.metaManager.Versions(name, from, size)

		if err != nil {
			render.OfError(err).To(ctx)
			return
		}

		for i := range metas {
			bytes, _ := json.Marshal(metas[i])
			ctx.Writer.Write(bytes)
			ctx.Writer.Write([]byte("\n"))
		}

		if len(metas) != size {
			return
		}

		from += size
	}
}
