FROM golang:1.23.3-alpine AS migrator
WORKDIR /app
COPY . .
RUN go mod download
CMD ["go", "run", "cmd/migrator/main.go"] 