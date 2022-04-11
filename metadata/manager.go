package metadata

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	"synod/conf"
	"synod/util/logx"
)

var (
	ErrMetaNotFound = errors.New("meta data not found")
)

type Manager interface {
	Get(name string, version int) (Meta, error)
	Put(name string, version int, size int64, hash string) error
	Remove(name string, version int)
	Versions(name string, from, size int) ([]Meta, error)
	AddVersion(name, hash string, size int64) error
	LatestVersion(name string) (Meta, error)
}

func New() Manager {
	return createElasticMetaManager()
}

func createElasticMetaManager() Manager {
	cli, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: conf.StringSlice("meta.elasticsearch.endpoints"),
		Username:  conf.String("meta.elasticsearch.username"),
		Password:  conf.String("meta.elasticsearch.password"),
	})

	if err != nil {
		logx.Fatalw("create meta manager error", "adapter", "elasticsearch", "error", err)
	}

	return &ElasticMetaManager{
		client:    cli,
		indexName: "metas",
	}
}
