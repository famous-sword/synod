package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

	f, err := os.Open(TempDir(name))

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

	f, err := os.Create(TempDir(originFileName))

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

func commit(tempName string, tmp *Temp) {
	_ = os.Rename(tempName, Workdir(tmp.Name))
}
