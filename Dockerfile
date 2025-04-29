FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o balancer-app ./cmd/balancer/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/balancer-app .
COPY configs/ configs/
RUN chmod +x /app/balancer-app

CMD ["/app/balancer-app"]
