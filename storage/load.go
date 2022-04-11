package storage

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"synod/util/render"
)

func (s *Service) load(ctx *gin.Context) {
	name := ctx.Param("name")

	if name == "" {
		render.Fail().WithMessage("invalid name").To(ctx)
		return
	}

	binary, err := os.Open(Workdir(name))

	if err != nil {
		render.Fail().WithError(err).To(ctx)
		return
	}

	defer binary.Close()

	io.Copy(ctx.Writer, binary)
}

func (s *Service) exists(ctx *gin.Context) {
	hash := ctx.Param("hash")

	r := gin.H{
		"exists": exists(hash),
	}

	render.Success().With(r).To(ctx)
}

func exists(hash string) bool {
	_, err := os.Stat(Workdir(hash))

	if err == nil {
		return true
	}

	return os.IsExist(err)
}
