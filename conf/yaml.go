package conf

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func (r *Repository) ReadYaml(path string) error {
	if path == "" {
		return r.readFromDefault()
	}

	r.yamlPath = path

	return r.readFromSpecifiedFile(path)
}

func (r *Repository) readFromDefault() error {
	r.yaml.SetConfigName("synod.yml")

	r.yaml.SetConfigType("yaml")
	r.yaml.AddConfigPath("./var/")
	r.yaml.AddConfigPath("var/")
	r.yaml.AddConfigPath("/etc/synod/")

	return r.yaml.ReadInConfig()
}

func (r *Repository) readFromSpecifiedFile(path string) error {
	realPath, err := filepath.Abs(path)

	if err != nil {
		return err
	}

	if _, err := os.Stat(realPath); errors.Is(err, os.ErrNotExist) {
		return ErrConfigNotFound
	}

	ext := filepath.Ext(realPath)

	if ext != ".yml" && ext != ".yaml" {
		return ErrMustYaml
	}

	r.yaml.SetConfigFile(realPath)

	return r.yaml.ReadInConfig()
}
