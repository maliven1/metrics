-- migrations/000001_create_metrics_table.down.sql
-- Откат создания таблицы фильмов
DROP INDEX IF EXISTS idx_gauge_name;
DROP INDEX IF EXISTS idx_count_value;
DROP TABLE IF EXISTS metrics;
