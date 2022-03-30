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

func (s *LocalStorage) putObject(ctx *gin.Context) {
	path := ctx.Param("path")

	if path == "" {
		render.Fail().WithMessage("invalid path").To(ctx)
		return
	}

	file, err := os.Create(StoreDisk(path))

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

func (s *LocalStorage) loadObject(ctx *gin.Context) {
	path := ctx.Param("path")

	if path == "" {
		render.Fail().WithMessage("invalid path").To(ctx)
		return
	}

	binary, err := os.Open(StoreDisk(path))

	if err != nil {
		render.Fail().WithError(err).To(ctx)
		return
	}

	defer binary.Close()

	io.Copy(ctx.Writer, binary)
}

func StoreDisk(path string) string {
	disk := conf.String("storage.local")
	return filepath.Join(disk, path)
}
