package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type (
	// searched response of elasticsearch
	searched struct {
		Hits struct {
			Total int
			Hits  []struct {
				Source Meta `json:"_source"`
			}
		}
	}

	ElasticMetaManager struct {
		baseURL string
	}
)

func (c *ElasticMetaManager) Get(name string, version int) (Meta, error) {
	if version == 0 {
		return c.LatestVersion(name)
	}

	return c.getMeta(name, version)
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

	u := c.getUrlBuilder(makeMetaDocName(name, version))
	queries := &url.Values{}
	queries.Add("op_type", "create")
	u.RawPath = queries.Encode()

	client := &http.Client{}
	request, err := http.NewRequest("PUT", u.String(), bytes.NewReader(b))

	if err != nil {
		return err
	}

	r, err := client.Do(request)

	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("elasticsearch responsed status: %d", r.StatusCode)
	}

	return nil
}

func (c *ElasticMetaManager) Remove(name string, version int) {
	u := c.getUrlBuilder(makeMetaDocName(name, version))
	request, _ := http.NewRequest("DELETE", u.String(), nil)
	client := &http.Client{}

	_, _ = client.Do(request)
}

func (c *ElasticMetaManager) Versions(name string, from, size int64) ([]Meta, error) {
	u := c.getUrlBuilder("_search")
	queries := &url.Values{}
	queries.Add("sort", "name,version")
	queries.Add("form", strconv.FormatInt(from, 10))
	queries.Add("size", strconv.FormatInt(size, 10))

	if name != "" {
		queries.Add("q", "name:"+name)
	}

	u.RawQuery = queries.Encode()

	response, err := http.Get(u.String())

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	metas := make([]Meta, 0)
	var results searched
	err = json.Unmarshal(body, &results)

	if err != nil {
		return nil, err
	}

	for i := range results.Hits.Hits {
		metas = append(metas, results.Hits.Hits[i].Source)
	}

	return metas, nil
}

func (c *ElasticMetaManager) AddVersion(name, hash string, size int64) error {
	meta, err := c.LatestVersion(name)

	if err != nil {
		return err
	}

	return c.Put(name, meta.Version+1, size, hash)
}

func (c *ElasticMetaManager) LatestVersion(name string) (meta Meta, err error) {
	u := c.getUrlBuilder("_search")

	queries := url.Values{}
	queries.Add("name", name)
	queries.Add("size", "1")
	queries.Add("sort", "version:desc")

	u.RawQuery = queries.Encode()

	response, err := http.Get(u.String())

	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to ger latest version by responsed %d", response.StatusCode)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	var r searched

	err = json.Unmarshal(body, &r)

	if err != nil {
		return
	}

	if len(r.Hits.Hits) != 0 {
		meta = r.Hits.Hits[0].Source
	}

	return
}

func (c *ElasticMetaManager) getUrlBuilder(op string) *url.URL {
	base := fmt.Sprintf("%s/metas/%s", c.baseURL, op)
	u, _ := url.Parse(base)

	return u
}

func (c *ElasticMetaManager) getMeta(name string, version int) (meta Meta, err error) {
	u := c.getUrlBuilder(makeMetaDocName(name, version) + "_source")
	response, err := http.Get(u.String())

	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("elasticsearch responsed status: %d", response.StatusCode)
		return
	}

	r, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(r, &meta)

	return
}

func makeMetaDocName(name string, version int) string {
	return fmt.Sprintf("/objects/%s-%d", name, version)
}
