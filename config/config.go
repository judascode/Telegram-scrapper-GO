package config

import "sync"

var (
	APP_ID   = 0
	APP_HASH = ""
)

type Config struct {
	APP_ID   int
	APP_HASH string
}

var configInstance *Config
var once sync.Once

func initializeConfig() {
	cfg := Config{
		APP_ID:   APP_ID,
		APP_HASH: APP_HASH,
	}
	configInstance = &cfg
}

func GetConfig() *Config {
	once.Do(initializeConfig)
	return configInstance
}
