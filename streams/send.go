package streams

import (
	"fmt"
	"io"
	"net/http"
)

// Sender send object to target service
type Sender struct {
	writer *io.PipeWriter
	err    chan error
}

func NewSender(target string) *Sender {
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

	return &Sender{writer: writer, err: e}
}

func (p *Sender) Write(b []byte) (n int, err error) {
	return p.writer.Write(b)
}

func (p *Sender) Close() error {
	p.writer.Close()
	return <-p.err
}
