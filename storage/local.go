package storage

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"synod/conf"
	"synod/discovery"
	"synod/render"
)

type LocalStorage struct {
	Name       string
	Addr       string
	Server     *http.Server
	Handler    http.Handler
	publisher  *discovery.Publisher
	subscriber *discovery.Subscriber
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (s *LocalStorage) Run() error {
	s.Name = "storage"
	s.Addr = conf.String("storage.addr")

	handler := gin.Default()
	handler.GET("/objects/*path", s.loadObject)
	handler.PUT("/objects/*path", s.putObject)

	server := &http.Server{
		Addr:    s.Addr,
		Handler: handler,
	}

	s.Server = server
	s.publisher = discovery.NewPublisher(s.Name, s.Addr)
	s.subscriber = discovery.NewSubscriber("api")

	s.publisher.Publish()
	s.subscriber.Subscribe()

	s.Server.RegisterOnShutdown(func() {
		s.publisher.Unpublished()
	})

	return s.Server.ListenAndServe()
}

func (s *LocalStorage) putObject(ctx *gin.Context) {
	path := ctx.Param("path")

	if path == "" {
		render.Fail().WithMessage("invalid path").To(ctx)
		return
	}

	file, err := os.Create(diskPath(path))

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

	binary, err := os.Open(diskPath(path))

	if err != nil {
		render.Fail().WithError(err).To(ctx)
		return
	}

	defer binary.Close()

	io.Copy(ctx.Writer, binary)
}

func diskPath(path string) string {
	disk := conf.String("storage.local")
	return filepath.Join(disk, path)
}
