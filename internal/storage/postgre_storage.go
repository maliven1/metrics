package storage

import (
	"context"
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
	if cfg.PostgreDSN == "" {
		return nil, fmt.Errorf("PostgreSQL DSN is empty")
	}
	db, err := sql.Open("pgx", cfg.PostgreDSN)
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

func (db *PostgreDB) SetGauge(key string, value float64, ctx context.Context) error {
	op := "path: storage/SetGauge."
	tx, err := db.DB.Begin()
	if err != nil {

		return fmt.Errorf("%s tarts a transaction err:%w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM metrics WHERE gauge = $1)", key).Scan(&exists)
	if err != nil {

		return fmt.Errorf("%s QueryRowContext err:%w", op, err)
	}

	if exists {
		_, err = tx.ExecContext(ctx, "UPDATE metrics SET gauge_value = $1 WHERE gauge = $2", value, key)
	} else {
		_, err = tx.ExecContext(ctx, "INSERT INTO metrics (gauge, gauge_value, count, count_value) VALUES ($1, $2, '', 0)", key, value)
	}

	if err != nil {
		return fmt.Errorf("%s ExecContext err:%w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s Commit err:%w", op, err)
	}
	return nil
}

func (db *PostgreDB) SetCounter(key string, value int64, ctx context.Context) error {
	op := "path: storage/SetCounter."
	tx, err := db.DB.Begin()
	if err != nil {
		return fmt.Errorf("%s tarts a transaction err:%w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM metrics WHERE count = $1)", key).Scan(&exists)
	if err != nil {

		return fmt.Errorf("%s QueryRowContext err:%w", op, err)
	}

	if exists {
		_, err = tx.ExecContext(ctx, "UPDATE metrics SET count_value = count_value + $1 WHERE count = $2", value, key)
	} else {
		_, err = tx.ExecContext(ctx, "INSERT INTO metrics (count, count_value) VALUES ($1, $2)", key, value)
	}

	if err != nil {
		return fmt.Errorf("%s ExecContext err:%w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s Commit err:%w", op, err)
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

func (db *PostgreDB) GetAllGauges() (map[string]float64, error) {
	gauges := make(map[string]float64)

	rows, err := db.DB.Query("SELECT gauge, gauge_value FROM metrics WHERE gauge IS NOT NULL AND gauge != ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var value float64
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		gauges[key] = value
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gauges, nil
}

func (db *PostgreDB) GetAllCounters() (map[string]int64, error) {
	counters := make(map[string]int64)

	rows, err := db.DB.Query("SELECT count, count_value FROM metrics WHERE count IS NOT NULL AND count != ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var value int64
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		counters[key] = value
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return counters, nil
}

func (db *PostgreDB) GetItemGauge(key string) (string, float64, error) {
	var value float64
	err := db.DB.QueryRow("SELECT gauge_value FROM metrics WHERE gauge = $1", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, fmt.Errorf("gauge with key %s not found", key)
		}
		return "", 0, err
	}
	return key, value, nil
}

func (db *PostgreDB) GetItemCounter(key string) (string, int64, error) {
	var value int64
	err := db.DB.QueryRow("SELECT count_value FROM metrics WHERE count = $1", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, fmt.Errorf("counter with key %s not found", key)
		}
		return "", 0, err
	}
	return key, value, nil
}

func (db *PostgreDB) CheckItemGauge(key string) (bool, error) {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM metrics WHERE gauge = $1)", key).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *PostgreDB) CheckCounter(key string) (bool, error) {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM metrics WHERE count = $1)", key).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
