package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"posts/internal/pkg/config"
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
				firstname, 
				lastname, 
				username, 
				email, 
				phone, 
				gender, 
				password, 
				role
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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
		     firstname, 
		     lastname, 
		     username, 
		     email, 
		     phone, 
		     gender, 
		     role,
		     created_at
		FROM users 
			WHERE id = $1
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
		conditions = append(conditions, fmt.Sprintf("firstname = $%d", len(arg)))
	}

	if request.LastName != "" && request.LastName != "string" {
		arg = append(arg, request.LastName)
		conditions = append(conditions, fmt.Sprintf("lastname = $%d", len(arg)))
	}

	if request.Username != "" && request.Username != "string" {
		arg = append(arg, request.Username)
		conditions = append(conditions, fmt.Sprintf("username = $%d", len(arg)))
	}

	if request.PhoneNumber != "" && request.PhoneNumber != "string" {
		arg = append(arg, request.PhoneNumber)
		conditions = append(conditions, fmt.Sprintf("phone = $%d", len(arg)))
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
		UPDATE users 
		SET password = $1, updated_at = NOW()
		WHERE email = $2 AND deleted_at IS NULL
	`

	_, err := u.DB.ExecContext(ctx, query, request.NewPassword, request.Email)
	if err != nil {
		return nil, fmt.Errorf("error updating password: %w", err)
	}

	return &pb.Void{}, nil
}

func (u *Repository) Delete(ctx context.Context, request *pb.ById) (*pb.Void, error) {
	query := `
		UPDATE users 
		SET deleted_at = NOW() 
		WHERE id = $1
	`
	_, err := u.DB.ExecContext(ctx, query, request.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *Repository) GetList(ctx context.Context, filter *pb.FilterUser) (*pb.UserGetAll, error) {
	var args []interface{}
	var conditions []string
	argID := 1

	conditions = append(conditions, "deleted_at IS NULL")

	if filter.Role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argID))
		args = append(args, filter.Role)
		argID++
	}
	if filter.Username != "" {
		conditions = append(conditions, fmt.Sprintf("username = $%d", argID))
		args = append(args, filter.Username)
		argID++
	}
	if filter.Firstname != "" {
		conditions = append(conditions, fmt.Sprintf("firstname = $%d", argID))
		args = append(args, filter.Firstname)
		argID++
	}
	if filter.Lastname != "" {
		conditions = append(conditions, fmt.Sprintf("lastname = $%d", argID))
		args = append(args, filter.Lastname)
		argID++
	}
	if filter.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", argID))
		args = append(args, filter.Gender)
		argID++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var limit, offset int64

	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Page > 0 && filter.Limit > 0 {
		offset = (filter.Page - 1) * filter.Limit
	}

	query := fmt.Sprintf(`
							SELECT 
								id, 
								firstname,
								lastname,
								username,
								email,
								phone,
								gender,
								role 
							FROM users %s 
							ORDER BY created_at DESC 
							LIMIT $%d OFFSET $%d`, whereClause, argID, argID+1)

	args = append(args, limit, offset)

	rows, err := u.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying users: %w", err)
	}
	defer rows.Close()

	var response []*pb.UserGetList
	for rows.Next() {
		var user pb.UserGetList
		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Email,
			&user.PhoneNumber,
			&user.Gender,
			&user.Role,
		)
		if err != nil {
			return nil, err
		}
		response = append(response, &user)
	}
	countArgs := args[:len(args)-2]
	countQuery := `SELECT COUNT(*) FROM users`
	if whereClause != "" {
		countQuery += " " + whereClause
	}

	var totalCount int32
	if len(countArgs) > 0 {
		err := u.DB.QueryRow(countQuery, countArgs...).Scan(&totalCount)
		if err != nil {
			return nil, fmt.Errorf("counting users: %w", err)
		}
	} else {
		err := u.DB.QueryRow(countQuery).Scan(&totalCount)
		if err != nil {
			return nil, fmt.Errorf("counting users: %w", err)
		}
	}

	return &pb.UserGetAll{
		User:  response,
		Count: totalCount,
	}, nil
}

func (u *Repository) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {

	query := `
		SELECT id, email, password
		FROM users WHERE email = $1 AND deleted_at is null
	`

	var user pb.LoginResponse
	var hashedPassword string

	err := u.DB.QueryRowContext(ctx, query, request.Email).Scan(
		&user.Id, &user.Email, &hashedPassword,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !config.CheckPasswordHash(request.Password, hashedPassword) {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
