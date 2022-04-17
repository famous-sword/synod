package conf

var schema = map[string]string{
	"app.id":           "SYNOD_APP_ID",
	"app.dev":          "SYNOD_APP_DEBUG",
	"api.addr":         "SYNOD_API_ADDR",
	"storage.addr":     "SYNOD_STORAGE_ADDR",
	"storage.data_dir": "SYNOD_DATA_DIR",
	"storage.temp_dir": "SYNOD_TEMP_DIR",
}

var defaults = map[string]string{
	"api.addr":         ":5555",
	"storage.addr":     ":6666",
	"storage.data_dir": "/data",
	"storage.temp_dir": "/tmp",
}

func (r *Repository) LoadEnvironments() error {
	r.env.AutomaticEnv()
	r.env.SetEnvPrefix(r.envPrefix)

	for key, value := range defaults {
		r.env.SetDefault(key, value)
	}

	for key, name := range schema {
		if err := r.env.BindEnv(key, name); err != nil {
			return err
		}
	}

	return nil
}
