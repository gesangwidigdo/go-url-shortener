package migration

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Up(databaseURL string, files fs.FS) error {
	m, err := newMigrator(databaseURL, files)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

func Down(databaseURL string, files fs.FS, steps int) error {
	if steps < 1 {
		return errors.New("steps must be at least 1")
	}

	m, err := newMigrator(databaseURL, files)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate down: %w", err)
	}
	return nil
}

func Version(databaseURL string, files fs.FS) (version uint, dirty bool, err error) {
	m, err := newMigrator(databaseURL, files)
	if err != nil {
		return 0, false, err
	}
	defer closeMigrator(m)

	version, dirty, err = m.Version()
	if errors.Is(err, migrate.ErrNoChange) {
		return 0, false, nil
	}
	
	return version, dirty, err
}

func newMigrator(databaseURL string, files fs.FS) (*migrate.Migrate, error) {
	src, err := iofs.New(files, ".")
	if err != nil {
		return nil, fmt.Errorf("load migration files: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, toPgx5URL(databaseURL))
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return m, nil
}

func closeMigrator(m *migrate.Migrate) {
	_, _ = m.Close()
}

func toPgx5URL(url string) string {
	for _, prefix := range []string{"postgresql://", "postgres://"} {
		if strings.HasPrefix(url, prefix) {
			return "pgx5://" + strings.TrimPrefix(url, prefix)
		}
	}

	return url
}