package app

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"posts/internal/pkg/config"
	pb "posts/internal/pkg/genproto"
	"posts/internal/pkg/postgres"
	"posts/internal/repository/postgres"
	"posts/internal/usecase/kafka"
	"posts/internal/usecase/service"
)

func Run(cf *config.Config) {
	// connect to postgres
	pgm, err := postgres.NewPostgresStorage(cf)
	cfg := config.New()
	log.Printf("Connecting to DB: %s:%s", cfg.PostgresHost, cfg.PostgresPort)

	// connect to kafka producer
	kf_p, err := kafka.NewKafkaProducer([]string{cf.KafkaUrl})
	if err != nil {
		log.Fatal(err)
	}
	// repo
	if pgm == nil || pgm.DB == nil {
		log.Fatal("Postgres connection is nil")
	}
	db := repo.NewStorage(pgm.DB)

	//db := repo.NewStorage(pgm.DB)
	log.Println("Connected to database")

	// register kafka handlers
	k_handler := KafkaHandler{
		log:  service.NewLogService(db, kf_p),
		post: service.NewPostService(db, kf_p),
	}

	if err := Registries(&k_handler, cf); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", cf.GRPCPort)
	if cf.GRPCPort == "" {
		log.Fatal("GRPC port is not set in config")
	}

	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}
	// set grpc server
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, service.NewUserService(db))
	pb.RegisterLogServiceServer(server, service.NewLogService(db, kf_p))
	pb.RegisterPostServiceServer(server, service.NewPostService(db, kf_p))

	// start server

	log.Println("Server started on", cf.GRPCPort)
	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	defer lis.Close()
}
