package service

import (
	"context"
	"posts/internal/usecase/kafka"

	pb "posts/internal/pkg/genproto"
	"posts/internal/repository"
)

type LogService struct {
	stg      repository.StorageI
	producer kafka.KafkaProducer
	pb.UnimplementedLogServiceServer
}

func NewLogService(stg repository.StorageI, kafka kafka.KafkaProducer) *LogService {
	return &LogService{stg: stg}
}
func (s *LogService) Create(ctx context.Context, request *pb.LogCreateRequest) (*pb.LogCreateResponse, error) {
	return s.stg.Log().Create(ctx, request)
}
func (s *LogService) GetDetail(ctx context.Context, request *pb.GetId) (*pb.LogGetResponse, error) {
	return s.stg.Log().GetDetail(ctx, request)
}
func (s *LogService) GetList(ctx context.Context, request *pb.FilterLog) (*pb.LogGetAll, error) {
	return s.stg.Log().GetList(ctx, request)
}
func (s *LogService) Update(ctx context.Context, request *pb.LogUpdateRequest) (*pb.LogVoid, error) {
	return s.stg.Log().Update(ctx, request)
}
func (s *LogService) Delete(ctx context.Context, request *pb.GetId) (*pb.LogVoid, error) {
	return s.stg.Log().Delete(ctx, request)
}
