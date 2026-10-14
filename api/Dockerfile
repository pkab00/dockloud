FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/api-bin ./api/cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-bin ./

EXPOSE 8080

CMD ["./api-bin"]