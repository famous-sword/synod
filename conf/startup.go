package conf

import (
	"github.com/pkg/errors"
)

var (
	ErrMustYaml       = errors.New("configuration file must be in YAML")
	ErrConfigNotFound = errors.New("repository file not found")
)

func Startup(config string) (err error) {
	repository = New()

	if err = repository.ReadYaml(config); err != nil {
		return err
	}

	if err = repository.LoadEnvironments(); err != nil {
		return err
	}

	return nil
}
