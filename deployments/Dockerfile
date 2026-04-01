# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Set working directory di dalam container
WORKDIR /app

# Copy go.mod dan go.sum dari ROOT project ke dalam container
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code dari ROOT project ke container
COPY . .

# Build binary (Perhatikan path cmd/server/main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta

# Ambil file binary dari stage builder
COPY --from=builder /app/main .

# Expose port backend Anda
EXPOSE 8080

# Jalankan server
CMD ["./main"]