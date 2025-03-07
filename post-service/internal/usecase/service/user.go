package service

import (
	"context"

	pb "posts/internal/pkg/genproto"
	"posts/internal/repository"
)

type UserService struct {
	stg repository.StorageI
	pb.UnimplementedUserServiceServer
}

func NewUserService(stg repository.StorageI) *UserService {
	return &UserService{stg: stg}
}
func (s *UserService) Create(ctx context.Context, request *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	return s.stg.User().Create(ctx, request)
}
func (s *UserService) GetDetail(ctx context.Context, request *pb.ById) (*pb.UserGetResponse, error) {
	return s.stg.User().GetDetail(ctx, request)
}
func (s *UserService) GetList(ctx context.Context, request *pb.FilterUser) (*pb.UserGetAll, error) {
	return s.stg.User().GetList(ctx, request)
}
func (s *UserService) Update(ctx context.Context, request *pb.UserUpdateRequest) (*pb.Void, error) {
	return s.stg.User().Update(ctx, request)
}
func (s *UserService) Delete(ctx context.Context, request *pb.ById) (*pb.Void, error) {
	return s.stg.User().Delete(ctx, request)
}
func (s *UserService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.stg.User().Login(ctx, request)
}
func (s *UserService) ChangeUserPassword(ctx context.Context, request *pb.UserRecoverPasswordRequest) (*pb.Void, error) {
	return s.stg.User().ChangeUserPassword(ctx, request)
}
