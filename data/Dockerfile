FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache wget

COPY . .

RUN go build -o /app/data-bin ./data/cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/data-bin ./

EXPOSE 8081

CMD ["./data-bin"]