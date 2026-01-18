-- Tabela para gerenciar módulos e suas ativações (feature flags)
CREATE TABLE `system_modules` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(50) UNIQUE NOT NULL COMMENT 'Nome único do módulo, ex: fiscal_nfse',
  `enabled` BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Se o módulo está ativado ou não',
  `enabled_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Quando o módulo foi ativado',
  `enabled_by` BIGINT NULL COMMENT 'ID do admin que ativou o módulo',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT='Gerencia a ativação de módulos do sistema.';

-- Inserir o registro inicial para o módulo fiscal, desativado por padrão
INSERT INTO `system_modules` (name, enabled) VALUES ('fiscal_nfse', FALSE);
