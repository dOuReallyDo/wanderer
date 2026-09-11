FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o wanderer ./cmd/wanderer/

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/wanderer .
EXPOSE 8899
CMD ["./wanderer", "-port", "8899"]