-- Migration para remover a coluna 'last_login_ip' da tabela de usuários, que não é mais necessária.
-- ATENÇÃO: Esta é uma operação destrutiva e causará perda de dados.

ALTER TABLE users DROP COLUMN last_login_ip;

-- A migration abaixo também seria detectada como arriscada
-- ALTER TABLE invoices ALTER COLUMN due_date TYPE VARCHAR(255);
