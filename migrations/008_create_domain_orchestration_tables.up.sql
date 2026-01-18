
-- Tabela para rastrear todos os domínios e seu estado atual
CREATE TABLE IF NOT EXISTS domains (
    id INT AUTO_INCREMENT PRIMARY KEY,
    client_id INT NOT NULL,
    service_id INT NOT NULL, -- Link para o serviço/produto adquirido
    domain_name VARCHAR(255) NOT NULL UNIQUE,
    type ENUM('register', 'transfer', 'existing') NOT NULL,
    status ENUM('pending_payment', 'pending_provisioning', 'active', 'failed', 'cancelled') NOT NULL DEFAULT 'pending_payment',
    provider VARCHAR(50) NULL, -- Ex: 'hostinger', 'cloudflare'
    provider_order_id VARCHAR(255) NULL, -- ID externo do registro/transferência
    expires_at DATE NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES clients(id),
    FOREIGN KEY (service_id) REFERENCES services(id) -- Assumindo que você tem uma tabela `services`
);

-- Tabela para auditoria de eventos, garantindo idempotência e rastreabilidade
CREATE TABLE IF NOT EXISTS domain_events (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    domain_id INT NOT NULL,
    type VARCHAR(100) NOT NULL, -- Ex: 'payment.succeeded', 'domain.registration.initiated', 'whm.creation.failed'
    message TEXT NULL, -- Detalhes do evento, como mensagens de erro da API
    raw_data JSON NULL, -- Para armazenar payloads de webhooks ou respostas de API
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE CASCADE
);
