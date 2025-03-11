package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"posts/internal/pkg/config"
)

type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage(cfg *config.Config) (*PostgresStorage, error) {
	// PostgreSQL DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDatabase,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Failed to open DB: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping DB: %v", err)
		return nil, err
	}

	return &PostgresStorage{DB: db}, nil
}

func (db *PostgresStorage) Close() error {
	return db.DB.Close()
}
