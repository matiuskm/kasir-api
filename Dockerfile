# build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

# runtime stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/kasir-api /app/kasirapi

# optional: timezone/certs
RUN apk add --no-cache ca-certificates

ENV PORT=8080
EXPOSE 8080
CMD ["/app/kasir-api"]
