package storage

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"synod/conf"
	"synod/discovery"
)

// LocalStorage is a local disk storage service
type LocalStorage struct {
	Name   string
	Addr   string
	Server *http.Server
	// implement features to handler storage
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
