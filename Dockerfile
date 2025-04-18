FROM golang:1.20-alpine

WORKDIR /app

# Copiază modulele Go
COPY relationship-helix-backend/go.mod relationship-helix-backend/go.sum ./
RUN go mod download

# Copiază codul sursă
COPY relationship-helix-backend/ ./

# Compilează aplicația
RUN go build -o main ./cmd/server

# Expune portul
EXPOSE 8080

# Comandă de start
CMD ["./main"]
