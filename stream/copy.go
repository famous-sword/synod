package stream

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"synod/util/urlbuilder"
)

var ErrInvalidServer = errors.New("invalid server address")

// Copier copy object from storage service
type Copier struct {
	reader io.Reader
}

func NewCopier(server, name string) (*Copier, error) {
	from := urlbuilder.Join(server, "objects", name).Build()

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
