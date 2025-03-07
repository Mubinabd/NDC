package service

import (
	"context"
	"posts/internal/usecase/kafka"

	pb "posts/internal/pkg/genproto"
	"posts/internal/repository"
)

type PostService struct {
	stg      repository.StorageI
	producer kafka.KafkaProducer
	pb.UnimplementedPostServiceServer
}

func NewPostService(stg repository.StorageI, kafka kafka.KafkaProducer) *PostService {
	return &PostService{stg: stg}
}
func (s *PostService) Create(ctx context.Context, request *pb.PostCreateRequest) (*pb.PostCreateResponse, error) {
	return s.stg.Post().Create(ctx, request)
}
func (s *PostService) GetDetail(ctx context.Context, request *pb.GetById) (*pb.PostGetResponse, error) {
	return s.stg.Post().GetDetail(ctx, request)
}
func (s *PostService) GetList(ctx context.Context, request *pb.FilterPost) (*pb.PostGetAll, error) {
	return s.stg.Post().GetList(ctx, request)
}
func (s *PostService) Update(ctx context.Context, request *pb.PostUpdateRequest) (*pb.PostVoid, error) {
	return s.stg.Post().Update(ctx, request)
}
func (s *PostService) Delete(ctx context.Context, request *pb.GetById) (*pb.PostVoid, error) {
	return s.stg.Post().Delete(ctx, request)
}
