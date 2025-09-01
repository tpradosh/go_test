CREATE TABLE IF NOT EXISTS results (
    id SERIAL PRIMARY KEY,
    watch_id INT REFERENCES watches(id),
    status_code INT,
    latency_ms INT,
    success BOOLEAN,
    checked_at TIMESTAMP
)