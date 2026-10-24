package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gesangwidigdo/url-shortener/internal/platform/config"
	"github.com/gesangwidigdo/url-shortener/internal/platform/migration"
	"github.com/gesangwidigdo/url-shortener/migrations"
	"github.com/joho/godotenv"
)

const usage = `usage:
	migrate up 		apply all pending migrations
	migrate down [N] 	revert the last N migrations (default 1)
	migrate version 	show the current version
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if len(os.Args) < 2 {
		return fmt.Errorf("missing command\n%s", usage)
	}

	if os.Args[1] == "create" {
		if len(os.Args) < 3 {
			return fmt.Errorf("missing migration name\n%s", usage)
		}

		files, err := migration.Create("migrations", os.Args[2])
		if err != nil {
			return err
		}

		for _, f := range files {
			fmt.Println("created", f)
			return nil
		}
	}

	url := cfg.DB.URL

	switch os.Args[1] {
	case "up":
		if err := migration.Up(url, migrations.FS); err != nil {
			return err
		}
		fmt.Println("migrations applied")

	case "down":
		steps := 1
		if len(os.Args) > 2 {
			steps, err = strconv.Atoi(os.Args[2])
			if err != nil {
				return fmt.Errorf("invalid step count %q", os.Args[2])
			}
		}

		if err := migration.Down(url, migrations.FS, steps); err != nil {
			return err
		}
		fmt.Printf("reverted %d migration(s)\n", steps)

	case "version":
		v, dirty, err := migration.Version(url, migrations.FS)
		if err != nil {
			return err
		}

		fmt.Printf("version: %d, dirty: %t\n", v, dirty)

	default:
		return fmt.Errorf("unknown command %q\n%s", os.Args[1], usage)
	}

	return nil
}
