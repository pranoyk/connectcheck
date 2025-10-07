FROM golang:1.20.1-alpine3.17 AS builder

WORKDIR /app

# Copy go source
COPY main.go .

# Build the application
RUN go build -o server main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]