package storage

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/maliven1/metrics/internal/config"
	"go.uber.org/zap"
)

type PostgreDB struct {
	DB *sql.DB
}

func NewPostgreDB(cfg config.ServerConfig, log *zap.SugaredLogger) (*PostgreDB, error) {
	if cfg.PostgreDNS == "" {
		return nil, fmt.Errorf("PostgreSQL DNS is empty")
	}
	db, err := sql.Open("pgx", cfg.PostgreDNS)
	if err != nil {
		return nil, err
	}
	err = UpMigrations(db, log)
	if err != nil {

		return nil, err
	}
	return &PostgreDB{DB: db}, nil
}

func (db *PostgreDB) Close() error {
	return db.DB.Close()
}

func (db *PostgreDB) CheckConnection() error {
	err := db.DB.Ping()
	if err != nil {
		return err
	}
	return nil
}

func UpMigrations(db *sql.DB, log *zap.SugaredLogger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Error("Failed to create migration driver", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Path to migrations directory
		"postgres", driver)
	if err != nil {
		log.Error("Failed to create migration instance", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error("Failed to apply migrations", err)
		return err
	}

	if err == migrate.ErrNoChange {
		log.Info("No migrations to apply")
	} else {
		log.Info("Migrations applied successfully")
	}
	return nil
}
