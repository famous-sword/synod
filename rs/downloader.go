package rs

import (
	"io"
	"synod/stream"
	"synod/util/logx"
)

type Downloader struct {
	*decoder
}

type Locates map[int]string

func NewDownloader(locates Locates, servers []string, hash string, size int64) (*Downloader, error) {
	if len(locates)+len(servers) != TotalShards {
		return nil, ErrNumOfServer
	}

	readers := make([]io.Reader, TotalShards)

	for i := 0; i < TotalShards; i++ {
		server := locates[i]

		if server == "" {
			locates[i] = servers[i]
			servers = servers[1:]
			continue
		}

		reader, err := stream.NewCopier(server, getShardSeq(hash, i))

		if err == nil {
			readers[i] = reader
		}
	}

	writers := make([]io.Writer, TotalShards)
	shardSize := (size + NumDataShard - 1) / NumDataShard
	var err error

	for i := range readers {
		writers[i], err = stream.NewTemp(locates[i], getShardSeq(hash, i), shardSize)

		if err != nil {
			return nil, err
		}
	}

	dec := newDecoder(readers, writers, size)

	return &Downloader{dec}, nil
}

func (d *Downloader) Close() {
	for _, writer := range d.writers {
		if writer != nil {
			writer.(*stream.Temp).Commit(true)
		}
	}
}

func (d *Downloader) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		logx.Fatalw("support io.SeekCurrent only")
	}

	if offset < 0 {
		logx.Fatalw("support forward seek only")
	}

	for offset != 0 {
		length := int64(BlockSize)
		if offset < length {
			length = offset
		}
		buffer := make([]byte, length)
		_, _ = io.ReadFull(d, buffer)
		offset -= length
	}

	return offset, nil
}
