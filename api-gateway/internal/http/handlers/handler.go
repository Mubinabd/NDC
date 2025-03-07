package handlers

import (
	"posts/internal/grpc"
	"posts/internal/pkg/kafka"
	"posts/internal/pkg/logger"
	// "github.com/go-redis/redis"
)

type Handler struct {
	Clients  grpc.Clients
	Producer kafka.KafkaProducer
	Logger   *logger.Logger
}

func NewHandler(clients grpc.Clients, producer kafka.KafkaProducer, logger *logger.Logger) *Handler {
	return &Handler{Clients: clients, Producer: producer, Logger: logger}
}
