FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/document-bin ./document/cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/document-bin ./

EXPOSE 8080

CMD ["./document-bin"]