package repo

import (
	"database/sql"
	"posts/internal/repository"
	log "posts/internal/repository/postgres/logs"
	post "posts/internal/repository/postgres/posts"
	user "posts/internal/repository/postgres/users"
)

type Storage struct {
	UserS repository.UserI
	LogS  repository.LogI
	PostS repository.PostI
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		UserS: user.NewRepository(db),
		LogS:  log.NewRepository(db),
		PostS: post.NewRepository(db),
	}
}

func (s *Storage) User() repository.UserI {
	return s.UserS
}

func (s *Storage) Log() repository.LogI {
	return s.LogS
}

func (s *Storage) Post() repository.PostI {
	return s.PostS
}
