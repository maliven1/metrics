package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/maliven1/metrics/internal/config"
)

type PostgreDB struct {
	DB *sql.DB
}

func NewPostgreDB(cfg config.ServerConfig) (*PostgreDB, error) {
	dbConnectString := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=%s",
		"localhost", cfg.PostgreDNS, "postgres", "metrics", "12345678", "disable")

	db, err := sql.Open("pgx", dbConnectString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
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
		return nil
	}
	return nil
}
