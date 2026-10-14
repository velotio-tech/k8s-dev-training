# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o k8s-crud .

# Runtime stage
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/k8s-crud .

ENTRYPOINT ["./k8s-crud"]
