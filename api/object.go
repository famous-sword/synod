package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"synod/conf"
	"synod/discovery"
)

var (
	ErrInvalidAddr = errors.New("invalid addr")
)

type Object struct {
	Name       string
	Addr       string
	Server     *http.Server
	publisher  *discovery.Publisher
	subscriber *discovery.Subscriber
}

func NewObject() *Object {
	obj := &Object{}
	obj.Name = "api"
	obj.Addr = conf.String("api.address")

	handler := gin.Default()
	handler.GET("/objects/*path", obj.loadObject)
	handler.PUT("/objects/*path", obj.putObject)

	obj.Server = &http.Server{
		Handler: handler,
	}

	return obj
}

func (o *Object) Run() error {
	if o.Addr == "" {
		return ErrInvalidAddr
	}

	o.Server.Addr = o.Addr

	o.Server.RegisterOnShutdown(func() {
		var err error
		if err = o.publisher.Unpublished(); err != nil {
			log.Println(err)
		}
		if err = o.subscriber.Unsubscribe(); err != nil {
			log.Println(err)
		}
	})

	o.publisher = discovery.NewPublisher(o.Name, o.Addr)
	o.publisher.Publish()
	// o.subscriber = discovery.NewSubscriber("")

	return o.Server.ListenAndServe()
}

func (o *Object) Close() {
	o.Server.Shutdown(context.TODO())
}

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
