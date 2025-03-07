package app

import (
	"log"
	"path/filepath"
	"runtime"

	"posts/internal/grpc"
	"posts/internal/http"
	"posts/internal/http/handlers"
	"posts/internal/pkg/config"
	"posts/internal/pkg/kafka"
	"posts/internal/pkg/logger"
	// "github.com/go-redis/redis"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func Run(cfg config.Config) {
	logger := logger.NewLogger(basepath, cfg.LogPath)
	clients, err := grpc.NewClients(&cfg)
	if err != nil {
		logger.ERROR.Println("Failed to create gRPC clients", err)
		log.Println("error")
		return
	}

	//connect to kafka
	broker := []string{cfg.KafkaUrl}
	kafka, err := kafka.NewKafkaProducer(broker)
	if err != nil {
		logger.ERROR.Println("Failed to connect to Kafka", err)
		log.Println("error")
		return
	}
	defer kafka.Close()

	// make handler
	h := handlers.NewHandler(*clients, kafka, logger)

	// make gin
	router := http.NewGin(h)

	// start server
	router.Run(":8080")
}
