# 🚀 Crypto Server API

A high-performance REST API service for cryptocurrency price tracking and management, built with Go, PostgreSQL, and Redis.

## ✨ Features

- **🔐 JWT Authentication** - Secure user registration and login
- **📊 Real-time Crypto Tracking** - Add and monitor cryptocurrency prices
- **📈 Price History** - Historical price data with Redis caching
- **📉 Statistical Analysis** - Min/max/average prices and change calculations
- **⚡ Auto-updates** - Configurable scheduled price updates
- **🔄 Manual Refresh** - On-demand price refreshing
- **📦 Containerized** - Full Docker support with Docker Compose
- **📈 Monitoring Stack** - Prometheus, Grafana, Loki, and Alertmanager
- **🏥 Health Checks** - Comprehensive health and readiness endpoints
- **🎯 Load Tested** - Production-ready with performance testing

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client App    │───▶│   Crypto API    │───▶│   CoinGecko     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                               │
                               ▼
                    ┌─────────────────┐    ┌─────────────────┐
                    │   PostgreSQL    │    │      Redis      │
                    │   (Users +      │    │  (Price Cache)  │
                    │    Crypto)      │    │                 │
                    └─────────────────┘    └─────────────────┘
```

### Tech Stack

- **Backend**: Go 1.23 with Chi router
- **Database**: PostgreSQL 15 with migrations
- **Cache**: Redis 7 for price history
- **Authentication**: JWT tokens with bcrypt hashing
- **External API**: CoinGecko API for real-time prices
- **Monitoring**: Prometheus + Grafana + Loki stack
- **Deployment**: Docker & Docker Compose

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development)

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd RESTCryptoServer

# Copy environment file
cp .env.example .env

# Edit .env with your values
nano .env
```

### 2. Environment Configuration

```bash
# .env file
SECRET_KEY_CRYPTO_SERVER=your-super-secret-jwt-key-here
DB_DSN=postgres://user:password@db:5432/myapp?sslmode=disable
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
```

### 3. Start Services

```bash
# Start main application stack
docker-compose up -d

# Optional: Start monitoring stack
docker-compose -f docker-compose.monitoring.yml up -d
```

### 4. Verify Installation

```bash
# Check service health
curl http://localhost:8080/health

# Run comprehensive tests
chmod +x tests/test_crypto_server.sh
./tests/test_crypto_server.sh
```

## 📚 API Documentation

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register new user |
| POST | `/auth/login` | Login user |

### Cryptocurrency Endpoints

All endpoints require `Authorization: Bearer <token>` header.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/crypto` | List all tracked cryptocurrencies |
| POST | `/crypto` | Add new cryptocurrency to tracking |
| GET | `/crypto/{symbol}` | Get specific cryptocurrency data |
| PUT | `/crypto/{symbol}/refresh` | Manually refresh price |
| GET | `/crypto/{symbol}/history` | Get price history (last 100 entries) |
| GET | `/crypto/{symbol}/stats` | Get price statistics |
| DELETE | `/crypto/{symbol}` | Remove cryptocurrency from tracking |

### Scheduler Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/schedule` | Get current schedule configuration |
| PUT | `/schedule` | Update schedule settings |
| POST | `/schedule/trigger` | Manually trigger update |

### Health & Monitoring

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Comprehensive health check |
| GET | `/ready` | Readiness probe |
| GET | `/live` | Liveness probe |
| GET | `/metrics` | Prometheus metrics |

## 🔧 Usage Examples

### 1. User Registration & Login

```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secure123"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secure123"}'
```

### 2. Cryptocurrency Management

```bash
# Add Bitcoin
curl -X POST http://localhost:8080/crypto \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"btc"}'

# Get all cryptocurrencies
curl http://localhost:8080/crypto \
  -H "Authorization: Bearer <your-token>"

# Get Bitcoin details
curl http://localhost:8080/crypto/btc \
  -H "Authorization: Bearer <your-token>"

# Get Bitcoin price history
curl http://localhost:8080/crypto/btc/history \
  -H "Authorization: Bearer <your-token>"

# Get Bitcoin statistics
curl http://localhost:8080/crypto/btc/stats \
  -H "Authorization: Bearer <your-token>"
```

### 3. Schedule Management

```bash
# Enable auto-updates every 60 seconds
curl -X PUT http://localhost:8080/schedule \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"enabled":true,"interval_seconds":60}'

# Trigger manual update
curl -X POST http://localhost:8080/schedule/trigger \
  -H "Authorization: Bearer <your-token>"
```

