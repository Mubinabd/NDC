package test

import (
	"context"
	"regexp"
	"testing"

	pb "posts/internal/pkg/genproto"
	"posts/internal/repository/postgres/logs"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetDetail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := logs.NewRepository(db)
	ctx := context.Background()
	request := &pb.GetId{Id: 1}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, level, message, service_name, created_at FROM logs WHERE id = $1")).
		WithArgs(request.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "level", "message", "service_name", "created_at"}).
			AddRow(1, "INFO", "Test log", "TestService", "2024-03-07 12:00:00"))

	resp, err := repo.GetDetail(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := logs.NewRepository(db)
	ctx := context.Background()
	request := &pb.LogUpdateRequest{
		Id:          1,
		Level:       "ERROR",
		Message:     "Updated log",
		ServiceName: "UpdatedService",
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE logs SET updated_at = NOW(), level = $1, message = $2, service_name = $3 WHERE id = $4")).
		WithArgs(request.Level, request.Message, request.ServiceName, request.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = repo.Update(ctx, request)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := logs.NewRepository(db)
	ctx := context.Background()
	request := &pb.GetId{Id: 1}

	mock.ExpectExec(regexp.QuoteMeta("update logs set deleted_at = NOW() WHERE id = $1")).
		WithArgs(request.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = repo.Delete(ctx, request)
	assert.NoError(t, err)
}
