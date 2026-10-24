package url

import (
	"context"
	"regexp"
	"strings"
	"time"
)

const (
	maxExpirationInHours = 24 * 365
	offset               = 56_800_235_584
)

var (
	customPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)
	reserved      = map[string]bool{"api": true, "healthz": true}
)

type Repository interface {
	FindByCode(ctx context.Context, code string) (*URL, error)
	Create(ctx context.Context, form CreateParam) (*URL, error)
	GetNextID(ctx context.Context) (uint64, error)
	Resolve(ctx context.Context, code string) (string, error)
}

type CodeEncoder interface {
	Encode(id uint64) (string, error)
}

type Service struct {
	repo    Repository
	baseURL string
}

func NewService(repo Repository, baseURL string) *Service {
	return &Service{
		repo:    repo,
		baseURL: strings.Trim(baseURL, "/"),
	}
}

func (s *Service) GetByCode(ctx context.Context, code string) (*URL, error) {
	return s.repo.FindByCode(ctx, code)
}

func (s *Service) Create(ctx context.Context, form CreateForm) (*URL, error) {
	// define expiration
	var expiresAt *time.Time
	if form.ExpiresAt != nil {
		h := *form.ExpiresAt
		if h < 1 || h > maxExpirationInHours {
			return nil, ErrInvalidInput
		}
		t := time.Now().Add(time.Duration(h) * time.Hour)
		expiresAt = &t
	}

	params := CreateParam{OriginalURL: form.OriginalURL, ExpiresAt: expiresAt}

	// if user use custsom alias
	if alias := strings.TrimSpace(form.Alias); alias != "" {
		if !customPattern.MatchString(alias) {
			return nil, ErrInvalidAlias
		}
		if reserved[alias] {
			return nil, ErrAliasTaken
		}
		params.Code = alias
		return s.insert(ctx, params)
	}

	// if user not provide custom alias, then generate code using encode the id using base62
	nextID, err := s.repo.GetNextID(ctx)
	if err != nil {
		return nil, err
	}

	genCode, err := Encode(nextID, offset)
	if err != nil {
		return nil, err
	}
	params.Code = genCode

	u, err := s.insert(ctx, params)
	return u, err
}

func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	return s.repo.Resolve(ctx, code)
}

func (s *Service) insert(ctx context.Context, param CreateParam) (*URL, error) {
	u, err := s.repo.Create(ctx, param)
	u.ShortURL = s.shortURL(u.Code)
	return u, err
}

func (s *Service) shortURL(code string) string {
	return s.baseURL + "/" + code
}
