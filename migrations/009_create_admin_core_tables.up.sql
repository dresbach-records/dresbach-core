-- Função de trigger para atualizar o campo updated_at (idempotente)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

-- Create ENUM types for PostgreSQL if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'plan_billing_cycle') THEN
        CREATE TYPE plan_billing_cycle AS ENUM('monthly', 'annually', 'quarterly', 'semiannually', 'biennially');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'plan_status') THEN
        CREATE TYPE plan_status AS ENUM('active', 'hidden', 'archived');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'api_credential_status') THEN
        CREATE TYPE api_credential_status AS ENUM('active', 'inactive');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'audit_log_result') THEN
        CREATE TYPE audit_log_result AS ENUM('success', 'failure');
    END IF;
END$$;

-- Módulo 1 e 3: Configurações do Site e Modo Manutenção
CREATE TABLE IF NOT EXISTS site_settings (
    id INT PRIMARY KEY DEFAULT 1,
    company_name VARCHAR(255) NULL,
    slogan VARCHAR(255) NULL,
    description TEXT NULL,
    phone_numbers JSONB NULL,
    whatsapp VARCHAR(50) NULL,
    institutional_email VARCHAR(255) NULL,
    address TEXT NULL,
    social_links JSONB NULL,
    logo_url VARCHAR(255) NULL,
    favicon_url VARCHAR(255) NULL,
    maintenance_enabled BOOLEAN DEFAULT FALSE,
    maintenance_message TEXT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT single_row_site CHECK (id = 1)
);
INSERT INTO site_settings (id) VALUES (1) ON CONFLICT (id) DO NOTHING;

-- Trigger for site_settings
DROP TRIGGER IF EXISTS update_site_settings_updated_at ON site_settings;
CREATE TRIGGER update_site_settings_updated_at
BEFORE UPDATE ON site_settings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Módulo 2: Gerenciamento de Planos
CREATE TABLE IF NOT EXISTS plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    category VARCHAR(100) NULL,
    price DECIMAL(10, 2) NOT NULL,
    billing_cycle plan_billing_cycle NOT NULL,
    features JSONB NULL,
    whm_package_name VARCHAR(255) NULL,
    status plan_status DEFAULT 'active',
    display_order INT DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Trigger for plans
DROP TRIGGER IF EXISTS update_plans_updated_at ON plans;
CREATE TRIGGER update_plans_updated_at
BEFORE UPDATE ON plans
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Módulo 4: Gerenciamento de APIs
CREATE TABLE IF NOT EXISTS api_credentials (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(50) NOT NULL UNIQUE,
    encrypted_key TEXT NOT NULL,
    status api_credential_status DEFAULT 'inactive',
    last_test_at TIMESTAMPTZ NULL,
    last_error TEXT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Trigger for api_credentials
DROP TRIGGER IF EXISTS update_api_credentials_updated_at ON api_credentials;
CREATE TRIGGER update_api_credentials_updated_at
BEFORE UPDATE ON api_credentials
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- Módulo 5: Configurações de Sistema
CREATE TABLE IF NOT EXISTS system_settings (
    id INT PRIMARY KEY DEFAULT 1,
    allow_new_registrations BOOLEAN DEFAULT TRUE,
    manual_order_approval BOOLEAN DEFAULT FALSE,
    auto_provisioning BOOLEAN DEFAULT TRUE,
    auto_emails BOOLEAN DEFAULT TRUE,
    default_language VARCHAR(10) DEFAULT 'pt-BR',
    default_currency VARCHAR(5) DEFAULT 'BRL',
    default_timezone VARCHAR(50) DEFAULT 'America/Sao_Paulo',
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT single_row_system CHECK (id = 1)
);
INSERT INTO system_settings (id) VALUES (1) ON CONFLICT (id) DO NOTHING;


-- Trigger for system_settings
DROP TRIGGER IF EXISTS update_system_settings_updated_at ON system_settings;
CREATE TRIGGER update_system_settings_updated_at
BEFORE UPDATE ON system_settings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- Módulo 7: Acesso Admin (RBAC)
CREATE TABLE IF NOT EXISTS staff (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Trigger for staff
DROP TRIGGER IF EXISTS update_staff_updated_at ON staff;
CREATE TRIGGER update_staff_updated_at
BEFORE UPDATE ON staff
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE IF NOT EXISTS roles ( id SERIAL PRIMARY KEY, name VARCHAR(50) NOT NULL UNIQUE, description TEXT NULL );
CREATE TABLE IF NOT EXISTS permissions ( id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL UNIQUE, description TEXT NULL );
CREATE TABLE IF NOT EXISTS staff_roles ( staff_id INT NOT NULL, role_id INT NOT NULL, PRIMARY KEY (staff_id, role_id), FOREIGN KEY (staff_id) REFERENCES staff(id) ON DELETE CASCADE, FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE );
CREATE TABLE IF NOT EXISTS role_permissions ( role_id INT NOT NULL, permission_id INT NOT NULL, PRIMARY KEY (role_id, permission_id), FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE, FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE );

-- Módulo 6: Logs de Auditoria
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    staff_id INT NULL,
    ip_address VARCHAR(45) NULL,
    action VARCHAR(255) NOT NULL,
    target_type VARCHAR(50) NULL,
    target_id VARCHAR(255) NULL,
    old_value JSONB NULL,
    new_value JSONB NULL,
    result audit_log_result NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (staff_id) REFERENCES staff(id) ON DELETE SET NULL
);