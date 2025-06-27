package config

type AppConfig struct {
	Name   string `mapstructure:"name"`
	Env    string `mapstructure:"env"`
	Locale string `mapstructure:"locale"`
}
