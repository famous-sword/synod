package stream

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Puller struct {
	writer *io.PipeWriter
	err    chan error
}

func NewPuller(server, name string) *Puller {
	reader, writer := io.Pipe()
	err := make(chan error)
	b := strings.Builder{}
	b.WriteString("http://")
	b.WriteString(server)
	b.WriteString("/objects/")
	b.WriteString(name)

	go func() {
		request, e := http.NewRequest("PUT", b.String(), reader)

		if e != nil {
			err <- e
			return
		}

		client := http.Client{}
		response, e := client.Do(request)

		if e == nil && response.StatusCode != http.StatusOK {
			e = fmt.Errorf("storage server responsed status: %d", response.StatusCode)
		}

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
