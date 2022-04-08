package streams

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

var ErrInvalidServer = errors.New("invalid server address")

// CopyStream pull object from target service
type CopyStream struct {
	reader io.Reader
}

func NewCopyStream(from string) (*CopyStream, error) {
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

	return &CopyStream{reader: response.Body}, err
}

func (r *CopyStream) Read(b []byte) (n int, err error) {
	return r.reader.Read(b)
}
