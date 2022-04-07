package metadata

import (
	"github.com/elastic/go-elasticsearch/v7"
	"testing"
)

func TestLastVersion(t *testing.T) {
	m := createTestingMetaManager()

	md, err := m.LatestVersion("/01.gif")

	if err != nil {
		t.Errorf("search error: %v", err)
	}

	t.Log("got meta: ", md)
}

func TestElasticMetaManager_Versions(t *testing.T) {
	c := createTestingMetaManager()

	got, err := c.Versions("", 0, 10)

	if err != nil {
		t.Fatalf("error: %v", err)
	}

	t.Log("got meta: ", got)
}

func TestElasticMetaManager_getMetaById(t *testing.T) {
	m := createTestingMetaManager()
	meta, err := m.getMetaById("01.gif", 1)

	if err != nil {
		t.Fatal(err)
	}

	t.Log("got meta:", meta)
}

func TestElasticMetaManager_Versions1(t *testing.T) {
	m := createTestingMetaManager()
	metas, err := m.Versions("01.gif", 0, 10)

	if err != nil {
		t.Fatal(err)
	}

	t.Log("got metas: ", metas)
}

func createTestingMetaManager() *ElasticMetaManager {
	cli, _ := elasticsearch.NewDefaultClient()

	return &ElasticMetaManager{
		client:    cli,
		indexName: "metas",
	}
}
