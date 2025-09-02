CREATE TABLE IF NOT EXISTS results (
    id SERIAL PRIMARY KEY,
    watch_id INT REFERENCES watches(id),
    response_time_ms INT NOT NULL,
    status INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
