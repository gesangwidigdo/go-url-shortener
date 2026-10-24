package url

import (
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("url not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrCodeTaken     = errors.New("code already taken")
	ErrAliasTaken    = errors.New("alias already taken")
	ErrInvalidAlias  = errors.New("alias invalid, valid: A-Z, a-z, 0-9")
	ErrInvalidURL    = errors.New("invalid url")
	ErrInvalidExpiry = errors.New("invalid expiry")
)

type URL struct {
	ID          uint64    `json:"id"`
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	ClickCount  int       `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type CreateForm struct {
	Alias       string `json:"alias"`
	OriginalURL string `json:"original_url"`
	ExpiresAt   *int   `json:"expires_at"`
}

type CreateParam struct {
	Code        string
	OriginalURL string
	ExpiresAt   *time.Time
}
