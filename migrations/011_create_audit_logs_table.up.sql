-- Tabela para registrar ações importantes e sensíveis no sistema
CREATE TABLE `audit_logs` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `user_id` BIGINT NOT NULL COMMENT 'ID do usuário (geralmente admin) que realizou a ação',
  `action` VARCHAR(100) NOT NULL COMMENT 'Ação realizada, ex: enable_module, update_fiscal_settings',
  `target_type` VARCHAR(50) COMMENT 'O tipo de entidade que foi afetada, ex: system_modules, fiscal_settings',
  `target_id` BIGINT COMMENT 'O ID da entidade que foi afetada',
  `details` TEXT COMMENT 'Um JSON ou texto com detalhes da mudança, ex: valores antigos e novos',
  `ip_address` VARCHAR(45) COMMENT 'Endereço IP de onde a ação foi originada',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) COMMENT='Registros de auditoria para ações críticas do sistema.';
