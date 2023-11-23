package dbrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/models"
	"go_link_shortener/pkg/repository"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
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
	_, err := dbrepo.DB.Exec("CREATE TABLE IF NOT EXISTS shorturls (id serial PRIMARY KEY, user_id TEXT NOT NULL, short_url TEXT NOT NULL, origin_url TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}

	_, err = dbrepo.DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS origin_url ON shorturls (origin_url)")
	if err != nil {
		panic(err)
	}
}

func (dbrepo *DBRepository) FindByID(ctx context.Context, encodedURL string) string {

	row := dbrepo.DB.QueryRowContext(ctx, "SELECT origin_url FROM shorturls WHERE short_url = $1", encodedURL)

	var originURL string
	err := row.Scan(&originURL)
	if err != nil {
		logger.Log.Error("Unrecognized data from the database", zap.Error(err))
		return ""
	}

	return originURL
}

func (dbrepo *DBRepository) FindByLink(ctx context.Context, link string) string {
	row := dbrepo.DB.QueryRowContext(ctx, "SELECT short_url FROM shorturls WHERE origin_url = $1", link)

	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil {
		logger.Log.Error("Unrecognized data from the database \n", zap.Error(err))
		return ""
	}

	return shortURL
}

func (dbrepo *DBRepository) FindByUser(ctx context.Context, uid string) []models.URLEntry {
	var URLEntries []models.URLEntry

	rows, err := dbrepo.DB.Query("SELECT short_url, origin_url FROM shorturls WHERE user_id = $1", uid)
	if err != nil {
		logger.Log.Error("Unrecognized data from the database \n", zap.Error(err))
	}

	defer rows.Close()

	for rows.Next() {
		var Entry models.URLEntry
		if err := rows.Scan(&Entry.ShortURL, &Entry.OriginalURL); err != nil {
			logger.Log.Error("Unable to parse the received value", zap.Error(err))
			continue
		}

		URLEntries = append(URLEntries, Entry)
	}

	if err = rows.Err(); err != nil {
		return URLEntries
	}

	return URLEntries
}

func (dbrepo *DBRepository) SetShortURL(ctx context.Context, userID, shortURL, origURL string) (string, error) {
	_, err := dbrepo.DB.Exec("INSERT INTO shorturls (user_id, short_url, origin_url) VALUES ($1, $2, $3)", userID, shortURL, origURL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return "", repository.ErrConflict
		}
	}
	return shortURL, err
}

func (dbrepo *DBRepository) BatchInsertShortURLS(ctx context.Context, uid string, urls []models.URLEntry) error {
	var (
		placeholders []string
		newUrls      []interface{}
	)

	for index, url := range urls {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d)",
			index*3+1,
			index*3+2,
			index*3+3,
		))

		newUrls = append(newUrls, uid, url.ShortURL, url.OriginalURL)
	}

	tx, err := dbrepo.DB.Begin()
	if err != nil {
		logger.Log.Error("Failed to start transaction", zap.Error(err))
	}

	insertStatement := fmt.Sprintf("INSERT INTO shorturls (user_id, short_url, origin_url) VALUES %s", strings.Join(placeholders, ","))
	_, err = tx.ExecContext(ctx, insertStatement, newUrls...)
	if err != nil {
		tx.Rollback()
		logger.Log.Error("Failed to insert multiple records", zap.Error(err))
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error("Failed to commit transaction", zap.Error(err))
	}

	return nil
}

func (dbrepo *DBRepository) GetNumberOfEntries(ctx context.Context) int {
	row := dbrepo.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM shorturls")

	var Num int
	err := row.Scan(&Num)
	if err != nil {
		return 0
	}

	return Num
}

func (dbrepo *DBRepository) PingConnect(ctx context.Context) error {
	err := dbrepo.DB.Ping()
	if err != nil {
		return err
	}

	return nil
}
