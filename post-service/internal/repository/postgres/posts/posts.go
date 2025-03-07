package posts

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

func (r *Repository) Create(ctx context.Context, request *pb.PostCreateRequest) (*pb.PostCreateResponse, error) {
	query := `INSERT INTO posts (
								user_id, 
								title, 
								content) 
	          VALUES (?, ?, ?)`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		request.UserId,
		request.Title,
		request.Content,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	response := pb.PostCreateResponse{
		UserId:  request.UserId,
		Title:   request.Title,
		Content: request.Content,
	}
	return &response, nil
}

func (r *Repository) GetDetail(ctx context.Context, request *pb.GetById) (*pb.PostGetResponse, error) {
	query := `SELECT 
					id, 
					user_id, 
					title, 
					content, 
					created_at,
					created_by
	          FROM posts 
			  WHERE id = ?`
	row := r.DB.QueryRowContext(ctx, query, request.Id)

	var post pb.PostGetResponse
	err := row.Scan(
		&post.Id,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *Repository) Update(ctx context.Context, request *pb.PostUpdateRequest) (*pb.PostVoid, error) {
	query := `UPDATE posts SET updated_at = NOW()`

	var arg []interface{}
	var conditions []string

	if request.UserId != 0 {
		arg = append(arg, request.UserId)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(arg)))
	}

	if request.Title != "" && request.Title != "string" {
		arg = append(arg, request.Title)
		conditions = append(conditions, fmt.Sprintf("title = $%d", len(arg)))
	}

	if request.Content != "" && request.Content != "string" {
		arg = append(arg, request.Content)
		conditions = append(conditions, fmt.Sprintf("content = $%d", len(arg)))
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
	return &pb.PostVoid{}, nil
}

func (r *Repository) Delete(ctx context.Context, request *pb.GetById) (*pb.PostVoid, error) {
	query := `update posts set deleted_at = EXTRACT(EPOCH FROM NOW) WHERE id = ? and deleted_at is null`
	_, err := r.DB.ExecContext(ctx, query, request.Id)
	return &pb.PostVoid{}, err
}

func (r *Repository) GetList(ctx context.Context, filter *pb.FilterPost) (*pb.PostGetAll, error) {

	where := fmt.Sprintf(`where p.deleted_at isnull`)
	if filter.UserId != 0 {
		where += fmt.Sprintf(" and p.user_id = '%d'", filter.UserId)
	}
	if filter.Content != "" {
		where += fmt.Sprintf(" and p.content = '%s'", filter.Content)
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
					COUNT(p.id) OVER () AS total_count,
					p.id, 
					p.user_id, 
					p.title, 
					p.content, 
					u.id as user_id,
					u.firstname as firstname,
					u.lastname as lastname,
					u.username as username,
					u.email as email,
					u.phone as phone,
					u.gender,
					u.role,
					p.created_at 
	          FROM posts p 
	          left join users u on u.id = p.user_id
			  %s order by created_at desc %s %s `, where, limit, offset)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.New("querying users")
	}
	defer rows.Close()

	var response []*pb.PostGet
	var count int32
	for rows.Next() {
		var post pb.PostGet
		post.User = &pb.User{}
		err := rows.Scan(
			&count,
			&post.Id,
			&post.UserId,
			&post.Title,
			&post.Content,
			&post.User.Id,
			&post.User.FirstName,
			&post.User.LastName,
			&post.User.Username,
			&post.User.Email,
			&post.User.PhoneNumber,
			&post.User.Gender,
			&post.User.Role,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		response = append(response, &post)
	}

	return &pb.PostGetAll{
		Post:  response,
		Count: count,
	}, nil
}
