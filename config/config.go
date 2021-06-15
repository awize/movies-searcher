package config

const ConfigFile = "config.yml"

type Config struct {
	Port           int         `mapstructure:"api_port"`
	MoviesAPI      ExternalAPI `mapstructure:"movies_api"`
	MoviesFileName string      `mapstructure:"movies_filename"`
}

type ExternalAPI struct {
	BaseUrl  string            `mapstructure:"base_url"`
	ApiKey   string            `mapstructure:"api_key"`
	Defaults map[string]string `mapstructure:"defaults"`
}
