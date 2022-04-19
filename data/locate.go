package data

import (
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"synod/util/logx"
)

var locator *Locator

type Locator struct {
	table *bloom.BloomFilter
	temps map[string]int
	mux   sync.RWMutex
}

func NewLocator() *Locator {
	return &Locator{
		table: bloom.NewWithEstimates(1000000, 0.01),
		temps: map[string]int{},
	}
}

func (l *Locator) LoadToTable() {
	files, err := filepath.Glob(Disk("*"))

	if err != nil {
		logx.Errorw("glob error", "msg", err)
	}

	for i := range files {
		segments := strings.Split(filepath.Base(files[i]), ".")

		if len(segments) != 3 {
			logx.Errorw("error data", "path", files[i])
			continue
		}

		hash := segments[0]
		id, _ := strconv.Atoi(segments[1])

		fmt.Printf("%s=>%d\n", hash, id)

		l.table.AddString(hash)
		l.temps[hash] = id
	}
}

func (l *Locator) Load(v string) {
	l.table.AddString(v)
}

func (l *Locator) Has(v string) bool {
	return l.table.Test([]byte(v))
}

func (l *Locator) AddTemp(hash string, id int) {
	l.mux.Lock()
	l.temps[hash] = id
	l.mux.Unlock()
}

func (l *Locator) HasTemp(hash string) bool {
	_, has := l.temps[hash]

	return has
}

func (l *Locator) TempId(hash string) int {
	if id, has := l.temps[hash]; has {
		return id
	}

	return -1
}

func (l *Locator) Forget(hash string) {
	l.mux.Lock()
	delete(l.temps, hash)
	l.mux.Unlock()
}
