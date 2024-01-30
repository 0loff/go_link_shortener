package repository

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Метод создания соединения с БД postgres
func NewPostgresDB(cfg string) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg)

	if err != nil {
		return nil, err
	}

	return db, err
}
