# Go Rate Limiter

A rate limiter middleware for Go web services that controls request flow based on IP address or access token, using Redis for persistence.

## Features

- **IP-based limiting**: Restricts the maximum number of requests per second from a single IP address
- **Token-based limiting**: Restricts requests based on an access token in the `API_KEY` header
- **Token precedence**: Token configurations override IP limits
- **Redis persistence**: Uses Redis for storing counters and block states
- **Strategy pattern**: Easily swappable persistence mechanisms
- **Configurable blocking**: Customizable block duration for exceeded limits
- **HTTP 429 responses**: Returns appropriate status codes and messages

## Architecture

The system follows a clean architecture with separation of concerns:

- **Middleware**: HTTP middleware that intercepts requests
- **Rate Limiter**: Core business logic for rate limiting decisions
- **Persistence Strategy**: Interface for storage operations (Redis implementation included)
- **Configuration**: Environment-based configuration management

## Configuration

All settings are configured via environment variables. Copy `.env.example` to `.env` and adjust as needed:

```env
# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Rate Limiting Configuration
RATE_LIMIT_IP=10          # Max requests per second per IP
RATE_LIMIT_TOKEN=100      # Max requests per second per token
BLOCK_DURATION=5m         # Block duration when limit exceeded

# Server Configuration
SERVER_PORT=8080
```

## Quick Start

### Using Makefile (Recommended)

```bash
# Setup development environment
make dev-setup

# Run application
make run

# Run tests
make test-docker

# Clean up
make clean
```

### Using Docker Compose

1. Clone the repository
2. Copy `.env.example` to `.env` and configure as needed
3. Run with Docker Compose:

```bash
docker-compose up --build
```

The application will be available at `http://localhost:8080`

#### Running Tests with Docker Compose

To run the automated tests using Docker Compose:

```bash
docker-compose -f compose.test.yml up --build --abort-on-container-exit
```

### Using Makefile for Testing

```bash
# Run tests with Docker
make test-docker

# Run tests locally
make test
```

### Local Development

1. Ensure Go 1.26+ is installed
2. Install dependencies:

```bash
go mod download
```

3. Start Redis (or use Docker):

```bash
docker run -p 6379:6379 redis:5
```

4. Run the application:

```bash
go run ./cmd
```

## API Usage

### Test Endpoint

```bash
curl http://localhost:8080/test
```

Response (when allowed):
```json
{"message": "Request allowed"}
```

### Rate Limiting Behavior

- **IP Limiting**: Requests are limited by client IP address
- **Token Limiting**: Include `API_KEY` header to use token-based limits
- **Block Response**: When limit exceeded, returns HTTP 429 with message:
  ```
  "you have reached the maximum number of requests or actions allowed within a certain time frame"
  ```

### Examples

```bash
# IP-based limiting (10 req/s default)
curl http://localhost:8080/test

# Token-based limiting (100 req/s default, overrides IP)
curl -H "API_KEY: mytoken" http://localhost:8080/test
```

## Testing

### Local Testing

Run the test suite locally:

```bash
go test -v ./...
```

Or use Makefile:

```bash
make test
```

### Docker Testing

Run tests using Docker Compose:

```bash
docker-compose -f compose.test.yml up --build --abort-on-container-exit
```

Or use Makefile:

```bash
make test-docker
```

The tests include:
- IP-based rate limiting
- Token precedence over IP limits
- Block behavior when limits exceeded
- Block behavior when limits exceeded

## Extending Persistence Strategies

The system uses the Strategy pattern for persistence. To implement a different storage mechanism:

1. Implement the `PersistenceStrategy` interface:

```go
type PersistenceStrategy interface {
    Increment(ctx context.Context, key string) (int64, error)
    GetTTL(ctx context.Context, key string) (time.Duration, error)
    SetTTL(ctx context.Context, key string, ttl time.Duration) error
    Block(ctx context.Context, key string, duration time.Duration) error
    IsBlocked(ctx context.Context, key string) (bool, error)
}
```

2. Create a new strategy implementation
3. Update the main.go to use your new strategy

## Project Structure

```
.
├── Dockerfile                # Main application container
├── Dockerfile.test           # Test container
├── docker-compose.yml        # Application services
├── compose.test.yml          # Test services
├── Makefile                  # Build and test automation
├── go.mod
├── .env.example
├── README.md
├── cmd/
│   └── main.go
└── internal/
    ├── config/
    │   └── config.go
    ├── middleware/
    │   └── ratelimit.go
    └── ratelimiter/
        ├── limiter.go
        ├── limiter_test.go
        ├── redis.go
        └── strategy.go
```
