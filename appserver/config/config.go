package config

import "os"

type Config struct {
	AppName     string `env:"APP_NAME"`
	AppVersion  string `env:"APP_VERSION"`
	BuildType   string `env:"BUILD_TYPE"`
	Node1Url    string `env:"NODE1_URL"`
	Node1ApiKey string `env:"NODE1_API_KEY"`
	Node2Url    string `env:"NODE2_URL"`
	Node2ApiKey string `env:"NODE2_API_KEY"`
}

func GetConfig() *Config {
	return &Config{
		AppName:     os.Getenv("APP_NAME"),
		AppVersion:  os.Getenv("APP_VERSION"),
		BuildType:   os.Getenv("BUILD_TYPE"),
		Node1Url:    os.Getenv("NODE1_URL"),
		Node2Url:    os.Getenv("NODE2_URL"),
		Node1ApiKey: os.Getenv("NODE1_API_KEY"),
		Node2ApiKey: os.Getenv("NODE2_API_KEY"),
	}
}
