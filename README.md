# Go Vehicle Tracker

A real-time fleet tracking system demonstrating **microservices architecture**, **event-driven design**, and **clean code principles** in Go. Built for scalability and maintainability with modern DevOps practices.

**Objective**: This project was created to learn Go while sharpening system design skills through hands-on implementation of microservices architecture and event-driven patterns.

## Architecture & Design

### System Overview
```
                      ┌─────────────────┐
                      │   IoT Devices   │
                      │   (GPS Units)   │
                      └─────────────────┘
                              │
                              ▼ MQTT
                      ┌─────────────────┐
                      │  Subscriber     │
                      │  Service (MQTT) │
                      └─────────────────┘
                              │
                              ▼
                      ┌─────────────────┐
                      │   Redis Queue   │
                      │   (Buffering)   │
                      └─────────────────┘
                              │
                              ▼
       ┌──────────────────────────────────────────────┐
       │               Worker Services                │
       │  ┌────────────────────┐   ┌────────────────┐ │
       │  │ Location Processor │   │ Event Log      │ │
       │  │ + Geofence Logic   │   │ Worker         │ │
       │  └────────────────────┘   └────────────────┘ │
       └──────────────────────────────────────────────┘
                              │
                              ▼
       ┌──────────────────────────────────────────────┐
       │               PostgreSQL Database            │
       │  ┌────────────────────┐  ┌─────────────────┐ │
       │  │ Vehicle Location   │  │ Event Log       │ │
       │  │ Table              │  │ Table           │ │
       │  └────────────────────┘  └─────────────────┘ │
       └──────────────────────────────────────────────┘
                              │
                              ▼
                      ┌─────────────────┐
                      │   HTTP API      │
                      │   (Gin)         │
                      └─────────────────┘
                              │
                              ▼
                      ┌─────────────────┐
                      │   Swagger UI    │
                      │   (Docs)        │
                      └─────────────────┘

       ┌─────────────────┐
       │ Dead Letter     │ ← Unprocessed Events or Failures
       │ Queue (Redis)   │
       └─────────────────┘
```

### Data Flow
```
1. IoT Devices publish GPS → Subscriber Service via MQTT
2. Subscriber pushes data → Redis Queue
3. Redis feeds → Worker Services:
    a. Location Worker → Vehicle Location Table
    b. Location Worker → Geofence Logic
    c. Geofence Logic → Event Log Table & Messaging (RabbitMQ)
    d. Event Log Worker → Event Log Table
4. Worker failures → Event Log Worker
5. Log Worker failure → Dead Letter Queue (Redis)
6. HTTP API queries → PostgreSQL Tables
7. Swagger UI reads from → HTTP API

*Optional*: Dead Letter Queue → Log Processor, Alerting, Analytics
```

### Project Flow Description

The system operates as a real-time event processing pipeline designed for high-throughput IoT data ingestion. GPS devices continuously publish location updates via MQTT to the Publisher Service, which acts as the entry point for all vehicle telemetry data. The Publisher Service buffers incoming messages in Redis queues to handle traffic spikes and ensure no data loss during processing bottlenecks.

The Worker Service operates as a stateless consumer that processes location events from Redis queues. It performs data validation, geofence boundary detection, and persists processed data to PostgreSQL for long-term storage. The worker also triggers real-time alerts via RabbitMQ when vehicles enter or exit predefined geofence zones.

The HTTP API serves as the primary interface for client applications, providing RESTful endpoints for querying vehicle locations and historical data. The API layer implements clean architecture principles with clear separation between delivery, application, and infrastructure concerns. All data retrieval operations query PostgreSQL directly, ensuring consistency and leveraging the database's query optimization capabilities.

### Fault Tolerance Design

The system handles failures through **Redis persistence** and **dead letter queues**. When the **Location Worker** fails to process data, it sends error events to the **Event Log Worker** via Redis. If the **Event Log Worker** itself fails, failed events are pushed to a **dead letter queue** for later retry. The **ArchiveDeadLetterWorker** continuously processes dead letter entries, moving them to a permanent failed list for manual investigation.

**Service isolation** prevents cascading failures by allowing each component to operate independently. If **RabbitMQ** becomes unavailable, geofence detection continues and events are logged locally. **PostgreSQL connection pooling** and automatic retry mechanisms handle transient database failures. The **stateless worker design** enables horizontal scaling - multiple worker instances can process from the same Redis queues without coordination.

### Data Design

**Redis** serves as a high-performance message queue for real-time data streams, while **PostgreSQL** provides durable storage for processed data. The **Location Worker** consumes from Redis queues and persists validated location data to the **VehicleLocation table**. The **Event Log Worker** handles system events and errors, storing them in the **EventLog table** with timestamps and source identification for audit trails.

