package dbrepository

import (
	"context"
	"database/sql"
	"go_link_shortener/internal/models"
	"go_link_shortener/pkg/repository"
	"time"

	"github.com/jackc/pgx/v5"
)

type DBRepository struct {
	DB *sql.DB
}

func NewRepository(DSNCfg string) *DBRepository {
	conn, err := pgx.ParseConfig(DSNCfg)
	if err != nil {
		panic(err)
	}

	db, err := repository.NewPostgresDB(conn.ConnString())
	if err != nil {
		panic(err)
	}

	DBRepo := &DBRepository{
		DB: db,
	}

	DBRepo.CreateTable()
	return DBRepo
}

func (dbrepo *DBRepository) CreateTable() {
	_, err := dbrepo.DB.Exec("CREATE TABLE IF NOT EXISTS shorturls (id serial PRIMARY KEY, short_url TEXT NOT NULL, origin_url TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}
}

func (dbrepo *DBRepository) FindByID(encodedURL string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := dbrepo.DB.QueryRowContext(ctx, "SELECT origin_url FROM shorturls WHERE short_url = $1", encodedURL)

	var originURL string
	err := row.Scan(&originURL)
	if err != nil {
		return ""
	}

	return originURL
}

func (dbrepo *DBRepository) FindByLink(link string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := dbrepo.DB.QueryRowContext(ctx, "SELECT short_url FROM shorturls WHERE origin_url = $1", link)

	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil {
		return ""
	}

	return shortURL
}

func (dbrepo *DBRepository) SetShortURL(shortURL, origURL string) {
	_, err := dbrepo.DB.Exec("INSERT INTO shorturls (short_url, origin_url) VALUES ($1, $2)", shortURL, origURL)
	if err != nil {
		panic(err)
	}
}

func (dbrepo *DBRepository) BatchInsertShortURLS(urls []models.BatchInsertURLEntry) error {
	ctx := context.Background()
	tx, err := dbrepo.DB.Begin()
	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO shorturls (short_url, origin_url) VALUES ($1, $2)")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	for _, u := range urls {
		_, err := stmt.ExecContext(ctx, u.ShortURL, u.OriginalURL)
		if err != nil {
			panic(err)
		}
	}

	return tx.Commit()
}

func (dbrepo *DBRepository) GetNumberOfEntries() int {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := dbrepo.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM shorturls")

	var Num int
	err := row.Scan(&Num)
	if err != nil {
		return 0
	}

	return Num
}

func (dbrepo *DBRepository) PingConnect() error {
	err := dbrepo.DB.Ping()
	if err != nil {
		return err
	}

	return nil
}
