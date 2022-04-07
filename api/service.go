package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"synod/conf"
	"synod/discovery"
	"synod/metadata"
)

var (
	ErrInvalidAddr = errors.New("invalid addr")
)

// Service is an api server for front user
// it provides rest api
type Service struct {
	// service name for register to discover
	// and other services subscribes
	Name string
	Addr string
	// a running http server
	Server *http.Server
	// used to publish self to discovery
	publisher *discovery.Publisher

	// used to subscribe other services
	subscriber *discovery.Subscriber

	metaManager metadata.Manager
}

func New() *Service {
	obj := &Service{}
	obj.Name = "api"
	obj.Addr = conf.String("api.addr")

	handler := gin.Default()
	handler.GET("/objects/:name", obj.load)
	handler.PUT("/objects/:name", obj.put)
	handler.GET("/versions/:name", obj.versions)
	handler.GET("/locates/:hash", obj.locate)

	obj.Server = &http.Server{
		Handler: handler,
	}

	obj.metaManager = metadata.New()

	return obj
}

func (s *Service) Run() error {
	if s.Addr == "" {
		return ErrInvalidAddr
	}

	s.Server.Addr = s.Addr
	s.publisher = discovery.NewPublisher(s.Name, s.Addr)
	s.publisher.Publish()
	s.subscriber = discovery.NewSubscriber("storage")

	s.Server.RegisterOnShutdown(func() {
		log.Println("on shutdown...")

		var err error
		if err = s.publisher.Unpublished(); err != nil {
			log.Println(err)
		}
		if err = s.subscriber.Unsubscribe(); err != nil {
			log.Println(err)
		}
	})

	s.subscriber.Subscribe()

	return s.Server.ListenAndServe()
}

func (s *Service) Close() {
	s.Server.Shutdown(context.TODO())
}
