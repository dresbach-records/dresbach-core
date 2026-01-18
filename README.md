# Backend em Go com Autenticação e MySQL

## Descrição

Este é o backend de uma aplicação, desenvolvido em Go. Ele fornece uma API com funcionalidades de autenticação e se conecta a um banco de dados MySQL. O projeto é containerizado usando Docker e Docker Compose para facilitar a configuração e o deploy.

## Funcionalidades

*   **Health Check**: Endpoint `/health` para verificar o status da API.
*   **Autenticação**: Endpoint `/auth/login` para autenticação de usuários.
*   **Banco de Dados**: Integração com banco de dados MySQL para persistência de dados.
*   **Containerização**: Uso de Docker e Docker Compose para um ambiente de desenvolvimento e produção consistente.
*   **Integração com Gemini**: Módulo para interação com a API do Gemini.

## Tecnologias Utilizadas

*   Go (versão 1.21)
*   Docker
*   Docker Compose
*   MySQL (imagem 8.0)
*   (Potentially) Google Gemini API

## Como Começar

### Pré-requisitos

*   Docker instalado
*   Docker Compose instalado

### Instalação e Execução

1.  Clone o repositório.
2.  Configure as variáveis de ambiente no arquivo `docker-compose.yml`, especialmente as senhas do banco de dados (`DB_PASSWORD`, `MYSQL_PASSWORD`, `MYSQL_ROOT_PASSWORD`).
3.  Execute o comando a seguir na raiz do projeto para construir e iniciar os contêineres:
    ```bash
    docker-compose up --build
    ```
4.  A API estará disponível em `http://localhost:8080`.

## Endpoints da API

### `GET /health`

Verifica a saúde da aplicação. Retorna `OK` com status `200` se a API estiver funcionando.

### `POST /auth/login`

Endpoint para autenticação de usuários. O corpo da requisição deve conter as credenciais do usuário (e.g., JSON com email e senha). A lógica de manipulação está em `internal/handlers/auth/auth_handler.go`.

## Configuração

As configurações da aplicação, como as credenciais do banco de dados e a porta, são gerenciadas através de variáveis de ambiente no arquivo `docker-compose.yml`.

| Variável        | Descrição                                  | Serviço |
| --------------- | ------------------------------------------ | ------- |
| `APP_PORT`      | Porta em que a aplicação Go será executada.    | `app`     |
| `DB_HOST`       | Host do banco de dados.                    | `app`     |
| `DB_PORT`       | Porta do banco de dados.                     | `app`     |
| `DB_NAME`       | Nome do banco de dados.                      | `app`     |
| `DB_USER`       | Usuário para conexão com o banco de dados. | `app`     |
| `DB_PASSWORD`   | Senha para o usuário do banco.             | `app`     |
| `MYSQL_DATABASE`| Nome do banco de dados a ser criado.       | `db`      |
| `MYSQL_USER`    | Nome do usuário do MySQL.                  | `db`      |
| `MYSQL_PASSWORD`| Senha para o usuário do MySQL.             | `db`      |
| `MYSQL_ROOT_PASSWORD`| Senha para o usuário root do MySQL.        | `db`      |

