Create TABLE IF NOT EXISTS watches(
    id SERIAL PRIMARY KEY,
    link TEXT,
    interval INT,
    created_at TIMESTAMP,
    expected_status INT
)