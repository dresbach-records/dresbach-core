-- No PostgreSQL, INT AUTO_INCREMENT é substituído por SERIAL.
-- A cláusula ON UPDATE CURRENT_TIMESTAMP não existe e é substituída por um trigger.

CREATE TABLE IF NOT EXISTS vps_orders (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    invoice_id INT NOT NULL,
    vps_instance_id VARCHAR(255), -- ID retornado pela Hostinger após o provisionamento
    plan_id VARCHAR(255) NOT NULL,
    location VARCHAR(50) NOT NULL,
    template VARCHAR(100) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL, -- Armazenar apenas o hash da senha
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, active, failed, suspended
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (client_id) REFERENCES clients(id),
    FOREIGN KEY (invoice_id) REFERENCES invoices(id)
);

-- Função de trigger para atualizar o campo updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger para a tabela vps_orders
DROP TRIGGER IF EXISTS update_vps_orders_updated_at ON vps_orders;
CREATE TRIGGER update_vps_orders_updated_at
BEFORE UPDATE ON vps_orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
