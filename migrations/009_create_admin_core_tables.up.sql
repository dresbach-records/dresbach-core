-- Módulo 1 e 3: Configurações do Site e Modo Manutenção
CREATE TABLE IF NOT EXISTS site_settings (
    id INT PRIMARY KEY DEFAULT 1,
    company_name VARCHAR(255) NULL,
    slogan VARCHAR(255) NULL,
    description TEXT NULL,
    phone_numbers JSON NULL,
    whatsapp VARCHAR(50) NULL,
    institutional_email VARCHAR(255) NULL,
    address TEXT NULL,
    social_links JSON NULL,
    logo_url VARCHAR(255) NULL,
    favicon_url VARCHAR(255) NULL,
    maintenance_enabled BOOLEAN DEFAULT FALSE,
    maintenance_message TEXT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT single_row_site CHECK (id = 1)
);
INSERT INTO site_settings (id) VALUES (1) ON DUPLICATE KEY UPDATE id=1;

-- Módulo 2: Gerenciamento de Planos
CREATE TABLE IF NOT EXISTS plans (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    category VARCHAR(100) NULL,
    price DECIMAL(10, 2) NOT NULL,
    billing_cycle ENUM('monthly', 'annually', 'quarterly', 'semiannually', 'biennially') NOT NULL,
    features JSON NULL,
    whm_package_name VARCHAR(255) NULL,
    status ENUM('active', 'hidden', 'archived') DEFAULT 'active',
    display_order INT DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Módulo 4: Gerenciamento de APIs
CREATE TABLE IF NOT EXISTS api_credentials (
    id INT AUTO_INCREMENT PRIMARY KEY,
    provider VARCHAR(50) NOT NULL UNIQUE,
    encrypted_key TEXT NOT NULL,
    status ENUM('active', 'inactive') DEFAULT 'inactive',
    last_test_at TIMESTAMP NULL,
    last_error TEXT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

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
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT single_row_system CHECK (id = 1)
);
INSERT INTO system_settings (id) VALUES (1) ON DUPLICATE KEY UPDATE id=1;

-- Módulo 7: Acesso Admin (RBAC)
CREATE TABLE IF NOT EXISTS staff (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS roles ( id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(50) NOT NULL UNIQUE, description TEXT NULL );
CREATE TABLE IF NOT EXISTS permissions ( id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100) NOT NULL UNIQUE, description TEXT NULL );
CREATE TABLE IF NOT EXISTS staff_roles ( staff_id INT NOT NULL, role_id INT NOT NULL, PRIMARY KEY (staff_id, role_id), FOREIGN KEY (staff_id) REFERENCES staff(id) ON DELETE CASCADE, FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE );
CREATE TABLE IF NOT EXISTS role_permissions ( role_id INT NOT NULL, permission_id INT NOT NULL, PRIMARY KEY (role_id, permission_id), FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE, FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE );

-- Módulo 6: Logs de Auditoria
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    staff_id INT NULL,
    ip_address VARCHAR(45) NULL,
    action VARCHAR(255) NOT NULL,
    target_type VARCHAR(50) NULL,
    target_id VARCHAR(255) NULL,
    old_value JSON NULL,
    new_value JSON NULL,
    result ENUM('success', 'failure') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (staff_id) REFERENCES staff(id) ON DELETE SET NULL
);
