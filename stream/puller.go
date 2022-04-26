package stream

import (
	"fmt"
	"io"
	"net/http"
	"synod/util/urlbuilder"
)

// Puller pull object from storage service
type Puller struct {
	writer *io.PipeWriter
	err    chan error
}

func NewPuller(server, name string) *Puller {
	reader, writer := io.Pipe()
	err := make(chan error)
	url := urlbuilder.Join(server, "objects", name)

	go func() {
		request, e := http.NewRequest("PUT", url.Build(), reader)

		if e != nil {
			err <- e
			return
		}

		client := http.Client{}
		response, e := client.Do(request)

		if e == nil && response.StatusCode != http.StatusOK {
			e = fmt.Errorf("storage server responsed status: %d", response.StatusCode)
		}

		_ = response.Body.Close()

		err <- e
	}()

	return &Puller{
		writer: writer,
		err:    err,
	}
}

func (p Puller) Write(b []byte) (n int, err error) {
	return p.writer.Write(b)
}

func (p Puller) Close() error {
	_ = p.writer.Close()

	return <-p.err
}
