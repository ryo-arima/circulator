FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/server cmd/server/main.go
RUN go build -o bin/grpc-server cmd/grpc-server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin ./bin
COPY --from=builder /app/etc ./etc

CMD ["./server"]
