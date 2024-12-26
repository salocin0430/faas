FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o faas-api ./cmd/api

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/faas-api .

EXPOSE 8080
CMD ["./faas-api"] 