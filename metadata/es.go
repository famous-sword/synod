package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"net/http"
	"strings"
)

type (
	// searched response of elasticsearch
	searched struct {
		Hits struct {
			Total struct {
				Value    int    `json:"value"`
				Relation string `json:"relation"`
			}
			Hits []struct {
				Source Meta `json:"_source"`
			}
		}
	}

	ElasticMetaManager struct {
		client    *elasticsearch.Client
		indexName string
	}
)

func (c *ElasticMetaManager) Get(name string, version int) (Meta, error) {
	if version == 0 {
		return c.LatestVersion(name)
	}

	return c.getMetaById(name, version)
}

func (c *ElasticMetaManager) Put(name string, version int, size int64, hash string) error {
	m := Meta{
		Name:    name,
		Version: version,
		Size:    size,
		Hash:    hash,
	}

	b, err := json.Marshal(m)

	if err != nil {
		return err
	}

	request := esapi.IndexRequest{
		Index:      c.indexName,
		DocumentID: generateId(name, version),
		Body:       bytes.NewReader(b),
		OpType:     "create",
	}

	_, err = request.Do(context.Background(), c.client)

	if err != nil {
		return err
	}

	return nil
}

func (c *ElasticMetaManager) Remove(name string, version int) {
	request := esapi.DeleteRequest{
		Index:      c.indexName,
		DocumentID: generateId(name, version),
	}

	_, _ = request.Do(context.Background(), c.client)
}

func (c *ElasticMetaManager) Versions(name string, from, size int) ([]Meta, error) {
	var query string

	if name == "" {
		query = `{"sort":[{"version":{"order":"desc"}}],"from":"%d","size":"%d"}`
		query = fmt.Sprintf(query, from, size)
	} else {
		query = `{"query":{"match":{"name":"%s"}},"sort":[{"version":{"order":"desc"}}],"from":"%d","size":"%d"}`
		query = fmt.Sprintf(query, name, from, size)
	}

	cli := c.client

	response, err := cli.Search(
		cli.Search.WithIndex(c.indexName),
		cli.Search.WithBody(strings.NewReader(query)),
		cli.Search.WithPretty(),
	)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrMetaNotFound
	}

	var result searched

	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	metas := make([]Meta, len(result.Hits.Hits))

	for i, hit := range result.Hits.Hits {
		metas[i] = hit.Source
	}

	return metas, nil
}

func (c *ElasticMetaManager) AddVersion(name, hash string, size int64) error {
	meta, err := c.LatestVersion(name)

	if err != nil && err != ErrMetaNotFound {
		return err
	}

	return c.Put(name, meta.Version+1, size, hash)
}

func (c *ElasticMetaManager) LatestVersion(name string) (meta Meta, err error) {
	query := `{"query":{"match":{"name":"%s"}},"sort":[{"version":{"order":"desc"}}],"size":1}`
	query = fmt.Sprintf(query, name)
	cli := c.client

	response, err := cli.Search(
		cli.Search.WithIndex(c.indexName),
		cli.Search.WithBody(strings.NewReader(query)),
		cli.Search.WithPretty(),
	)

	if err != nil {
		return
	}

	if response.StatusCode == http.StatusNotFound {
		err = ErrMetaNotFound
		return
	}

	defer response.Body.Close()

	var result searched

	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return
	}

	if len(result.Hits.Hits) != 0 {
		meta = result.Hits.Hits[0].Source
	}

	return meta, nil
}

func (c *ElasticMetaManager) getMetaById(name string, version int) (meta Meta, err error) {
	id := generateId(name, version)
	query := `{"query":{"term":{"_id":"%s"}}}`
	query = fmt.Sprintf(query, id)

	cli := c.client

	response, err := cli.Search(
		cli.Search.WithIndex(c.indexName),
		cli.Search.WithBody(strings.NewReader(query)),
		cli.Search.WithPretty(),
	)

	if err != nil {
		return
	}

	if response.StatusCode == http.StatusNotFound {
		err = ErrMetaNotFound
		return
	}

	var result searched

	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return
	}

	if len(result.Hits.Hits) != 0 {
		meta = result.Hits.Hits[0].Source
	}

	return meta, nil
}

func generateId(name string, version int) string {
	return fmt.Sprintf("%s-%d", strings.Trim(name, "/"), version)
}
