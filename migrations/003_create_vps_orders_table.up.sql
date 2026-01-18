CREATE TABLE IF NOT EXISTS vps_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    client_id INT NOT NULL,
    invoice_id INT NOT NULL,
    vps_instance_id VARCHAR(255), -- ID retornado pela Hostinger após o provisionamento
    plan_id VARCHAR(255) NOT NULL,
    location VARCHAR(50) NOT NULL,
    template VARCHAR(100) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL, -- Armazenar apenas o hash da senha
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, active, failed, suspended
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES clients(id),
    FOREIGN KEY (invoice_id) REFERENCES invoices(id)
);