package conf

import (
	"github.com/spf13/viper"
	"os/user"
)

func Startup() error {
	return loadDefaultConfig()
}

func loadDefaultConfig() error {
	viper.SetConfigName("config.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./var/")
	viper.AddConfigPath("var/")
	viper.AddConfigPath("/etc/synod/")

	current, err := user.Current()

	if err != nil {
		viper.AddConfigPath(current.HomeDir)
	}

	err = viper.ReadInConfig()

	if err != nil {
		return err
	}

	return nil
}
