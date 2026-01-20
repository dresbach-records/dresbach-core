CREATE TABLE IF NOT EXISTS password_resets (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);
