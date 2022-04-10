package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

type decoder struct {
	readers   []io.Reader
	writers   []io.Writer
	enc       reedsolomon.Encoder
	size      int64
	cache     []byte
	cacheSize int
	total     int64
}

func newDecoder(readers []io.Reader, writers []io.Writer, size int64) *decoder {
	enc, _ := reedsolomon.New(NumDataShard, NumParityShard)

	return &decoder{
		readers:   readers,
		writers:   writers,
		enc:       enc,
		size:      size,
		cache:     nil,
		cacheSize: 0,
		total:     0,
	}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	if d.cacheSize == 0 {
		if err = d.read(); err != nil {
			return 0, err
		}
	}

	length := len(p)

	if d.cacheSize < length {
		length = d.cacheSize
	}

	d.cacheSize -= length

	copy(p, d.cache[:length])

	d.cache = d.cache[length:]

	return length, nil
}

func (d *decoder) read() error {
	if d.total == d.size {
		return io.EOF
	}

	shards := make([][]byte, TotalShards)
	repairIds := make([]int, 0)

	var (
		err error
		n   int
	)

	for idx := range shards {
		if d.readers[idx] == nil {
			repairIds = append(repairIds, idx)
			continue
		}

		shards[idx] = make([]byte, BlockPerShard)
		n, err = io.ReadFull(d.readers[idx], shards[idx])

		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			shards[idx] = nil
		} else if n != BlockPerShard {
			shards[idx] = shards[idx][:n]
		}
	}

	err = d.enc.Reconstruct(shards)

	if err != nil {
		return err
	}

	for _, id := range repairIds {
		_, _ = d.writers[id].Write(shards[id])
	}

	for i := 0; i < NumDataShard; i++ {
		shardSize := int64(len(shards[i]))

		if d.total+shardSize > d.size {
			shardSize -= d.total + shardSize - d.size
		}

		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.cacheSize += int(shardSize)
		d.total += shardSize
	}

	return nil
}
