# ==============================
# Stage 1: Build Go binary
# ==============================
FROM golang:1.25-alpine AS builder

# Install dependency yang diperlukan untuk build Go
RUN apk add --no-cache git openssh tzdata build-base python3 net-tools

# Set working directory
WORKDIR /app

# Copy go.mod & go.sum terlebih dahulu (agar dependency cache)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy semua source code
COPY . .

# Install tool tambahan (opsional)
RUN go install github.com/buu700/gin@latest

# Build binary lewat Makefile
RUN make build


# ==============================
# Stage 2: Final runtime image
# ==============================
FROM alpine:latest

# Install dependency ringan untuk runtime
RUN apk add --no-cache tzdata curl

# Set timezone (opsional)
ENV TZ=Asia/Jakarta

# Buat direktori kerja
WORKDIR /app

# Copy hasil build dari stage builder
COPY --from=builder /app /app

# Copy file .env ke container
COPY .env .env

# Expose port aplikasi
EXPOSE 8001

# Jalankan aplikasi
ENTRYPOINT ["/app/user-service"]
