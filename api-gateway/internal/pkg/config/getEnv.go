package config

import (
	"log"
	"os"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPCPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	LogPath          string
	KafkaUrl         string
	HttpPort         string

	DefaultOffset string
	DefaultLimit  string
}

func Load() Config {

	config := Config{}

	file, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// PostgreSQL configuration
	config.PostgresHost = cast.ToString(getEnv("POSTGRES_HOST", "localhost"))
	config.PostgresPort = cast.ToString(getEnv("POSTGRES_PORT", "5432"))
	config.PostgresUser = cast.ToString(getEnv("POSTGRES_USER", "postgres"))
	config.PostgresPassword = cast.ToString(getEnv("POSTGRES_PASSWORD", "1"))
	config.PostgresDatabase = cast.ToString(getEnv("POSTGRES_DATABASE", "posts"))
	config.KafkaUrl = cast.ToString(getEnv("KAFKA_URL", "kafka_posts:9092"))
	config.GRPCPort = cast.ToString(getEnv("GRPC_PORT", "posts_service:7001"))
	config.HttpPort = cast.ToString(getEnv("HTTP_PORT", ":8080"))

	config.DefaultOffset = cast.ToString(getEnv("DEFAULT_OFFSET", "0"))
	config.DefaultLimit = cast.ToString(getEnv("DEFAULT_LIMIT", "10"))

	return config
}

func getEnv(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
