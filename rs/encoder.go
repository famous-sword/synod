package rs

import (
	"github.com/klauspost/reedsolomon"
	"io"
)

type encoder struct {
	writers []io.Writer
	cache   []byte
	enc     reedsolomon.Encoder
}

func newEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(NumDataShard, NumParityShard)

	return &encoder{
		writers: writers,
		cache:   nil,
		enc:     enc,
	}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	current := 0

	for length != 0 {
		next := BlockSize - len(e.cache)

		if next > length {
			next = length
		}

		e.cache = append(e.cache, p[current:current+next]...)

		if len(e.cache) == BlockSize {
			e.Flush()
		}

		current += next
		length -= next
	}

	return len(p), nil
}

func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}

	shards, _ := e.enc.Split(e.cache)
	_ = e.enc.Encode(shards)

	for index := range shards {
		_, _ = e.writers[index].Write(shards[index])
	}

	e.cache = []byte{}
}
