-- Tabela de Clientes
CREATE TABLE `clients` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `user_id` INT NOT NULL,
  `company_name` VARCHAR(255),
  `contact_name` VARCHAR(255),
  `email` VARCHAR(255) NOT NULL UNIQUE,
  `phone` VARCHAR(50),
  `address` VARCHAR(255),
  `city` VARCHAR(100),
  `state` VARCHAR(100),
  `zip` VARCHAR(50),
  `country` VARCHAR(100),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
);

-- Tabela de Serviços
CREATE TABLE `services` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `client_id` INT NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `description` TEXT,
  `billing_cycle` ENUM('monthly', 'quarterly', 'semi-annually', 'annually', 'one-time') NOT NULL,
  `price` DECIMAL(10, 2) NOT NULL,
  `status` ENUM('active', 'suspended', 'canceled', 'pending') NOT NULL DEFAULT 'pending',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `next_due_date` DATE,
  FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`)
);

-- Tabela de Faturas
CREATE TABLE `invoices` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `client_id` INT NOT NULL,
  `issue_date` DATE NOT NULL,
  `due_date` DATE NOT NULL,
  `total_amount` DECIMAL(10, 2) NOT NULL,
  `status` ENUM('paid', 'unpaid', 'overdue', 'canceled') NOT NULL DEFAULT 'unpaid',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`)
);

-- Tabela de Itens da Fatura (para associar serviços a faturas)
CREATE TABLE `invoice_items` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `invoice_id` INT NOT NULL,
    `service_id` INT NOT NULL,
    `description` VARCHAR(255) NOT NULL,
    `amount` DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (`invoice_id`) REFERENCES `invoices`(`id`),
    FOREIGN KEY (`service_id`) REFERENCES `services`(`id`)
);


-- Tabela de Transações
CREATE TABLE `transactions` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `invoice_id` INT,
  `client_id` INT NOT NULL,
  `date` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `description` VARCHAR(255),
  `amount` DECIMAL(10, 2) NOT NULL,
  `gateway` VARCHAR(50),
  `transaction_id_gateway` VARCHAR(255),
  FOREIGN KEY (`invoice_id`) REFERENCES `invoices`(`id`),
  FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`)
);

-- Tabela de Logs de Auditoria
CREATE TABLE `audit_logs` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `user_id` INT,
  `action` VARCHAR(255) NOT NULL,
  `target_id` INT,
  `target_type` VARCHAR(100),
  `ip_address` VARCHAR(45),
  `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
);
