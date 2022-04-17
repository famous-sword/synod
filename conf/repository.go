package conf

import (
	"github.com/spf13/viper"
)

var repository *Repository

type Repository struct {
	env       *viper.Viper
	envPrefix string
	yaml      *viper.Viper
	yamlPath  string
}

func New() *Repository {
	return &Repository{
		env:       viper.New(),
		envPrefix: "SYNOD",
		yaml:      viper.New(),
	}
}

func (r *Repository) String(key string) string {
	if r.env.IsSet(key) {
		return r.env.GetString(key)
	}

	return r.yaml.GetString(key)
}

func (r *Repository) StringSlice(key string) []string {
	return r.yaml.GetStringSlice(key)
}

func (r *Repository) IntSlice(key string) []int {
	return r.yaml.GetIntSlice(key)
}

func (r *Repository) Integer(key string) int {
	if r.env.IsSet(key) {
		return r.env.GetInt(key)
	}

	return r.yaml.GetInt(key)
}

func (r *Repository) Bool(key string) bool {
	if r.env.IsSet(key) {
		return r.env.GetBool(key)
	}

	return r.yaml.GetBool(key)
}

func (r *Repository) Int32(key string) int32 {
	if r.env.IsSet(key) {
		r.env.GetInt32(key)
	}

	return r.yaml.GetInt32(key)
}

func (r *Repository) Int64(key string) int64 {
	if r.env.IsSet(key) {
		r.env.GetInt64(key)
	}

	return r.yaml.GetInt64(key)
}

func (r *Repository) Uint(key string) uint {
	if r.env.IsSet(key) {
		r.env.GetUint(key)
	}

	return r.yaml.GetUint(key)
}

func (r *Repository) Uint32(key string) uint32 {
	if r.env.IsSet(key) {
		r.env.GetUint32(key)
	}

	return r.yaml.GetUint32(key)
}

func (r *Repository) Uint64(key string) uint64 {
	if r.env.IsSet(key) {
		r.env.GetUint64(key)
	}

	return r.yaml.GetUint64(key)
}

func (r *Repository) Float(key string) float64 {
	if r.env.IsSet(key) {
		r.env.GetFloat64(key)
	}

	return r.yaml.GetFloat64(key)
}

func (r *Repository) Set(key string, value interface{}) {
	r.env.Set(key, value)
}
