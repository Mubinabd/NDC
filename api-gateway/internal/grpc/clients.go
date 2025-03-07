package grpc

import (
	"posts/internal/pkg/config"
	pb "posts/internal/pkg/genproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Post pb.PostServiceClient
	Log  pb.LogServiceClient
	User pb.UserServiceClient
}

func NewClients(cfg *config.Config) (*Clients, error) {
	post_conn, err := grpc.NewClient(cfg.GRPCPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	userClient := pb.NewUserServiceClient(post_conn)
	logClient := pb.NewLogServiceClient(post_conn)
	postClient := pb.NewPostServiceClient(post_conn)

	return &Clients{
		Post: postClient,
		Log:  logClient,
		User: userClient,
	}, nil
}
