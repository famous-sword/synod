package conf

var schema = map[string]string{
	"app.id":        "SYNOD_APP_ID",
	"app.dev":       "SYNOD_APP_DEBUG",
	"api.addr":      "SYNOD_API_ADDR",
	"data.addr":     "SYNOD_DATA_ADDR",
	"data.data_dir": "SYNOD_DATA_DIR",
	"data.temp_dir": "SYNOD_TEMP_DIR",
}

var defaults = map[string]string{
	"api.addr":      ":5555",
	"data.addr":     ":6666",
	"data.data_dir": "/data",
	"data.temp_dir": "/tmp",
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
