-- migrations/000001_create_metrics_table.up.sql
-- Создание таблицы метрик
CREATE TABLE metrics (
    id SERIAL PRIMARY KEY,
    gauge varchar(255) UNIQUE,
    gauge_value double precision ,
    count varchar(255) UNIQUE,
    count_value BIGINT  
); 

-- Базовый индекс для поиска gauge_name 
CREATE INDEX idx_gauge ON metrics(gauge);

-- Индекс для поиска count_value
CREATE INDEX idx_count ON metrics(count);