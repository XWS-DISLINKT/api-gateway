package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port           string
	ProfileHost    string
	ProfilePort    string
	PostHost       string
	PostPort       string
	AuthHost       string
	AuthPort       string
	ConnectionHost string
	ConnectionPort string
}

func NewConfig() *Config {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		fmt.Println("docker")

		return &Config{
			Port:           os.Getenv("GATEWAY_PORT"),
			ProfileHost:    os.Getenv("PROFILE_SERVICE_HOST"),
			ProfilePort:    os.Getenv("PROFILE_SERVICE_PORT"),
			PostHost:       os.Getenv("POST_SERVICE_HOST"),
			PostPort:       os.Getenv("POST_SERVICE_PORT"),
			AuthHost:       os.Getenv("AUTHENTICATION_SERVICE_HOST"),
			AuthPort:       os.Getenv("AUTHENTICATION_SERVICE_PORT"),
			ConnectionHost: os.Getenv("CONNECTION_SERVICE_HOST"),
			ConnectionPort: os.Getenv("CONNECTION_SERVICE_PORT"),
		}
	} else {
		fmt.Println("local")

		return &Config{
			Port:           "8000",
			ProfileHost:    "localhost",
			ProfilePort:    "8001",
			PostHost:       "localhost",
			PostPort:       "8002",
			AuthHost:       "localhost",
			AuthPort:       "8003",
			ConnectionHost: "localhost",
			ConnectionPort: "8004",
		}
	}
}
