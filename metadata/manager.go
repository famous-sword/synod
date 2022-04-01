package metadata

import "synod/conf"

type Manager interface {
	LatestVersion(name string) (Meta, error)
	Get(name string, version int) (Meta, error)
	Put(name string, version int, size int64, hash string) error
	AddVersion(name, hash string, size int64) error
	Versions(name string, from, size int64) ([]Meta, error)
	Remove(name string, version int)
}

func New() Manager {
	return &ElasticMetaManager{
		baseURL: conf.String("meta.elasticsearch.endpoint"),
	}
}
