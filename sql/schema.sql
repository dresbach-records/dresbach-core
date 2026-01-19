-- =========================
-- TABELA DE USUÁRIOS
-- =========================
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role ENUM('admin','client') NOT NULL DEFAULT 'client',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- =========================
-- TABELA DE CLIENTES
-- =========================
CREATE TABLE clients (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  company_name VARCHAR(255),
  contact_name VARCHAR(255),
  email VARCHAR(255) NOT NULL UNIQUE,
  phone VARCHAR(50),
  address VARCHAR(255),
  city VARCHAR(100),
  state VARCHAR(100),
  zip VARCHAR(50),
  country VARCHAR(100),
  cpf_cnpj VARCHAR(20),
  asaas_customer_id VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL,
  CONSTRAINT fk_clients_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- =========================
-- TABELA DE SERVIÇOS
-- =========================
CREATE TABLE services (
  id INT AUTO_INCREMENT PRIMARY KEY,
  client_id INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  billing_cycle ENUM('monthly','annually') NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  status ENUM('active','suspended','canceled','pending') DEFAULT 'pending',
  next_due_date DATE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_services_client FOREIGN KEY (client_id) REFERENCES clients(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- =========================
-- TABELA DE FATURAS
-- =========================
CREATE TABLE invoices (
  id INT AUTO_INCREMENT PRIMARY KEY,
  client_id INT NOT NULL,
  issue_date DATE NOT NULL,
  due_date DATE NOT NULL,
  total_amount DECIMAL(10,2) NOT NULL,
  status ENUM('paid','unpaid','overdue','canceled') DEFAULT 'unpaid',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_invoices_client FOREIGN KEY (client_id) REFERENCES clients(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- =========================
-- ITENS DA FATURA
-- =========================
CREATE TABLE invoice_items (
  id INT AUTO_INCREMENT PRIMARY KEY,
  invoice_id INT NOT NULL,
  service_id INT NOT NULL,
  description VARCHAR(255),
  amount DECIMAL(10,2) NOT NULL,
  CONSTRAINT fk_items_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id),
  CONSTRAINT fk_items_service FOREIGN KEY (service_id) REFERENCES services(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- =========================
-- TRANSAÇÕES
-- =========================
CREATE TABLE transactions (
  id INT AUTO_INCREMENT PRIMARY KEY,
  invoice_id INT,
  client_id INT NOT NULL,
  gateway VARCHAR(50),
  transaction_id_gateway VARCHAR(255),
  amount DECIMAL(10,2) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_transactions_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id),
  CONSTRAINT fk_transactions_client FOREIGN KEY (client_id) REFERENCES clients(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- =========================
-- LOGS DE AUDITORIA
-- =========================
CREATE TABLE audit_logs (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT,
  action VARCHAR(255) NOT NULL,
  target_type VARCHAR(100),
  target_id INT,
  ip_address VARCHAR(45),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_audit_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
