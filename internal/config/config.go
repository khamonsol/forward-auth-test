package config

import (
	"log"
	"os"
)

type ServerConfig struct {
	Address string
}

type Config struct {
	Server ServerConfig
}

func LoadConfig() (*Config, error) {
	address := getEnv("SERVER_ADDRESS", ":8080")

	config := &Config{
		Server: ServerConfig{
			Address: address,
		},
	}

	log.Printf("Config loaded: %+v", config)
	return config, nil
}

// getEnv reads an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
