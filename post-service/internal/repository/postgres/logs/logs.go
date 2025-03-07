package logs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	pb "posts/internal/pkg/genproto"
)

type Repository struct {
	*sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{DB: database}
}

func (r *Repository) Create(ctx context.Context, request *pb.LogCreateRequest) (*pb.LogCreateResponse, error) {
	query := `INSERT INTO posts (
								level, 
								message, 
								service_name) 
	          VALUES (?, ?, ?)`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		request.Level,
		request.Message,
		request.ServiceName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	response := pb.LogCreateResponse{
		Level:       request.Level,
		Message:     request.Message,
		ServiceName: request.ServiceName,
	}
	return &response, nil
}

func (r *Repository) GetDetail(ctx context.Context, request *pb.GetId) (*pb.LogGetResponse, error) {
	query := `SELECT 
					id, 
					level, 
					message, 
					service_name, 
					created_at,
					created_by
	          FROM logs 
			  WHERE id = ?`
	row := r.DB.QueryRowContext(ctx, query, request.Id)

	var log pb.LogGetResponse
	err := row.Scan(
		&log.Id,
		&log.Level,
		&log.Message,
		&log.ServiceName,
		&log.CreatedAt,
		&log.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (r *Repository) Update(ctx context.Context, request *pb.LogUpdateRequest) (*pb.LogVoid, error) {
	query := `UPDATE logs SET updated_at = NOW()`

	var arg []interface{}
	var conditions []string

	if request.Level != "" && request.Level != "string" {
		arg = append(arg, request.Level)
		conditions = append(conditions, fmt.Sprintf("level = $%d", len(arg)))
	}

	if request.Message != "" && request.Message != "string" {
		arg = append(arg, request.Message)
		conditions = append(conditions, fmt.Sprintf("message = $%d", len(arg)))
	}

	if request.ServiceName != "" && request.ServiceName != "string" {
		arg = append(arg, request.ServiceName)
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", len(arg)))
	}

	if len(conditions) > 0 {
		query += ", " + strings.Join(conditions, ", ")
	}

	query += fmt.Sprintf(" WHERE id = ?", len(arg)+1)
	arg = append(arg, request.Id)

	_, err := r.DB.ExecContext(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	return &pb.LogVoid{}, nil
}

func (r *Repository) Delete(ctx context.Context, request *pb.GetId) (*pb.LogVoid, error) {
	query := `update logs set deleted_at = EXTRACT(EPOCH FROM NOW) WHERE id = ? and deleted_at is null`
	_, err := r.DB.ExecContext(ctx, query, request.Id)
	return &pb.LogVoid{}, err
}

func (r *Repository) GetList(ctx context.Context, filter *pb.FilterLog) (*pb.LogGetAll, error) {

	where := fmt.Sprintf(`where deleted_at isnull`)
	if filter.Level != "" {
		where += fmt.Sprintf(" and level = '%s'", filter.Level)
	}
	if filter.ServiceName != "" {
		where += fmt.Sprintf(" and service_name = '%s'", filter.ServiceName)
	}

	var limit, offset int64
	if filter.Limit != 0 {
		limit = filter.Limit
	}
	if filter.Page != 0 && filter.Limit != 0 {
		offst := (filter.Page - 1) * (filter.Limit)
		filter.Offset = offst
	}
	if filter.Offset != 0 {
		offset = filter.Offset
	}

	query := fmt.Sprintf(`SELECT 
					COUNT(id) OVER () AS total_count,
					id, 
					level, 
					message, 
					service_name, 
					created_at 
	          FROM logs 
			  %s order by created_at desc %s %s `, where, limit, offset)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.New("querying logs")
	}
	defer rows.Close()

	var response []*pb.LogGetResponse
	var count int32
	for rows.Next() {
		var log pb.LogGetResponse
		err := rows.Scan(
			&count,
			&log.Id,
			&log.Level,
			&log.Message,
			&log.ServiceName,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		response = append(response, &log)
	}

	return &pb.LogGetAll{
		Log:   response,
		Count: count,
	}, nil
}
