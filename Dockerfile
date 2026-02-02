# Stage 1: Build Vue.js frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web

# Copy package files
COPY web/package*.json ./

# Install dependencies
RUN npm ci

# Copy source files
COPY web/ ./

# Build frontend
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source files
COPY . .

# Copy built frontend
COPY --from=frontend-builder /app/web/dist ./web/dist

# Generate go.sum and build the binary
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o listenbucket ./cmd/listenbucket

# Stage 3: Final image
FROM alpine:3.19

# Install runtime dependencies
# Note: Install yt-dlp via pip to get the latest version (Alpine's package is often outdated)
RUN apk add --no-cache \
    ffmpeg \
    python3 \
    py3-pip \
    ca-certificates \
    tzdata \
    && pip3 install --break-system-packages --no-cache-dir yt-dlp

# Create non-root user
RUN adduser -D -h /app listenbucket

WORKDIR /app

# Copy binary
COPY --from=backend-builder /app/listenbucket .

# Create data directory
RUN mkdir -p /data && chown -R listenbucket:listenbucket /data

# Switch to non-root user
USER listenbucket

# Environment variables
ENV PORT=8080
ENV DATA_DIR=/data
ENV BASE_URL=http://localhost:8080

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/feeds || exit 1

# Run
CMD ["./listenbucket"]
