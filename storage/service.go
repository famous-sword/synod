package storage

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"synod/conf"
	"synod/discovery"
	"time"
)

// Service is a local disk storage service
type Service struct {
	Name   string
	Addr   string
	Server *http.Server
	// implement features to handler storage
	Handler    http.Handler
	publisher  *discovery.Publisher
	subscriber *discovery.Subscriber
}

func New() *Service {
	return &Service{
		Name: "storage",
		Addr: conf.String("storage.addr"),
	}
}

func (s *Service) Run() error {
	handler := gin.Default()
	handler.GET("/objects/:name", s.load)
	handler.GET("/locates/:hash", s.exists)

	handler.POST("/tmp/:name", s.createTemp)
	handler.PATCH("/tmp/:uuid", s.patchTemp)
	handler.PUT("/tmp/:uuid", s.putTemp)
	handler.DELETE("/tmp/:uuid", s.removeTemp)

	server := &http.Server{
		Addr:         s.Addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.SetKeepAlivesEnabled(false)

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

// withWorkdir generate full path in work dir
func withWorkdir(name string) string {
	return filepath.Join(conf.String("storage.workdir"), name)
}

// withTemp generate full path in temp dir
func withTemp(name string) string {
	return filepath.Join(conf.String("storage.temp"), name)
}
