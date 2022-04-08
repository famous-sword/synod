package api

import (
	"github.com/gin-gonic/gin"
	"synod/util/render"
)

func (s *Service) delete(c *gin.Context) {
	name := c.Param("name")
	meta, err := s.metaManager.LatestVersion(name)

	if err != nil {
		render.OfError(err).To(c)
		return
	}

	err = s.metaManager.Put(name, meta.Version + 1, 0, "")

	if err != nil {
		render.OfError(err).To(c)
		return
	}

	render.Success().To(c)
}
