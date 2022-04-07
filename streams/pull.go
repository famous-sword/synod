package streams

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

var ErrInvalidServer = errors.New("invalid server address")

// Puller pull object from target service
type Puller struct {
	reader io.Reader
}

func NewPuller(from string) (*Puller, error) {
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

	return &Puller{reader: response.Body}, err
}

func (r *Puller) Read(b []byte) (n int, err error) {
	return r.reader.Read(b)
}
