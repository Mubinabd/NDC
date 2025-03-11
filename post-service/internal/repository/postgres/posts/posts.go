package posts

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	          VALUES ($1, $2, $3)`

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
					created_at
	          FROM posts 
			  WHERE id = $1`
	row := r.DB.QueryRowContext(ctx, query, request.Id)

	var post pb.PostGetResponse
	err := row.Scan(
		&post.Id,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *Repository) Update(ctx context.Context, request *pb.PostUpdateRequest) (*pb.PostVoid, error) {
	query := `UPDATE posts SET updated_at = NOW()`
	var args []interface{}
	var updates []string

	if request.UserId != 0 {
		args = append(args, request.UserId)
		updates = append(updates, fmt.Sprintf("user_id = $%d", len(args)))
	}

	if request.Title != "" && request.Title != "string" {
		args = append(args, request.Title)
		updates = append(updates, fmt.Sprintf("title = $%d", len(args)))
	}

	if request.Content != "" && request.Content != "string" {
		args = append(args, request.Content)
		updates = append(updates, fmt.Sprintf("content = $%d", len(args)))
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += ", " + strings.Join(updates, ", ")

	args = append(args, request.Id)
	query += fmt.Sprintf(" WHERE id = $%d", len(args))

	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("updating post: %w", err)
	}
	return &pb.PostVoid{}, nil
}

func (r *Repository) Delete(ctx context.Context, request *pb.GetById) (*pb.PostVoid, error) {
	query := `update posts set deleted_at =Now() WHERE id = $1 `
	_, err := r.DB.ExecContext(ctx, query, request.Id)
	return &pb.PostVoid{}, err
}

func (r *Repository) GetList(ctx context.Context, filter *pb.FilterPost) (*pb.PostGetAll, error) {
	var args []interface{}
	var conditions []string
	argID := 1

	conditions = append(conditions, "p.deleted_at IS NULL")

	if filter.UserId != 0 {
		conditions = append(conditions, fmt.Sprintf("p.user_id = $%d", argID))
		args = append(args, filter.UserId)
		argID++
	}
	if filter.Content != "" {
		conditions = append(conditions, fmt.Sprintf("p.content ILIKE $%d", argID)) // Case-insensitive search
		args = append(args, "%"+filter.Content+"%")                                // Like search
		argID++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Default limit
	var limit, offset int64
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Page > 0 && filter.Limit > 0 {
		offset = (filter.Page - 1) * filter.Limit
	}

	query := fmt.Sprintf(`
		SELECT 
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
		LEFT JOIN users u ON u.id = p.user_id
		%s 
		ORDER BY p.created_at DESC 
		LIMIT $%d OFFSET $%d`, whereClause, argID, argID+1)

	args = append(args, limit, offset)

	log.Println("Generated SQL Query:", query)
	log.Println("Query Args:", args)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying posts: %w", err)
	}
	defer rows.Close()

	var response []*pb.PostGet
	for rows.Next() {
		var post pb.PostGet
		post.User = &pb.User{}
		err := rows.Scan(
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

	countQuery := `SELECT COUNT(p.id) FROM posts p`
	if whereClause != "" {
		countQuery += " " + whereClause
	}

	var totalCount int32
	err = r.DB.QueryRow(countQuery, args[:len(args)-2]...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("counting posts: %w", err)
	}

	return &pb.PostGetAll{
		Post:  response,
		Count: totalCount, // `count` o‘rniga `totalCount`
	}, nil
}
