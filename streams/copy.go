package streams

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

var ErrInvalidServer = errors.New("invalid server address")

// Copier pull object from target service
type Copier struct {
	reader io.Reader
}

func NewCopier(from string) (*Copier, error) {
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

	return &Copier{reader: response.Body}, err
}

func (r *Copier) Read(b []byte) (n int, err error) {
	return r.reader.Read(b)
}
