package stream

import (
	"bytes"
	"fmt"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
)

type Temp struct {
	Server string
	Uuid   string
	url    string
}

func NewTemp(server, name string, size int64) (*Temp, error) {
	tmpUrl := fmt.Sprintf("http://%s/tmp/%s", server, name)
	request, err := http.NewRequest("POST", tmpUrl, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	sender := &Temp{
		Server: server,
		Uuid:   fastjson.GetString(bytes, "data"),
	}

	sender.url = fmt.Sprintf("http://%s/tmp/%s", sender.Server, sender.Uuid)

	return sender, nil
}

// Problem to be solved：
// Large files need to be written many times
// which may cause the storage service to open too many files
func (t *Temp) Write(p []byte) (n int, err error) {
	request, err := http.NewRequest("PATCH", t.url, bytes.NewReader(p))

	if err != nil {
		return n, err
	}

	client := http.Client{}
	r, err := client.Do(request)

	if err != nil {
		return n, err
	}

	if r.StatusCode != http.StatusOK {
		return n, fmt.Errorf("storage service responsed: %d", r.StatusCode)
	}

	return len(p), nil
}

func (t *Temp) Commit(nice bool) {
	method := "DELETE"

	if nice {
		method = "PUT"
	}

	request, _ := http.NewRequest(method, t.url, nil)
	client := http.Client{}
	_, _ = client.Do(request)
}