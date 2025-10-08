-- migrations/000001_create_metrics_table.up.sql
-- Создание таблицы фильмов
CREATE TABLE metrics (
    id SERIAL PRIMARY KEY,
    gauge_name varchar(255),
    gauge_value double precision NOT NULL,
    count_name varchar(255),
    count_value INTEGER NOT NULL
); 

-- Базовый индекс для поиска по названию
CREATE INDEX idx_gauge_name ON gauge_name(name);

-- Индекс для поиска по году
CREATE INDEX idx_count_value ON count_value(year);