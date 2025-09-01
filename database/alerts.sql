Create TABLE IF NOT EXISTS alerts(
    id SERIAL PRIMARY KEY,
    watch_id INT REFERENCES watches(id),
    msg TEXT,
    sent_at TIMESTAMP
)