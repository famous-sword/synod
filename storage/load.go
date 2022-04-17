package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strings"
	"synod/util/logx"
	"synod/util/render"
)

func (s *Service) load(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		render.Fail().WithMessage("invalid name").To(c)
		return
	}

	path := getFilePath(name)

	if path == "" {
		render.Fail().WithMessage("file not found").To(c)
		return
	}

	smartCopy(c.Writer, path)
}

func (s *Service) locate(c *gin.Context) {
	hash := c.Param("hash")
	id := locator.TempId(hash)

	render.Success().With(id).To(c)
}

func (s *Service) exists(c *gin.Context) {
	hash := c.Param("hash")

	r := gin.H{
		"exists": locator.Has(hash),
	}

	render.Success().With(r).To(c)
}

func getFilePath(name string) string {
	files, _ := filepath.Glob(DataPath(name + ".*"))

	if len(files) != 1 {
		return ""
	}

	file := files[0]
	device := sha256.New()
	smartCopy(device, file)
	expected := hex.EncodeToString(device.Sum(nil))
	hash := strings.Split(file, ".")[2]

	if hash != expected {
		logx.Errorw("hash check error", "file", file)
		_ = os.Remove(file)
		locator.Forget(hash)
		return ""
	}

	return file
}

func smartCopy(w io.Writer, name string) {
	file, _ := os.Open(name)
	defer file.Close()
	_, _ = io.Copy(w, file)
}
