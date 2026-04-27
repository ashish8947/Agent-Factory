# ----------- Build Stage -----------
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o app main.go

# ----------- Run Stage -----------
FROM alpine:latest

WORKDIR /root/

# Install CA certificates (required for HTTPS calls)
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/app .

# Copy .env (optional for local runs)
COPY .env .env

# Run the app
CMD ["./app"]