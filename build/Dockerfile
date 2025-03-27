FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o skybox-backend .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/skybox-backend .
CMD ["./skybox-backend"]
