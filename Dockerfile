# Usa uma imagem base do Go
FROM golang:1.21-alpine

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos de módulo e baixa as dependências
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copia o resto do código fonte do projeto
COPY . .

# Compila a aplicação
# O -o ./out/api constrói o executável 'api' dentro de uma pasta 'out'
RUN go build -o ./out/api ./cmd/api

# Expõe a porta que o nosso servidor Go usa
EXPOSE 8080

# Comando para executar a aplicação quando o contêiner iniciar
CMD [ "./out/api" ]
