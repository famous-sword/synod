package storage

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"synod/conf"
	"synod/render"
)

func (s *Service) put(ctx *gin.Context) {
	name := ctx.Param("name")

	if name == "" {
		render.Fail().WithMessage("invalid name").To(ctx)
		return
	}

	file, err := os.Create(withWorkdir(name))

	if err != nil {
		render.Fail().WithError(err).To(ctx)
		return
	}

	defer file.Close()

	written, err := io.Copy(file, ctx.Request.Body)

	if err != nil {
		render.Fail().WithError(err).To(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"written": written})
}

func (s *Service) load(ctx *gin.Context) {
	name := ctx.Param("name")

	if name == "" {
		render.Fail().WithMessage("invalid name").To(ctx)
		return
	}

	binary, err := os.Open(withWorkdir(name))

	if err != nil {
		render.Fail().WithError(err).To(ctx)
		return
	}

	defer binary.Close()

	io.Copy(ctx.Writer, binary)
}

// withWorkdir join full path with workdir
func withWorkdir(name string) string {
	return filepath.Join(conf.String("storage.workdir"), name)
}
