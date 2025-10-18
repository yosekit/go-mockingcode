# MockingCode

Backend-as-a-Service platform for mock data generation and API prototyping.

## Architecture

Microservices-based platform built with Go.

**Services:**
- Gateway (8080) - API Gateway with authentication
- Auth Service (8081) - User authentication and JWT tokens
- Project Service (8082) - Project and schema management
- Data Service (8083) - Mock data CRUD operations

**Databases:**
- PostgreSQL - Users, projects, schemas
- MongoDB - Mock data storage
- Redis - Caching

## Quick Start

```bash
# Start all services
docker-compose -f docker/docker-compose.dev.yml up -d

# Check status
docker-compose -f docker/docker-compose.dev.yml ps

# View logs
docker-compose -f docker/docker-compose.dev.yml logs -f

# Stop services
docker-compose -f docker/docker-compose.dev.yml down
```

## API Endpoints

Base URL: `http://localhost:8080`

**Authentication:**
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh token

**Projects (requires JWT):**
- `GET/POST /api/projects` - List/Create projects
- `GET/PUT/DELETE /api/projects/{id}` - Manage project
- `GET/POST /api/projects/{id}/collections` - List/Create collections
- `GET/PUT/DELETE /api/projects/{id}/collections/{collectionId}` - Manage collection

**Data (requires JWT):**
- `GET/POST /data/{collection}` - List/Create documents
- `GET/PUT/DELETE /data/{collection}/{id}` - Manage document

**Admin:**
- `GET /admin/log-level` - Get current log level
- `PUT /admin/log-level?level=debug` - Change log level (debug|info|warn|error)

## Testing

```bash
./tests/gateway_test.sh
```

## Development

**Structure:**
```
mockingcode/
├── auth/           - Authentication service
├── project/        - Project management service  
├── data/           - Data service
├── gateway/        - API Gateway
├── pkg/            - Shared modules (models, logger)
├── docker/         - Docker configuration
└── tests/          - Integration tests
```

**Commands:**
```bash
# Rebuild specific service
docker-compose -f docker/docker-compose.dev.yml build gateway

# View logs
docker-compose -f docker/docker-compose.dev.yml logs -f gateway

# Run local service
cd auth && go run cmd/server/main.go
```

## Configuration

Environment variables in `docker/.env`:

```env
DB_NAME=mockdb
DB_USER=mockuser
DB_PASSWORD=mockpass
MONGO_PORT=27017
AUTH_PORT=8081
PROJECT_PORT=8082
DATA_PORT=8083
GATEWAY_PORT=8080
JWT_SECRET=your-secret-key
LOG_LEVEL=info
LOG_FORMAT=text
```

## Documentation

Swagger UI available for each service:
- `http://localhost:8080/swagger/` - Gateway
- `http://localhost:8081/swagger/` - Auth
- `http://localhost:8082/swagger/` - Project
- `http://localhost:8083/swagger/` - Data

## Ports

- 8080 - Gateway
- 8081 - Auth Service
- 8082 - Project Service
- 8083 - Data Service
- 5432 - PostgreSQL
- 27017 - MongoDB
- 6379 - Redis

## Features

- JWT-based authentication
- API Gateway Pattern with centralized auth
- Structured logging with dynamic log levels
- Docker-based deployment
- RESTful API
- Mock data generation
- Go workspace monorepo

## License

MIT

