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

func (db *PostgreDB) SetGauge(key string, value float64) {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM metrics WHERE gauge = $1)", key).Scan(&exists)
	if err != nil {

		return
	}

	if exists {
		_, err = db.DB.Exec("UPDATE metrics SET gauge_value = $1 WHERE gauge = $2", value, key)
	} else {
		_, err = db.DB.Exec("INSERT INTO metrics (gauge, gauge_value, count, count_value) VALUES ($1, $2, '', 0)", key, value)
	}

	if err != nil {

		return
	}
}

func (db *PostgreDB) SetCounter(key string, value int64) {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM metrics WHERE count = $1)", key).Scan(&exists)
	if err != nil {

		return
	}

	if exists {
		_, err = db.DB.Exec("UPDATE metrics SET count_value = $1 WHERE count = $2", value, key)
	} else {
		_, err = db.DB.Exec("INSERT INTO metrics (count, count_value, gauge, gauge_value) VALUES ($1, $2, '', 0.0)", key, value)
	}

	if err != nil {

		return
	}
}

func (db *PostgreDB) GetGauge() map[string]float64 {
	// First get all gauge keys
	rows, err := db.DB.Query("SELECT gauge FROM metrics WHERE gauge IS NOT NULL AND gauge != ''")
	if err != nil {
		return make(map[string]float64)
	}
	defer rows.Close()

	// Collect all keys
	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			continue
		}
		keys = append(keys, key)
	}

	// For each key, use QueryRow to get the value
	gauge := make(map[string]float64)
	for _, key := range keys {
		var value float64
		err := db.DB.QueryRow("SELECT gauge_value FROM metrics WHERE gauge = $1", key).Scan(&value)
		if err != nil {
			continue
		}
		gauge[key] = value
	}

	return gauge
}

func (db *PostgreDB) GetCounter() map[string]int64 {
	// First get all counter keys
	rows, err := db.DB.Query("SELECT count FROM metrics WHERE count IS NOT NULL AND count != ''")
	if err != nil {
		return make(map[string]int64)
	}
	defer rows.Close()

	// Collect all keys
	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			continue
		}
		keys = append(keys, key)
	}

	// For each key, use QueryRow to get the value
	counter := make(map[string]int64)
	for _, key := range keys {
		var value int64
		err := db.DB.QueryRow("SELECT count_value FROM metrics WHERE count = $1", key).Scan(&value)
		if err != nil {
			continue
		}
		counter[key] = value
	}

	return counter
}

func (db *PostgreDB) AddCounter(key string, value int64) bool {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM metrics WHERE count = $1)", key).Scan(&exists)
	if err != nil {

		return false
	}

	if exists {
		_, err = db.DB.Exec("UPDATE metrics SET count_value = count_value + $1 WHERE count = $2", value, key)
	} else {
		_, err = db.DB.Exec("INSERT INTO metrics (count, count_value, gauge, gauge_value) VALUES ($1, $2, '', 0.0)", key, value)
	}

	if err != nil {

		return false
	}
	return true
}

func (db *PostgreDB) GetItemGauge(key string) (string, float64) {
	var value float64
	err := db.DB.QueryRow("SELECT gauge_value FROM metrics WHERE gauge = $1", key).Scan(&value)
	if err != nil {

		return "", 0
	}
	return key, value
}

func (db *PostgreDB) GetItemCounter(key string) (string, int64) {
	var value int64
	err := db.DB.QueryRow("SELECT count_value FROM metrics WHERE count = $1", key).Scan(&value)
	if err != nil {

		return "", 0
	}
	return key, value
}

func (db *PostgreDB) CheckCounter(key string) bool {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM metrics WHERE count = $1)", key).Scan(&exists)
	if err != nil {

		return false
	}
	return exists
}

func (db *PostgreDB) CheckItemGauge(key string) bool {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM metrics WHERE gauge = $1)", key).Scan(&exists)
	if err != nil {

		return false
	}
	return exists
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
