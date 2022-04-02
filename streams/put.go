package streams

import (
	"fmt"
	"io"
	"net/http"
)

// PutStream used to put objects to other service
type PutStream struct {
	writer *io.PipeWriter
	err    chan error
}

func NewPutStream(target string) *PutStream {
	reader, writer := io.Pipe()
	e := make(chan error)

	go func() {
		request, _ := http.NewRequest("PUT", target, reader)
		client := http.Client{}
		r, err := client.Do(request)

		if err == nil && r.StatusCode != http.StatusOK {
			err = fmt.Errorf("%s responsed status code: %d", target, r.StatusCode)
		}
		e <- err
	}()

	return &PutStream{writer: writer, err: e}
}

func (p *PutStream) Write(b []byte) (n int, err error) {
	return p.writer.Write(b)
}

func (p *PutStream) Close() error {
	p.writer.Close()
	return <-p.err
}
