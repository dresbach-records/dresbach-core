-- Tabela para armazenar informações sobre os documentos fiscais emitidos
CREATE TABLE `fiscal_documents` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `invoice_id` BIGINT NOT NULL COMMENT 'ID da fatura interna relacionada a este documento fiscal',
  `provider` VARCHAR(50) DEFAULT NULL COMMENT 'Provedor que emitiu a nota',
  `nf_number` VARCHAR(50) DEFAULT NULL COMMENT 'Número da Nota Fiscal de Serviço Eletrônica',
  `verification_code` VARCHAR(50) DEFAULT NULL COMMENT 'Código de verificação da autenticidade da NFS-e',
  `status` ENUM('pending', 'processing', 'authorized', 'denied', 'error') NOT NULL DEFAULT 'pending' COMMENT 'Status da emissão da nota fiscal',
  `pdf_url` TEXT DEFAULT NULL COMMENT 'URL para o PDF da nota fiscal',
  `xml_url` TEXT DEFAULT NULL COMMENT 'URL para o XML da nota fiscal',
  `error_message` TEXT DEFAULT NULL COMMENT 'Mensagem de erro em caso de falha na emissão',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `idx_invoice_id` (`invoice_id`)
) COMMENT='Armazena os documentos fiscais (NFS-e) gerados.';
