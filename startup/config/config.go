package config

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
	return &Config{
		Port:           "8000",      //os.Getenv("GATEWAY_PORT"),
		ProfileHost:    "localhost", //os.Getenv("PROFILE_SERVICE_HOST"),
		ProfilePort:    "8001",      //os.Getenv("PROFILE_SERVICE_PORT"),
		PostHost:       "localhost", //os.Getenv("POST_SERVICE_HOST"),
		PostPort:       "8002",      //os.Getenv("POST_SERVICE_PORT"),
		AuthHost:       "localhost", //os.Getenv("AUTH_SERVICE_HOST"),
		AuthPort:       "8003",      //os.Getenv("AUTH__SERVICE_PORT"),
		ConnectionHost: "localhost",
		ConnectionPort: "8004",
	}
}
