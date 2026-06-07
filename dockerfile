FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Recommended: Add CGO_ENABLED=0 to ensure the binary runs smoothly on Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]
