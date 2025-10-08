-- migrations/000001_create_metrics_table.down.sql
-- Откат создания таблицы 
DROP INDEX IF EXISTS idx_gauge;
DROP INDEX IF EXISTS idx_count;
