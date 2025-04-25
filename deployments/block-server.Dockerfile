FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/block-server

EXPOSE 9090

CMD ["go", "run", "main.go"]

