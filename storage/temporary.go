package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Temp struct {
	Uuid string
	Name string
	Size int64
}

func ofUuid(u string) (*Temp, error) {
	name := u + ".json"

	f, err := os.Open(withTemp(name))

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
	originFileName := t.Uuid + ".json"

	f, err := os.Create(withTemp(originFileName))

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
	_ = os.Rename(tempName, withWorkdir(tmp.Name))
}
