package url

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type PostgresRepository struct {
	db DBTX
}

func NewRepository(db DBTX) *PostgresRepository {
	return &PostgresRepository{db: db}
}

const columns = `id, code, original_url, click_count, created_at, expires_at`

func scan(row pgx.Row) (*URL, error) {
	var u URL
	err := row.Scan(&u.ID, &u.Code, &u.OriginalURL, &u.ClickCount, &u.CreatedAt, &u.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *PostgresRepository) FindByCode(ctx context.Context, code string) (*URL, error) {
	query := `SELECT ` + columns + ` FROM urls WHERE code = $1`
	return scan(r.db.QueryRow(ctx, query, code))
}

func (r *PostgresRepository) Create(ctx context.Context, form CreateParam) (*URL, error) {
	query := `INSERT INTO urls (code, original_url, expires_at) VALUES ($1, $2, $3) RETURNING ` + columns

	u, err := scan(r.db.QueryRow(ctx, query, form.Code, form.OriginalURL, form.ExpiresAt))
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, ErrCodeTaken
	}
	return u, err
}

func (r *PostgresRepository) GetNextID(ctx context.Context) (uint64, error) {
	var id uint64
	query := `SELECT nextval(pg_get_serial_sequence('urls', 'id'))`
	err := r.db.QueryRow(ctx, query).Scan(&id)
	return id, err
}

func (r *PostgresRepository) Resolve(ctx context.Context, code string) (string, error) {
	var target string
	query := `SELECT original_url FROM urls WHERE code = $1 AND (expires_at IS NULL OR expires_at > now())`
	err := r.db.QueryRow(ctx, query, code).Scan(&target)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	return target, err
}
