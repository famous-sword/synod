package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"synod/conf"
)

type Object struct{}

func (o *Object) putObject(ctx *gin.Context) {
	path := ctx.Param("path")

	if path == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid file path",
		})
		return
	}

	file, err := os.Create(diskPath(path))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer file.Close()

	written, err := io.Copy(file, ctx.Request.Body)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"written": written})
}

func (o *Object) loadObject(ctx *gin.Context) {
	path := ctx.Param("path")

	if path == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid file path",
		})
		return
	}

	binary, err := os.Open(diskPath(path))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer binary.Close()

	io.Copy(ctx.Writer, binary)
}

func diskPath(path string) string {
	disk := conf.String("storage.local")
	return filepath.Join(disk, path)
}
