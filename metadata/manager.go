package metadata

import "synod/conf"

type Manager interface {
	Get(name string, version int) (Meta, error)
	Put(name string, version int, size int64, hash string) error
	Remove(name string, version int)
	Versions(name string, from, size int64) ([]Meta, error)
	AddVersion(name, hash string, size int64) error
	LatestVersion(name string) (Meta, error)
}

func New() Manager {
	return &ElasticMetaManager{
		baseURL: conf.String("meta.elasticsearch.endpoint"),
	}
}
