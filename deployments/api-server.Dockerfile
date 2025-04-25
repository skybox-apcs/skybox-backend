FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/api-server

EXPOSE 8080

CMD ["go", "run", "main.go"]

