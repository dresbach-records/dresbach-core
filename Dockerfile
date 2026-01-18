# Estágio de Build
# Usamos uma imagem Go para compilar nossa aplicação.
FROM golang:1.21-alpine AS builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências
COPY go.mod go.sum ./

# Baixa as dependências. O comando `go mod download` resolve as dependências e as armazena em cache.
RUN go mod download

# Copia todo o código-fonte da aplicação para o diretório de trabalho.
COPY . .

# Compila a aplicação. 
# - CGO_ENABLED=0 desabilita o CGO para criar um binário estático.
# -o /bin/app especifica que o output da compilação será um arquivo chamado 'app' no diretório /bin.
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/api

# --- Estágio Final ---
# Usamos uma imagem Alpine mínima para a imagem final, para mantê-la pequena e segura.
FROM alpine:latest

# Copia o binário compilado do estágio 'builder' para o diretório /bin da imagem final.
COPY --from=builder /bin/app /bin/app

# Copia os arquivos de migrations e o .env.example para que estejam disponíveis para o Analyzer no próximo deploy
COPY migrations /migrations
COPY .env.example .env.example

# Expõe a porta 8080, que é a porta padrão da nossa aplicação.
EXPOSE 8080

# Define o comando que será executado quando o container iniciar.
# Isso executa o binário da nossa aplicação.
ENTRYPOINT ["/bin/app"]
