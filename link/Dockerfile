FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/link-bin ./link/cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/link-bin ./

EXPOSE 8082

CMD ["./link-bin"]