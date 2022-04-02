package streams

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

var ErrInvalidServer = errors.New("invalid server address")

// FetchStream used to fetch objects from other service
type FetchStream struct {
	reader io.Reader
}

func NewFetchStream(from string) (*FetchStream, error) {
	if from == "" {
		return nil, ErrInvalidServer
	}

	response, err := http.Get(from)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responsed status code: %d", from, response.StatusCode)
	}

	return &FetchStream{reader: response.Body}, err
}

func (r *FetchStream) Read(b []byte) (n int, err error) {
	return r.reader.Read(b)
}
