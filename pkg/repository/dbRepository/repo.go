package dbrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/models"
	"github.com/0loff/go_link_shortener/pkg/repository"
)

type DBRepository struct {
	DB     *sql.DB
	DSNcfg string
}

func NewRepository(DSNstring string) *DBRepository {
	conn, err := pgx.ParseConfig(DSNstring)
	if err != nil {
		panic(err)
	}

	db, err := repository.NewPostgresDB(conn.ConnString())
	if err != nil {
		panic(err)
	}

	DBRepo := &DBRepository{
		DB:     db,
		DSNcfg: DSNstring,
	}

	DBRepo.CreateTable()
	return DBRepo
}

func (dbrepo *DBRepository) CreateTable() {
	_, err := dbrepo.DB.Exec("CREATE TABLE IF NOT EXISTS shorturls (id serial PRIMARY KEY, user_id TEXT NOT NULL, short_url TEXT NOT NULL, origin_url TEXT NOT NULL, is_deleted BOOL DEFAULT false)")
	if err != nil {
		panic(err)
	}

	_, err = dbrepo.DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS origin_url ON shorturls (origin_url)")
	if err != nil {
		panic(err)
	}
}

func (dbrepo *DBRepository) FindByID(ctx context.Context, encodedURL string) (string, error) {
	var Entry models.URLEntry
	rows, err := dbrepo.DB.QueryContext(ctx, "SELECT short_url, origin_url, is_deleted FROM shorturls WHERE short_url = $1", encodedURL)
	if err != nil {
		logger.Log.Error("Request execution error", zap.Error(err))
	}

	if err = rows.Err(); err != nil {
		logger.Log.Error("Request execution error", zap.Error(err))
	}

	for rows.Next() {
		err = rows.Scan(&Entry.ShortURL, &Entry.OriginalURL, &Entry.IsDeleted)
		if err != nil {
			logger.Log.Error("Unrecognized data from the database", zap.Error(err))
			return "", err
		}
	}

	if Entry.OriginalURL == "" {
		return "", repository.ErrURLNotFound
	}

	if Entry.IsDeleted {
		return "", repository.ErrURLGone
	}

	return Entry.OriginalURL, nil
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

	rows, err := dbrepo.DB.Query("SELECT short_url, origin_url, is_deleted FROM shorturls WHERE user_id = $1", uid)
	if err != nil {
		logger.Log.Error("Unrecognized data from the database \n", zap.Error(err))
	}

	defer rows.Close()

	for rows.Next() {
		var Entry models.URLEntry
		if err := rows.Scan(&Entry.ShortURL, &Entry.OriginalURL, &Entry.IsDeleted); err != nil {
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

func (dbrepo *DBRepository) SetDelShortURLS(ShortURLsList []models.DelURLEntry) error {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, dbrepo.DSNcfg)
	if err != nil {
		logger.Log.Error("Failed to run pgxpool connection", zap.Error(err))
	}
	defer dbpool.Close()

	tx, err := dbpool.Begin(ctx)
	if err != nil {
		logger.Log.Error("Failed to start transaction", zap.Error(err))
	}

	b := &pgx.Batch{}

	updateStatement := `UPDATE shorturls SET is_deleted = true WHERE user_id = $1 AND short_url = $2;`

	for _, URLEntry := range ShortURLsList {
		b.Queue(updateStatement, URLEntry.UserID, URLEntry.ShortURL)
	}

	results := tx.SendBatch(ctx, b)
	_, err = results.Exec()
	if err != nil {
		results.Close()
		logger.Log.Error("Failed to update entries", zap.Error(err))
	}

	results.Close()
	return tx.Commit(ctx)
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
