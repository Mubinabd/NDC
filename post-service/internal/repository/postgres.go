package repository

import (
	"context"
	pb "posts/internal/pkg/genproto"
)

type StorageI interface {
	User() UserI
	Log() LogI
	Post() PostI
}
type LogI interface {
	Create(ctx context.Context, request *pb.LogCreateRequest) (*pb.LogCreateResponse, error)
	GetDetail(ctx context.Context, id *pb.GetId) (*pb.LogGetResponse, error)
	GetList(ctx context.Context, req *pb.FilterLog) (*pb.LogGetAll, error)
	Update(ctx context.Context, req *pb.LogUpdateRequest) (*pb.LogVoid, error)
	Delete(ctx context.Context, Id *pb.GetId) (*pb.LogVoid, error)
}
type PostI interface {
	Create(ctx context.Context, request *pb.PostCreateRequest) (*pb.PostCreateResponse, error)
	GetDetail(ctx context.Context, id *pb.GetById) (*pb.PostGetResponse, error)
	GetList(ctx context.Context, req *pb.FilterPost) (*pb.PostGetAll, error)
	Update(ctx context.Context, req *pb.PostUpdateRequest) (*pb.PostVoid, error)
	Delete(ctx context.Context, id *pb.GetById) (*pb.PostVoid, error)
}
type UserI interface {
	Create(ctx context.Context, request *pb.UserCreateRequest) (*pb.UserCreateResponse, error)
	GetDetail(ctx context.Context, id *pb.ById) (*pb.UserGetResponse, error)
	GetList(ctx context.Context, req *pb.FilterUser) (*pb.UserGetAll, error)
	Update(ctx context.Context, req *pb.UserUpdateRequest) (*pb.Void, error)
	Delete(ctx context.Context, id *pb.ById) (*pb.Void, error)
	Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error)
	ChangeUserPassword(ctx context.Context, request *pb.UserRecoverPasswordRequest) (*pb.Void, error)
}
