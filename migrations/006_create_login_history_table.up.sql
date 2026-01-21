CREATE TABLE IF NOT EXISTS login_history (
    id SERIAL PRIMARY KEY,
    client_id INT NULL,
    ip_address VARCHAR(45) NOT NULL,
    user_agent VARCHAR(255) NULL,
    was_successful BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_login_history_client_id ON login_history(client_id);
CREATE INDEX IF NOT EXISTS idx_login_history_ip_address ON login_history(ip_address);
