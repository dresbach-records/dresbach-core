-- Tabela para armazenar as configurações fiscais da empresa
CREATE TABLE `fiscal_settings` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `provider` VARCHAR(50) DEFAULT NULL COMMENT 'Provedor de NFS-e (ex: eNotas, FocusNFe)',
  `company_name` VARCHAR(255) DEFAULT NULL COMMENT 'Razão Social da empresa',
  `cnpj` VARCHAR(20) DEFAULT NULL COMMENT 'CNPJ da empresa',
  `municipal_registration` VARCHAR(50) DEFAULT NULL COMMENT 'Inscrição Municipal',
  `city` VARCHAR(100) DEFAULT NULL COMMENT 'Cidade da empresa',
  `state` CHAR(2) DEFAULT NULL COMMENT 'Estado da empresa (UF)',
  `iss_rate` DECIMAL(5,2) DEFAULT NULL COMMENT 'Alíquota de ISS em porcentagem',
  `environment` ENUM('sandbox','production') DEFAULT 'sandbox' COMMENT 'Ambiente de emissão das notas (teste ou produção)',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT='Configurações fiscais para emissão de notas.';

-- Inserir um registro inicial vazio para que o frontend possa fazer um PUT em vez de um POST
INSERT INTO `fiscal_settings` (id) VALUES (1);
