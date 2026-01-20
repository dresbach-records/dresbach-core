CREATE TABLE IF NOT EXISTS allowed_ips (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    description VARCHAR(255) NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (client_id, ip_address),
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);

ALTER TABLE clients
ADD COLUMN enforce_ip_whitelist BOOLEAN NOT NULL DEFAULT FALSE;
