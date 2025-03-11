package test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	pb "posts/internal/pkg/genproto"
	repository "posts/internal/repository/postgres/posts"
)

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *repository.Repository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %s", err)
	}
	repo := repository.NewRepository(db)
	return db, mock, repo
}

func TestCreatePost(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	request := &pb.PostCreateRequest{
		UserId:  1,
		Title:   "Test Title",
		Content: "Test Content",
	}

	mock.ExpectExec("INSERT INTO posts").
		WithArgs(request.UserId, request.Title, request.Content).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := repo.Create(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, request.UserId, resp.UserId)
	assert.Equal(t, request.Title, resp.Title)
	assert.Equal(t, request.Content, resp.Content)
}

func TestUpdatePost(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	request := &pb.PostUpdateRequest{
		Id:      1,
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	mock.ExpectExec("UPDATE posts SET").
		WithArgs(request.Title, request.Content, request.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := repo.Update(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
