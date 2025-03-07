package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Context struct {
		Timeout string
	}
	GRPCPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	KafkaUrl         string
}

func New() *Config {
	var config Config

	// Yaml faylni ochish
	file, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	// YAML faylni o'qish
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// General configuration
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// PostgreSQL configuration
	config.PostgresHost = cast.ToString(getEnv("POSTGRES_HOST", "localhost"))
	config.PostgresPort = cast.ToString(getEnv("POSTGRES_PORT", "5432"))
	config.PostgresUser = cast.ToString(getEnv("POSTGRES_USER", "postgres"))
	config.PostgresPassword = cast.ToString(getEnv("POSTGRES_PASSWORD", "1"))
	config.PostgresDatabase = cast.ToString(getEnv("POSTGRES_DATABASE", "posts"))
	config.KafkaUrl = cast.ToString(getEnv("KAFKA_URL", ":9092"))
	config.GRPCPort = cast.ToString(getEnv("GRPC_PORT", ":7001"))

	return &config
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
