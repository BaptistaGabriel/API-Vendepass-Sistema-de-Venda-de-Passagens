# Usar imagem base do Go
FROM golang:1.20-alpine

# Criar diretório de trabalho
WORKDIR /app

# Copiar o arquivo go.mod antes de copiar o código
COPY go.mod ./
RUN go mod download

# Copiar o restante dos arquivos (incluindo src)
COPY ./src .

# Compilar o binário
RUN go build -o server .

# Expor a porta 8080
EXPOSE 8080

# Rodar o servidor
CMD ["./server"]
