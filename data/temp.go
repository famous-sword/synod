package data

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"synod/util"
)

const (
	extInfo = ".json"
	extTemp = ".tmp"
)

type Temp struct {
	Uuid string
	Name string
	Size int64
}

func ofUuid(u string) (*Temp, error) {
	name := u + extInfo

	f, err := os.Open(TempPath(name))

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(f)

	if err != nil {
		return nil, err
	}

	var tmp Temp

	if err = json.Unmarshal(bytes, &tmp); err != nil {
		return nil, err
	}

	return &tmp, nil
}

func (t *Temp) saveInfo() error {
	originFileName := t.Uuid + extInfo

	f, err := os.Create(TempPath(originFileName))

	if err != nil {
		return err
	}

	defer f.Close()

	bytes, err := json.Marshal(t)

	if err != nil {
		return err
	}

	_, err = f.Write(bytes)

	if err != nil {
		return err
	}

	return nil
}

func (t *Temp) hash() string {
	segments := strings.Split(t.Name, ".")

	return segments[0]
}

func (t *Temp) id() int {
	segments := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(segments[1])

	return id
}

func commit(tempName string, tmp *Temp) {
	file, _ := os.Open(tempName)
	tempHash := url.PathEscape(util.SumHash(file))
	_ = file.Close()

	_ = os.Rename(tempName, Disk(tmp.Name+"."+tempHash))

	locator.AddTemp(tmp.hash(), tmp.id())
}
