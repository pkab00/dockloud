FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6
# RUN swag init -g cmd/app/main.go -o docs

RUN go build -o /app/data ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/data .

EXPOSE 8081

CMD ["./data"]