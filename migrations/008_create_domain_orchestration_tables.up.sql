-- Criando os tipos ENUM necessários para o PostgreSQL
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'domain_type') THEN
        CREATE TYPE domain_type AS ENUM('register', 'transfer', 'existing');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'domain_status') THEN
        CREATE TYPE domain_status AS ENUM('pending_payment', 'pending_provisioning', 'active', 'failed', 'cancelled');
    END IF;
END$$;

-- Tabela para rastrear todos os domínios e seu estado atual
CREATE TABLE IF NOT EXISTS domains (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    service_id INT NOT NULL, -- Link para o serviço/produto adquirido
    domain_name VARCHAR(255) NOT NULL UNIQUE,
    type domain_type NOT NULL,
    status domain_status NOT NULL DEFAULT 'pending_payment',
    provider VARCHAR(50) NULL, -- Ex: 'hostinger', 'cloudflare'
    provider_order_id VARCHAR(255) NULL, -- ID externo do registro/transferência
    expires_at DATE NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (client_id) REFERENCES clients(id),
    FOREIGN KEY (service_id) REFERENCES services(id) -- Assumindo que você tem uma tabela `services`
);

-- Tabela para auditoria de eventos, garantindo idempotência e rastreabilidade
CREATE TABLE IF NOT EXISTS domain_events (
    id BIGSERIAL PRIMARY KEY,
    domain_id INT NOT NULL,
    type VARCHAR(100) NOT NULL, -- Ex: 'payment.succeeded', 'domain.registration.initiated', 'whm.creation.failed'
    message TEXT NULL, -- Detalhes do evento, como mensagens de erro da API
    raw_data JSONB NULL, -- Usar JSONB é geralmente melhor no PostgreSQL
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

-- Trigger para a tabela domains
DROP TRIGGER IF EXISTS update_domains_updated_at ON domains;
CREATE TRIGGER update_domains_updated_at
BEFORE UPDATE ON domains
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
