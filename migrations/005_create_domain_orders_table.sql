CREATE TABLE IF NOT EXISTS `domain_orders` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `client_id` INT NOT NULL,
    `domain_name` VARCHAR(255) NOT NULL UNIQUE,
    `document` VARCHAR(20) NOT NULL, -- CPF ou CNPJ para o registro
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending_payment', -- pending_payment, pending_registration, completed, failed
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
