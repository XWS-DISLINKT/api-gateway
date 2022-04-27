package config

import "os"

type Config struct {
	Port        string
	ProfileHost string
	ProfilePort string
}

func NewConfig() *Config {
	return &Config{
		Port:        os.Getenv("GATEWAY_PORT"),
		ProfileHost: os.Getenv("PROFILE_SERVICE_HOST"),
		ProfilePort: os.Getenv("PROFILE_SERVICE_PORT"),
	}
}
