-- Tabela para gerenciar módulos e suas ativações (feature flags)
CREATE TABLE IF NOT EXISTS system_modules (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT FALSE,
  enabled_at TIMESTAMPTZ NULL DEFAULT NULL,
  enabled_by BIGINT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Adicionar comentários às colunas
COMMENT ON TABLE system_modules IS 'Gerencia a ativação de módulos do sistema.';
COMMENT ON COLUMN system_modules.name IS 'Nome único do módulo, ex: fiscal_nfse';
COMMENT ON COLUMN system_modules.enabled IS 'Se o módulo está ativado ou não';
COMMENT ON COLUMN system_modules.enabled_at IS 'Quando o módulo foi ativado';
COMMENT ON COLUMN system_modules.enabled_by IS 'ID do admin que ativou o módulo';


-- Inserir o registro inicial para o módulo fiscal, desativado por padrão
INSERT INTO system_modules (name, enabled) VALUES ('fiscal_nfse', FALSE) ON CONFLICT (name) DO NOTHING;

-- Trigger para a tabela system_modules
DROP TRIGGER IF EXISTS update_system_modules_updated_at ON system_modules;
CREATE TRIGGER update_system_modules_updated_at
BEFORE UPDATE ON system_modules
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
