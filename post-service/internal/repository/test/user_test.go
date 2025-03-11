package test

import (
	"context"
	us "posts/internal/repository/postgres/users"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	pb "posts/internal/pkg/genproto"
)

func TestNewRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	repo := us.NewRepository(db)
	assert.NotNil(t, repo)
}

func TestCreateUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := us.NewRepository(db)
	ctx := context.Background()

	mock.ExpectExec("INSERT INTO users").
		WithArgs("John", "Doe", "johndoe", "john@example.com", "+1234567890", "M", sqlmock.AnyArg(), "ADMIN").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := &pb.UserCreateRequest{
		FirstName:   "John",
		LastName:    "Doe",
		Username:    "johndoe",
		Email:       "john@example.com",
		PhoneNumber: "+1234567890",
		Gender:      "M",
		Password:    "password123",
		Role:        "ADMIN",
	}

	resp, err := repo.Create(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.FirstName, resp.FirstName)
}

func TestUpdateUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := us.NewRepository(db)
	ctx := context.Background()

	mock.ExpectExec("UPDATE users").
		WithArgs("John", "Doe", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := &pb.UserUpdateRequest{
		Id:        1,
		FirstName: "John",
		LastName:  "Doe",
	}

	_, err := repo.Update(ctx, req)
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := us.NewRepository(db)
	ctx := context.Background()

	mock.ExpectExec("UPDATE users SET deleted_at").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := &pb.ById{Id: 1}
	_, err := repo.Delete(ctx, req)
	assert.NoError(t, err)
}
