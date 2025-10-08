package storage

import (
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/maliven1/metrics/internal/config"
)

type PostgreDB struct {
	DB *sql.DB
}

func NewPostgreDB(cfg config.ServerConfig) (*PostgreDB, error) {
	if cfg.PostgreDNS == "" {
		return nil, nil
	}
	db, err := sql.Open("pgx", cfg.PostgreDNS)
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
		db.Close()
		return err
	}
	return nil
}
