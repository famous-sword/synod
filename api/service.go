package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"synod/conf"
	"synod/discovery"
	"synod/metadata"
	"synod/util/logx"
)

var (
	ErrInvalidAddr    = errors.New("invalid addr")
	ErrHashRequired   = errors.New("required object hash in digest header")
	ErrInvalidName    = errors.New("invalid name")
	ErrNoPeer         = errors.New("no peer available")
	ErrNotEnoughPeers = errors.New("cannot find enough peers")
	ErrHashCheckFail  = errors.New("hash check failed")
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
	return &Service{
		Name: "api",
		Addr: conf.String("api.addr"),
	}
}

func (s *Service) Run() error {
	handler := gin.Default()
	handler.GET("/objects/:name", s.load)
	handler.PUT("/objects/:name", s.put)
	handler.DELETE("/objects/:name", s.versions)
	handler.GET("/versions/:name", s.versions)

	s.Server = &http.Server{
		Handler: handler,
	}

	s.metaManager = metadata.New()

	if s.Addr == "" {
		return ErrInvalidAddr
	}

	s.Server.Addr = s.Addr
	s.publisher = discovery.NewPublisher(s.Name, s.Addr)
	if err := s.publisher.Publish(); err != nil {
		return err
	}

	s.subscriber = discovery.NewSubscriber("data")
	if err := s.subscriber.Subscribe(); err != nil {
		return err
	}

	return s.Server.ListenAndServe()
}

func (s *Service) Shutdown() {
	var err error
	if err = s.publisher.Unpublished(); err != nil {
		logx.Errorw("unpublished service error", "error", err)
	}
	if err = s.subscriber.Unsubscribe(); err != nil {
		logx.Errorw("unsubscribe service error", "error", err)
	}
}
