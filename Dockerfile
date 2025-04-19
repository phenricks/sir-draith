# Build stage
FROM golang:1.21-alpine AS builder

# Instalar dependências de build
RUN apk add --no-cache git make

# Configurar diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependência
COPY go.mod go.sum ./

# Download das dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN make build

# Final stage
FROM alpine:latest

# Instalar certificados CA (necessário para HTTPS)
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root
RUN adduser -D -h /app appuser
USER appuser

# Configurar diretório de trabalho
WORKDIR /app

# Copiar o binário compilado do stage anterior
COPY --from=builder /app/bin/sir-draith .
COPY --from=builder /app/configs/config.example.yaml ./configs/config.yaml

# Expor porta se necessário
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./sir-draith"] 