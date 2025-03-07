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

func TestGetDetailPost(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	request := &pb.GetById{Id: 1}
	expectedPost := pb.PostGetResponse{
		Id:        1,
		UserId:    2,
		Title:     "Sample Post",
		Content:   "Sample Content",
		CreatedAt: "2024-03-07",
		CreatedBy: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "created_at", "created_by"}).
		AddRow(expectedPost.Id, expectedPost.UserId, expectedPost.Title, expectedPost.Content, expectedPost.CreatedAt, expectedPost.CreatedBy)

	mock.ExpectQuery("SELECT (.+) FROM posts WHERE id = ?").
		WithArgs(request.Id).
		WillReturnRows(rows)

	resp, err := repo.GetDetail(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedPost, *resp)
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

func TestGetListPost(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	filter := &pb.FilterPost{UserId: 1, Limit: 10, Page: 1}

	rows := sqlmock.NewRows([]string{"total_count", "id", "user_id", "title", "content", "user_id", "firstname", "lastname", "username", "email", "phone", "gender", "role", "created_at"}).
		AddRow(1, 1, 1, "Title", "Content", 1, "First", "Last", "user1", "user@example.com", "123456", "M", "Admin", "2024-03-07")

	mock.ExpectQuery("SELECT (.+) FROM posts p").
		WillReturnRows(rows)

	resp, err := repo.GetList(context.Background(), filter)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.Count)
	assert.Equal(t, "Title", resp.Post[0].Title)
}