**Data consistency** is achieved through eventual consistency patterns where real-time data flows through Redis queues before being committed to PostgreSQL. The **Repository Pattern** abstracts data access, enabling easy testing and potential database migrations. **Event sourcing** principles are applied through the EventLog model, capturing all system events as immutable records that enable debugging and potential event replay for system recovery.

### Key Design Decisions

**Clean Architecture**
- **Domain Layer**: Pure business logic, no external dependencies
- **Application Layer**: Use cases and orchestration
- **Infrastructure Layer**: Database, messaging, external services
- **Delivery Layer**: HTTP, MQTT, message queues

**Event-Driven Processing**
- MQTT → Redis → Worker pipeline for real-time processing
- Asynchronous geofence detection and alerting
- Decoupled services for independent scaling

**Microservices Benefits**
- **Location Service**: GPS data processing and caching
- **Geofence Service**: Boundary detection and alerts
- **Event Service**: Audit trail and logging
- **API Gateway**: RESTful endpoints and documentation

## Technology Stack

**Backend**: Go 1.23, Gin, GORM, PostgreSQL, Redis, RabbitMQ  
**Infrastructure**: Docker, Docker Compose, MQTT  
**Documentation**: Swagger/OpenAPI  
**Testing**: Unit tests with mocks, integration tests  

## Project Structure

```
go-vehicle-tracker/
├── cmd/                  # Application entry points
│   ├── api/              # HTTP API server
│   ├── worker/           # Background processor
│   └── publisher/        # MQTT publisher
├── internal/             # Private application code
│   ├── app/              # Business logic & services
│   ├── delivery/         # HTTP/MQTT handlers
│   ├── repository/       # Data access layer
│   ├── model/            # Domain models
│   └── geo/              # Geographic utilities
├── tests/                
│   └── integration/      # Integration tests
├── docs/                 # Swagger documentation
└── docker/               # Container configurations
```

## Quick Start

```bash
# Start all services
make docker-up

# Run migrations
make migrate

# Access API docs
open http://localhost:8080/swagger/index.html
```

## API Documentation

### Interactive Documentation
Access the full interactive API documentation at: http://localhost:8080/swagger/index.html

### Key Endpoints

#### Get Latest Vehicle Location
```http
GET /vehicles/{vehicle_id}/location
```

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 123,
    "vehicle_id": "TRUCK-001",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

#### Get Vehicle Location History
```http
GET /vehicles/{vehicle_id}/history?start=2024-01-15T00:00:00Z&end=2024-01-15T23:59:59Z
```

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "vehicle_id": "TRUCK-001",
    "start_time": "2024-01-15T00:00:00Z",
    "end_time": "2024-01-15T23:59:59Z",
    "count": 1440,
    "locations": [...]
  }
}
```

## Key Features

- **Real-time GPS Tracking** with Redis caching
- **Geofence Detection** with automatic alerts
- **Event-Driven Architecture** for scalability
- **RESTful API** with Swagger documentation
- **MQTT Integration** for IoT devices
- **Clean Architecture** for maintainability

## Testing Strategy

The project implements a comprehensive testing strategy with multiple layers of test coverage:

**Unit Tests**
- **Handler Tests**: HTTP endpoint testing with mocked repositories using testify/mock
- **Service Tests**: Business logic testing for geofence detection and event processing
- **Repository Tests**: Data access layer testing with mocked database connections
- **Mock Strategy**: Uses testify/mock for dependency injection and isolated testing

**Integration Tests**
- **API Integration**: End-to-end HTTP API testing with real database connections
- **Event Log Testing**: Database integration tests for event logging functionality
- **Data Persistence**: Tests for vehicle location and event log data persistence
- **Error Handling**: Tests for dead letter queue and error recovery mechanisms

**Test Coverage**
- **Coverage Reporting**: Comprehensive coverage analysis with HTML reports
- **Race Detection**: Concurrent code testing with Go's race detector
- **Test Categories**: Organized tests by layer (unit, integration, e2e)

**Testing Commands**
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run with race detection
make test-race

# Run only HTTP tests
make test-http

# Run integration tests
make test-integration
```

## Project Intent

This project demonstrates:
- **System Design Skills**: Microservices and event-driven architecture
- **Clean Code**: SOLID principles and clean architecture
- **DevOps Practices**: Docker, testing, documentation
- **Real-world Patterns**: Repository pattern, service layer, dependency injection

## AI Assistance Disclosure

This project was built as a learning and portfolio piece with **AI collaboration**.  
I retained full control over the system architecture and business logic.

### My Contributions:
- System design and service orchestration
- Technology stack and architecture decisions
- Clean architecture implementation
- Business logic and data flow
- API and testing strategy

### AI Assisted With:
- Boilerplate code generation (handlers, services, repositories)
- Docker Compose and container setup
- Swagger / OpenAPI integration
- Testing scaffolds with mocks and integration setup
- Error handling patterns

---

**Let's Connect!** [LinkedIn](https://www.linkedin.com/in/satryo-pramahardi/) - Feel free to reach me!
