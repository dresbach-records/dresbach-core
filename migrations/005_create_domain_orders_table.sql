CREATE TABLE IF NOT EXISTS domain_orders (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    domain_name VARCHAR(255) NOT NULL UNIQUE,
    document VARCHAR(20) NOT NULL, -- CPF ou CNPJ para o registro
    status VARCHAR(50) NOT NULL DEFAULT 'pending_payment', -- pending_payment, pending_registration, completed, failed
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (client_id) REFERENCES clients(id)
);

-- Trigger para a tabela domain_orders
CREATE TRIGGER update_domain_orders_updated_at
BEFORE UPDATE ON domain_orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
