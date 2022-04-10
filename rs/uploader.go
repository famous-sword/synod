package rs

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"synod/streams"
)

var (
	ErrNumOfServer = errors.New("incorrect number of servers")
)

type Uploader struct {
	*encoder
}

func NewUploader(servers []string, hash string, size int64) (*Uploader, error) {
	if len(servers) != TotalShards {
		return nil, ErrNumOfServer
	}

	shard := (size + NumDataShard - 1) / NumDataShard
	writers := make([]io.Writer, TotalShards)
	var err error

	for i := range writers {
		serial := fmt.Sprintf("%s.%d", hash, i)
		writers[i], err = streams.NewTempStream(servers[i], serial, shard)

		if err != nil {
			return nil, err
		}
	}

	enc := newEncoder(writers)

	return &Uploader{enc}, err
}

func (u *Uploader) Commit(success bool) {
	u.Flush()

	for _, writer := range u.writers {
		writer.(*streams.TempStream).Commit(success)
	}
}
