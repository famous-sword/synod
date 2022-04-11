package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"os"
	"strconv"
	"synod/util/logx"
	"synod/util/render"
)

func (s *Service) createTemp(ctx *gin.Context) {
	u, _ := uuid.NewUUID()
	name := ctx.Param("name")
	size, err := strconv.ParseInt(ctx.GetHeader("size"), 0, 64)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	tmp := &Temp{Uuid: u.String(), Name: name, Size: size}

	if err = tmp.saveInfo(); err != nil {
		render.OfError(err).To(ctx)
		return
	}

	_, _ = os.Create(TempDir(tmp.Uuid + extTemp))

	render.Success().With(u.String()).To(ctx)
}

func (s *Service) patchTemp(ctx *gin.Context) {
	u := ctx.Param("uuid")
	origin, err := ofUuid(u)

	if err != nil {
		logx.Errorw("patch temp error", "error", err)
		render.OfError(err).To(ctx)
		return
	}

	tempFileName := TempDir(u + extTemp)

	f, err := os.OpenFile(tempFileName, os.O_WRONLY|os.O_APPEND, 0)

	if err != nil {
		logx.Errorw("patch temp error on open file", "error", err)
		render.OfError(err).To(ctx)
		return
	}

	defer f.Close()

	_, err = io.Copy(f, ctx.Request.Body)

	if err != nil {
		logx.Errorw("patch temp error on copy", "error", err)
		render.OfError(err).To(ctx)
		return
	}

	info, err := f.Stat()

	if err != nil {
		logx.Errorw("patch temp error on stat", "error", err)
		render.OfError(err).To(ctx)
		return
	}

	// size not match, remove all file
	if info.Size() > origin.Size {
		_ = os.Remove(tempFileName)
		_ = os.Remove(TempDir(u + extInfo))
	}

	render.Success().To(ctx)
}

func (s *Service) putTemp(ctx *gin.Context) {
	u := ctx.Param("uuid")

	info, err := ofUuid(u)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	tmp := TempDir(u + extTemp)

	f, err := os.Open(tmp)

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	defer f.Close()
	stat, err := f.Stat()

	if err != nil {
		render.OfError(err).To(ctx)
		return
	}

	_ = os.Remove(TempDir(u + extInfo))
	if info.Size != stat.Size() {
		_ = os.Remove(tmp)
		render.Fail().WithMessage("size not match").To(ctx)
		return
	}

	commit(tmp, info)
}

func (s *Service) removeTemp(c *gin.Context) {
	u := c.Param("uuid")
	_ = os.Remove(TempDir(u + extTemp))
	_ = os.Remove(TempDir(u + extInfo))
}
