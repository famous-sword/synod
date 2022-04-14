package conf

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
)

var (
	ErrMustYaml       = errors.New("configuration file must be in YAML")
	ErrConfigNotFound = errors.New("config file not found")
)

func Startup(config string) error {
	if config == "" {
		return scanConfigAtDefaultPaths()
	}

	return loadConfigFrom(config)
}

func loadConfigFrom(config string) error {
	real, err := filepath.Abs(config)

	if err != nil {
		return err
	}

	if _, err := os.Stat(real); errors.Is(err, os.ErrNotExist) {
		return ErrConfigNotFound
	}

	ext := filepath.Ext(real)

	if ext != ".yml" && ext != ".yaml" {
		return ErrMustYaml
	}

	viper.SetConfigFile(real)

	return viper.ReadInConfig()
}

func scanConfigAtDefaultPaths() error {
	viper.SetConfigName("synod.yml")

	viper.SetConfigType("yaml")
	viper.AddConfigPath("./var/")
	viper.AddConfigPath("var/")
	viper.AddConfigPath("/etc/synod/")

	current, err := user.Current()

	if err != nil {
		viper.AddConfigPath(current.HomeDir)
	}

	return viper.ReadInConfig()
}
