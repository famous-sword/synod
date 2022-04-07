package storage

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"synod/conf"
	"synod/discovery"
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
	handler.PUT("/objects/:name", s.put)
	handler.GET("/locates/:hash", s.exists)

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
