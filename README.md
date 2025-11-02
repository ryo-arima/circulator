# circulator

**⚠️ This is an experimental and educational application for learning purposes. It is not intended for production use.**

Circulator with gRPC and REST API support.

## Architecture

This project supports both gRPC and REST APIs:

- **REST API Server**: HTTP/JSON based API using Gin framework
- **gRPC Server**: High-performance gRPC API using Protocol Buffers
- **CLI Client**: Command-line interface for management operations

## Usage

### REST API Server
```bash
go run cmd/server/main.go
# or
make run-server
```

### gRPC Server
```bash
go run cmd/grpc-server/main.go
# or  
make run-grpc-server
```

### CLI Client
```bash
go run cmd/client/main.go --help
# or
make run-client
```

## Development

### Protocol Buffers

Generate protobuf files:
```bash
make proto
```

Install protobuf tools:
```bash
make proto-install
```

### Build

Build all binaries:
```bash
make build
```

### Docker

Start all services:
```bash
docker-compose up -d
```

## Structure

This project follows a clean architecture with:

- **REST API**: `pkg/server/` - HTTP/JSON API using Gin framework  
- **gRPC API**: `pkg/agent/` - High-performance gRPC services
- **CLI Client**: `pkg/client/` - Command-line management interface
- **Protocol Buffers**: `pkg/agent/proto/` - gRPC service definitions
- **Entities**: `pkg/entity/` - Domain models and data transfer objects
- **Configuration**: `pkg/config/` - Application configuration management

### Architecture Layers

- **Controllers**: Handle HTTP/gRPC requests
- **Use Cases**: Business logic layer  
- **Repositories**: Data access layer
- **Models**: Domain entities and database models
# circulator
