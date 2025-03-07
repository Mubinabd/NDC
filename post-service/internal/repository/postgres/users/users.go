package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	pb "posts/internal/pkg/genproto"
	"strings"
)

type Repository struct {
	*sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{DB: database}
}

func (u *Repository) Create(ctx context.Context, request *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("hashing password")
	}
	hashedPassword := string(hash)

	role := strings.ToUpper(strings.TrimSpace(request.Role))
	if role != "ADMIN" && role != "VIEWER" {
		return nil, errors.New("role should be ADMIN or VIEWER")
	}

	var gen string
	if request.Gender != "" {
		if (strings.ToUpper(request.Gender) != "M") && (strings.ToUpper(request.Gender) != "F") {
			return nil, errors.New("incorrect gender. gender should be M (male) or F (female)")
		}
		gen = strings.ToUpper(request.Gender)
	}

	query := `
		INSERT INTO users (
		                   first_name, 
		                   last_name, 
		                   username, 
		                   email, 
		                   phone, 
		                   gender, 
		                   password, 
		                   role)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = u.DB.ExecContext(
		ctx,
		query,
		request.FirstName,
		request.LastName,
		request.Username,
		request.Email,
		request.PhoneNumber,
		gen,
		hashedPassword,
		request.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	response := pb.UserCreateResponse{
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		Username:    request.Username,
		Gender:      request.Gender,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		Password:    hashedPassword,
		Role:        request.Role,
	}

	return &response, nil
}

func (u *Repository) GetDetail(ctx context.Context, request *pb.ById) (*pb.UserGetResponse, error) {
	query := `
		SELECT 
		     id,
		     first_name, 
		     last_name, 
		     username, 
		     email, 
		     phone, 
		     gender, 
		     role,
		     created_at,
		     created_by
		FROM users 
			WHERE id = ? and deleted_at = 0
	`
	var user pb.UserGetResponse

	err := u.DB.QueryRowContext(ctx, query, request.Id).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.Gender,
		&user.Role,
		&user.CreatedAt,
		&user.CreatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (u *Repository) Update(ctx context.Context, request *pb.UserUpdateRequest) (*pb.Void, error) {

	query := `UPDATE users SET updated_at = NOW()`

	var arg []interface{}
	var conditions []string

	if request.FirstName != "" && request.FirstName != "string" {
		arg = append(arg, request.FirstName)
		conditions = append(conditions, fmt.Sprintf("first_name = $%d", len(arg)))
	}

	if request.LastName != "" && request.LastName != "string" {
		arg = append(arg, request.LastName)
		conditions = append(conditions, fmt.Sprintf("last_name = $%d", len(arg)))
	}

	if request.Username != "" && request.Username != "string" {
		arg = append(arg, request.Username)
		conditions = append(conditions, fmt.Sprintf("username = $%d", len(arg)))
	}

	if request.PhoneNumber != "" && request.PhoneNumber != "string" {
		arg = append(arg, request.PhoneNumber)
		conditions = append(conditions, fmt.Sprintf("phone_number = $%d", len(arg)))
	}

	if request.Email != "" && request.Email != "string" {
		arg = append(arg, request.Email)
		conditions = append(conditions, fmt.Sprintf("email = $%d", len(arg)))
	}

	if request.Gender != "" && request.Gender != "string" {
		arg = append(arg, request.Gender)
		conditions = append(conditions, fmt.Sprintf("gender = $%d", len(arg)))
	}

	if request.Role != "" && request.Role != "string" {
		arg = append(arg, request.Role)
		conditions = append(conditions, fmt.Sprintf("role = $%d", len(arg)))
	}

	if len(conditions) > 0 {
		query += ", " + strings.Join(conditions, ", ")
	}

	query += fmt.Sprintf(" WHERE id = $%d", len(arg)+1)
	arg = append(arg, request.Id)

	_, err := u.DB.ExecContext(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *Repository) ChangeUserPassword(ctx context.Context, request *pb.UserRecoverPasswordRequest) (*pb.Void, error) {
	query := `
		UPDATE users SET password = ?, updated_at = time.Now()
		WHERE email = ? AND deleted_at is null
	`

	_, err := u.DB.ExecContext(ctx, query, request.NewPassword, request.Email)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *Repository) Delete(ctx context.Context, request *pb.ById) (*pb.Void, error) {
	query := `
		UPDATE users SET deleted_at = EXTRACT(EPOCH FROM NOW) WHERE id = ?
	`
	_, err := u.DB.ExecContext(ctx, query, request.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *Repository) GetList(ctx context.Context, filter *pb.FilterUser) (*pb.UserGetAll, error) {

	where := fmt.Sprintf(`where deleted_at isnull`)
	if filter.Role != "" {
		where += fmt.Sprintf(" and role = '%s'", filter.Role)
	}
	if filter.Username != "" {
		where += fmt.Sprintf(" and username = '%s'", filter.Username)
	}
	if filter.Firstname != "" {
		where += fmt.Sprintf(" and first_name = '%s'", filter.Firstname)
	}
	if filter.Lastname != "" {
		where += fmt.Sprintf(" and last_name = '%s'", filter.Lastname)
	}
	if filter.Gender != "" {
		where += fmt.Sprintf(" and gender = '%s'", filter.Gender)
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

	query := fmt.Sprintf(`
				select 
				    COUNT(id) OVER () AS total_count,
				    id, 
				    first_name,
				    last_name,
				    username,
				    email,
				    phone,
				    gender,
				    role 
				from users %s 
				order by created_at desc %s %s`, where, limit, offset)

	rows, err := u.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.New("querying users")
	}
	defer rows.Close()

	var response []*pb.UserGetList
	var count int32
	for rows.Next() {
		var user pb.UserGetList
		err := rows.Scan(
			&count,
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Email,
			&user.PhoneNumber,
			&user.Email,
			&user.Gender,
			&user.Role,
		)
		if err != nil {
			return nil, err
		}
		response = append(response, &user)
	}

	return &pb.UserGetAll{
		User:  response,
		Count: count,
	}, nil
}

func (u *Repository) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	query := `
		SELECT id, email
		FROM users WHERE email = ? AND password = ? AND deleted_at is null
	`

	var user pb.LoginResponse

	err := u.DB.QueryRowContext(ctx, query, request.Email, request.Password).Scan(
		&user.Id, &user.Username,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
