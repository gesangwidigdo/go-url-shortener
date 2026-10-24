package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Env     string
	BaseURL string
	DB      DBConfig
	HTTP    HTTPConfig
}

type DBConfig struct {
	URL      string
	MaxConns int32
}

type HTTPConfig struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		DB: DBConfig{
			URL:      os.Getenv("DB_URL"),
			MaxConns: int32(getIntEnv("DB_MAX_OPEN_CONNS", 25)),
		},
		HTTP: HTTPConfig{
			Addr:            ":" + strconv.Itoa(getIntEnv("HTTP_ADDR", 8080)),
			ReadTimeout:     getDuration("HTTP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    getDuration("HTTP_WRITE_TIMEOUT", 5*time.Second),
			ShutdownTimeout: getDuration("HTTP_SHUTDOWN_TIMEOUT", 15*time.Second),
		},
		BaseURL: getEnv("BASE_URL", "http://localhost:8080"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.DB.URL == "" {
		return errors.New("DB_URL is required!")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
