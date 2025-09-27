package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/thekrauss/beto-shared/pkg/errors"
)

// runs migrations from a path
func RunMigrations(sqlDB *sql.DB, driver, migrationsPath string) error {
	var dbDriver string
	switch driver {
	case "postgres":
		dbDriver = "postgres"
	case "mysql":
		dbDriver = "mysql"
	default:
		return errors.New(errors.CodeDBError, "unsupported driver for migrations")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		dbDriver,
	)
	if err != nil {
		return errors.Wrap(err, errors.CodeDBError, "failed to init migrations")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, errors.CodeDBError, "migration failed")
	}

	return nil
}

func RunMigrationsWithURL(databaseURL, migrationsPath string) error {
	m, err := migrate.New(
		"file://"+migrationsPath,
		databaseURL,
	)
	if err != nil {
		return errors.Wrap(err, errors.CodeDBError, "failed to init migrate instance")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, errors.CodeDBError, "failed to run migrations")
	}
	return nil
}
