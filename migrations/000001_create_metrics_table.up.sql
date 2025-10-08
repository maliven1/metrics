-- migrations/000001_create_metrics_table.up.sql
-- Создание таблицы метрик
CREATE TABLE metrics (
    id SERIAL PRIMARY KEY,
    gauge_name varchar(255),
    gauge_value double precision ,
    count_name varchar(255),
    count_value INTEGER 
); 

-- Базовый индекс для поиска gauge_name 
CREATE INDEX idx_gauge_name ON metrics(gauge_name);

-- Индекс для поиска count_value
CREATE INDEX idx_count_value ON metrics(count_name);