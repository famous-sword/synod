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

type ObjectServer struct {
	Name       string
	Addr       string
	Server     *http.Server
	publisher  *discovery.Publisher
	subscriber *discovery.Subscriber
}

func NewObjectServer() *ObjectServer {
	obj := &ObjectServer{}
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

func (s *ObjectServer) Run() error {
	if s.Addr == "" {
		return ErrInvalidAddr
	}

	s.Server.Addr = s.Addr

	s.Server.RegisterOnShutdown(func() {
		log.Println()
		
		var err error
		if err = s.publisher.Unpublished(); err != nil {
			log.Println(err)
		}
		if err = s.subscriber.Unsubscribe(); err != nil {
			log.Println(err)
		}
	})

	s.publisher = discovery.NewPublisher(s.Name, s.Addr)
	s.publisher.Publish()
	// s.subscriber = discovery.NewSubscriber("")

	return s.Server.ListenAndServe()
}

func (s *ObjectServer) Close() {
	s.Server.Shutdown(context.TODO())
}

func (s *ObjectServer) putObject(ctx *gin.Context) {
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

func (s *ObjectServer) loadObject(ctx *gin.Context) {
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
