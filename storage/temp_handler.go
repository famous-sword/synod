package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"strconv"
	"synod/render"
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

	_, _ = os.Create(withTemp(tmp.Uuid) + ".tmp")

	render.Success().With(u.String()).To(ctx)
}

func (s *Service) patchTemp(ctx *gin.Context) {
	u := ctx.Param("uuid")
	origin, err := ofUuid(u)

	if err != nil {
		log.Println(err)
		render.OfError(err).To(ctx)
		return
	}

	tempFileName := withTemp(u + ".tmp")

	f, err := os.OpenFile(tempFileName, os.O_WRONLY|os.O_APPEND, 0)

	if err != nil {
		log.Println(err)
		render.OfError(err).To(ctx)
		return
	}

	defer f.Close()

	_, err = io.Copy(f, ctx.Request.Body)

	if err != nil {
		log.Println(err)
		render.OfError(err).To(ctx)
		return
	}

	info, err := f.Stat()

	if err != nil {
		log.Println(err)
		render.OfError(err).To(ctx)
		return
	}

	// size not match, remove all file
	if info.Size() > origin.Size {
		_ = os.Remove(tempFileName)
		_ = os.Remove(withTemp(u + ".json"))
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

	tmp := withTemp(u + ".tmp")

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

	_ = os.Remove(withTemp(u + ".json"))
	if info.Size != stat.Size() {
		_ = os.Remove(tmp)
		render.Fail().WithMessage("size not match").To(ctx)
		return
	}

	commit(tmp, info)
}

func (s *Service) removeTemp(c *gin.Context) {
	u := c.Param("uuid")
	_ = os.Remove(withTemp(u + ".tmp"))
	_ = os.Remove(withTemp(u + ".json"))
}