## 📊 Monitoring & Observability

Access the monitoring stack after starting with `docker-compose.monitoring.yml`:

- **Grafana**: http://localhost:3000 (admin/admin123)
- **Prometheus**: http://localhost:9090
- **Alertmanager**: http://localhost:9093

### Key Metrics

- HTTP request rate and latency
- Database connection health
- External API call performance
- Cache hit/miss ratios
- Active cryptocurrency count

## 🧪 Testing

### Run All Tests

```bash
./tests/test_crypto_server.sh
```

### Load Testing

```bash
./tests/load_test.sh
```

### Manual Database Inspection

```bash
# PostgreSQL
docker exec -it $(docker-compose ps -q db) psql -U user -d myapp

# Redis
docker exec -it $(docker-compose ps -q redis) redis-cli
```

## 🛠️ Development

### Local Development Setup

```bash
# Install dependencies
go mod download

# Set up local environment
export SECRET_KEY_CRYPTO_SERVER=your-local-secret
export DB_DSN=postgres://user:password@localhost:5432/myapp?sslmode=disable
export REDIS_HOST=localhost

# Start dependencies only
docker-compose up -d db redis

# Run migrations and start server
go run cmd/server/main.go
```

### Project Structure

```
├── cmd/server/              # Application entry point
├── internal/
│   ├── auth/               # Authentication service
│   ├── crypto/             # Cryptocurrency management
│   ├── db/                 # Database layer (PostgreSQL)
│   ├── redis/              # Cache layer (Redis)
│   ├── coingecko/          # External API client
│   └── updater/            # Scheduled update service
├── monitoring/             # Observability configuration
├── tests/                  # Test scripts
└── docker-compose.yml      # Service orchestration
```

## 🔒 Security Features

- **Password Hashing**: bcrypt with default cost
- **JWT Authentication**: Secure token-based auth
- **Rate Limiting**: 100 requests per minute
- **Input Validation**: Comprehensive request validation
- **CORS Configuration**: Configurable cross-origin support
- **Non-root Container**: Security-focused Docker image

## 📈 Performance

- **Response Times**: Sub-100ms for cached requests
- **Throughput**: 1000+ requests/second
- **Cache Strategy**: Redis with LRU eviction
- **Connection Pooling**: Optimized database connections
- **Graceful Shutdown**: Clean resource cleanup

## 🚀 Deployment

### Production Deployment

1. **Set production environment variables**:
```bash
export ENV=production
export LOG_LEVEL=warn
export SECRET_KEY_CRYPTO_SERVER=<strong-production-secret>
```

2. **Use production-ready database**:
```bash
export DB_DSN=postgres://user:pass@prod-db:5432/crypto_prod?sslmode=require
```

3. **Deploy with Docker**:
```bash
docker-compose up -d
```

### Scaling Considerations

- **Database**: Use read replicas for better performance
- **Cache**: Redis Cluster for high availability
- **Load Balancing**: Multiple API instances behind load balancer
- **Rate Limiting**: Adjust throttle limits based on traffic

## 🔧 Configuration

### Supported Cryptocurrencies

The service supports all cryptocurrencies available on CoinGecko API. Popular ones include:
- Bitcoin (BTC)
- Ethereum (ETH)
- Tether (USDT)
- BNB (BNB)
- Solana (SOL)
- And 10,000+ more...

### Rate Limits

- **API Requests**: 100 requests/minute per connection
- **Price Updates**: Configurable (10-3600 seconds)
- **History Storage**: Last 100 price points per crypto

## 🐛 Troubleshooting

### Common Issues

**Database Connection Failed**:
```bash
# Check if PostgreSQL is running
docker-compose ps db

# Check logs
docker-compose logs db
```

**Redis Connection Failed**:
```bash
# Verify Redis is accessible
docker exec -it $(docker-compose ps -q redis) redis-cli ping
```

**CoinGecko API Errors**:
- Check internet connectivity
- Verify cryptocurrency symbol exists
- Review API rate limits

### Debug Mode

```bash
export LOG_LEVEL=debug
docker-compose up
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [CoinGecko API](https://coingecko.com/api) for cryptocurrency data
- [Chi Router](https://github.com/go-chi/chi) for HTTP routing
- [PostgreSQL](https://postgresql.org) for reliable data storage
- [Redis](https://redis.io) for high-performance caching

---

**Built with ❤️ using Go**