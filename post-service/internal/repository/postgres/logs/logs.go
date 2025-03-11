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
	query := `INSERT INTO logs (
								level, 
								message, 
								service_name) 
	          VALUES ($1, $2, $3)`

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
					created_at
	          FROM logs 
			  WHERE id = $1`
	row := r.DB.QueryRowContext(ctx, query, request.Id)

	var log pb.LogGetResponse
	err := row.Scan(
		&log.Id,
		&log.Level,
		&log.Message,
		&log.ServiceName,
		&log.CreatedAt,
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

	if len(conditions) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += ", " + strings.Join(conditions, ", ")

	arg = append(arg, request.Id)
	query += fmt.Sprintf(" WHERE id = $%d", len(arg))

	_, err := r.DB.ExecContext(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	return &pb.LogVoid{}, nil
}

func (r *Repository) Delete(ctx context.Context, request *pb.GetId) (*pb.LogVoid, error) {
	query := `update logs set deleted_at = NOW() WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, request.Id)
	return &pb.LogVoid{}, err
}

func (r *Repository) GetList(ctx context.Context, filter *pb.FilterLog) (*pb.LogGetAll, error) {
	// Asosiy WHERE sharti
	where := `WHERE deleted_at IS NULL`
	var args []interface{}

	// Agar Level berilgan bo‘lsa, shartga qo‘shamiz
	if filter.Level != "" {
		args = append(args, filter.Level)
		where += fmt.Sprintf(" AND level = $%d", len(args))
	}

	// Agar ServiceName berilgan bo‘lsa, shartga qo‘shamiz
	if filter.ServiceName != "" {
		args = append(args, filter.ServiceName)
		where += fmt.Sprintf(" AND service_name = $%d", len(args))
	}

	// Limit va Offset hisoblash
	limit := int64(10) // Default limit
	offset := int64(0) // Default offset

	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Page > 0 && filter.Limit > 0 {
		offset = (filter.Page - 1) * filter.Limit
	}

	args = append(args, limit, offset)

	// Yakuniy SQL so‘rovi
	query := fmt.Sprintf(`
		SELECT 
			COUNT(id) OVER () AS total_count,
			id, 
			level, 
			message, 
			service_name, 
			created_at 
		FROM logs 
		%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	// Queryni ishlatamiz
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.New("querying logs")
	}
	defer rows.Close()

	// Natijalarni o‘qish
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
