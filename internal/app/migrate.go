//go:build migrate

package app

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {
	databaseURL, ok := os.LookupEnv("MYSQL_URL")
	if !ok || len(databaseURL) == 0 {
		log.Fatalf("migrate: environment variable not declared: MYSQL_URL")
	}

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			defer m.Close()
			break
		}

		log.Printf("Migrate: mysql is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: mysql connect error: %s", err)
	}

	err = m.Up()
	m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}
