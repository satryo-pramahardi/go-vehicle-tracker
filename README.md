# Go Vehicle Tracker

An MQTT-based real-time fleet tracking system with geofence monitoring built with Go, featuring microservices architecture, event-driven design, and modern DevOps practices.

## Architecture Overview

This project demonstrates a **Clean Architecture** implementation with **Domain-Driven Design** principles:

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   HTTP API  │  │   MQTT      │  │  RabbitMQ   │         │
│  │   (Gin)     │  │  Publisher  │  │  Consumer   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Location    │  │ Geofence    │  │ Event       │         │
│  │ Worker      │  │ Service     │  │ Service     │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                   Domain Layer                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Vehicle     │  │ Geofence    │  │ Event       │         │
│  │ Location    │  │ Models      │  │ Models      │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ PostgreSQL  │  │   Redis     │  │  RabbitMQ   │         │
│  │ Repository  │  │   Cache     │  │  Message    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

## Key Features

- **Real-time Vehicle Tracking** - GPS location monitoring with Redis caching
- **Geofence Detection** - Automatic detection of vehicle entry/exit events
- **Geofence Alert Publication** - Real-time alert publishing via RabbitMQ
- **Event-Driven Architecture** - Asynchronous processing with RabbitMQ
- **Microservices Design** - Independent services for different concerns
- **RESTful API** - HTTP endpoints for data access and management
- **MQTT Integration** - IoT device communication support
- **Docker Containerization** - Complete containerized deployment
- **PostgreSQL Database** - Reliable data persistence
- **Redis Caching** - High-performance data caching

## Technology Stack

### Backend
- **Go 1.23** - High-performance server-side language
- **Gin** - Fast HTTP web framework
- **GORM** - Object-relational mapping
- **PostgreSQL** - Primary database
- **Redis** - Caching and message queuing
- **RabbitMQ** - Message broker for event processing

### Infrastructure
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **MQTT** - IoT communication protocol

### Architecture Patterns
- **Clean Architecture** - Separation of concerns
- **Domain-Driven Design** - Business logic modeling
- **Event-Driven Architecture** - Asynchronous processing
- **Repository Pattern** - Data access abstraction
- **Service Layer Pattern** - Business logic encapsulation

## Project Structure

```
go-vehicle-tracker/
├── cmd/                          # Application entry points
│   ├── api/                     # HTTP API server
│   ├── worker/                  # Background location processor
│   ├── rabbitmq_consumer/       # Message queue consumer
│   ├── publisher/               # MQTT message publisher
│   ├── subscriber/              # MQTT message subscriber
│   ├── migrate/                 # Database migration tool
│   └── sandbox/                 # Development/testing utilities
├── internal/                    # Private application code
│   ├── app/                     # Application layer
│   │   ├── service/             # Business logic services
│   │   │   ├── location_worker.go    # Location processing
│   │   │   ├── geofence_service.go   # Geofence detection
│   │   │   ├── rabbitmq_service.go   # Message publishing
│   │   │   └── event_service.go      # Event handling
│   │   └── usecase/             # Application use cases
│   ├── delivery/                # Delivery layer
│   │   ├── http/                # HTTP handlers and routes
│   │   └── mqtt/                # MQTT client and handlers
│   ├── repository/              # Data access layer
│   │   ├── interfaces.go        # Repository interfaces
│   │   └── postgres/            # PostgreSQL implementations
│   ├── model/                   # Domain models
│   │   ├── location.go          # Vehicle location model
│   │   ├── geofence.go          # Geofence model
│   │   └── event_log.go         # Event logging model
│   ├── db/                      # Database connections
│   │   └── gorm.go              # GORM configuration
│   └── geo/                     # Geographic utilities
│       └── haversine.go         # Distance calculation
├── docker/                      # Docker configurations
│   ├── mqtt/                    # MQTT broker config
│   ├── migration/               # Database migration
│   └── subscriber/              # MQTT subscriber
├── docker-compose.yaml          # Multi-service orchestration
├── go.mod                       # Go module definition
└── README.md                    # Project documentation
```

## Design Decisions

### **Clean Architecture Implementation**
- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Separation of Concerns**: Clear boundaries between layers
- **Testability**: Easy to unit test business logic
- **Maintainability**: Changes in one layer don't affect others

### **Event-Driven Design**
- **Asynchronous Processing**: Non-blocking event handling
- **Loose Coupling**: Services communicate via events
- **Scalability**: Easy to add new event consumers
- **Reliability**: Event persistence and retry mechanisms

### **Microservices Approach**
- **Single Responsibility**: Each service has one clear purpose
- **Independent Deployment**: Services can be deployed separately
- **Technology Flexibility**: Different services can use different tech stacks
- **Team Autonomy**: Teams can work on different services

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.23+

### Running the Application
```bash
# Clone the repository
git clone <repository-url>
cd go-vehicle-tracker

# Start all services
docker compose up -d

# Run database migrations
docker compose run --rm migrate

# Check service status
docker compose ps
```

### Accessing Services
- **API Server**: http://localhost:8080
- **RabbitMQ Management**: http://localhost:15672 (admin/password)
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## Data Flow Diagrams

### Vehicle Location Data Flow

The vehicle location data flows through multiple services for processing and storage:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   MQTT Device   │───▶│   MQTT Broker   │───▶│  Location       │───▶│   PostgreSQL    │
│   (Publisher)   │    │   (Mosquitto)   │    │   Worker        │    │   (vehicle_     │
│                 │    │                 │    │                 │    │   locations)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │     Redis       │    │  Geofence       │
                       │   (Cache)       │    │  Service        │
                       └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
                                               ┌─────────────────┐
                                               │   RabbitMQ      │
                                               │   (geofence.    │
                                               │   event)        │
                                               └─────────────────┘
```

#### Vehicle Location Flow Steps:

1. **MQTT Publisher** sends vehicle location data to MQTT broker
2. **MQTT Broker** (Mosquitto) receives and distributes messages
3. **Location Worker** processes location data from Redis queue
4. **PostgreSQL** stores the location in `vehicle_locations` table
5. **Redis** caches frequently accessed location data
6. **Geofence Service** checks if location triggers geofence events
7. **RabbitMQ** receives geofence alerts for further processing

#### Example Sequence Flow:

```
Vehicle sends location → MQTT Subscriber receives → Location Worker processes → 
Geofence Service checks → Insert to DB and send event → RabbitMQ Consumer alerts
```

**Detailed Example:**
1. **Vehicle TJ001** sends location `(-6.193125, 106.820233)` via MQTT
2. **MQTT Subscriber** receives and stores in Redis queue
3. **Location Worker** processes the location and saves to PostgreSQL
4. **Geofence Service** detects vehicle entered "Bundaran HI" geofence
5. **Event Service** creates geofence entry event and publishes to RabbitMQ
6. **RabbitMQ Consumer** receives alert and logs: "🚨 GEOFENCE ALERT RECEIVED!"

### Event Log Data Flow

The `event_log` serves as the **central messaging format** that unifies all system events into a standardized structure. This design provides consistency, auditability, and extensibility across the entire system.

#### Central Event Log Model

All system events follow a standardized `EventEnvelope` structure:

```go
type EventEnvelope struct {
    EventType string          `json:"event_type"`    // Type of event (e.g., "geofence_entry", "location_update")
    Timestamp time.Time       `json:"timestamp"`     // When the event occurred
    Payload   json.RawMessage `json:"payload"`       // Event-specific data (flexible JSON)
    Source    string          `json:"source"`        // Service that generated the event
}
```

#### Event Log Flow

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Various       │───▶│   Event         │───▶│   Redis         │───▶│   PostgreSQL    │
│   Services      │    │   Service       │    │   (event_log:   │    │   (event_logs)  │
│                 │    │                 │    │   queue)        │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
        │                        │                        │
        ▼                        ▼                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Geofence      │    │   Location      │    │   API           │
│   Events        │    │   Worker        │    │   Requests      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

#### Event Log Flow Steps:

1. **Various Services** generate events (geofence, location, API, etc.)
2. **Event Service** formats and enriches event data into `EventEnvelope`
3. **Redis Queue** temporarily stores events for processing
4. **PostgreSQL** permanently stores events in `event_logs` table
5. **Event Logs** provide audit trail and system monitoring

#### Benefits of Central Event Log Model

**Consistency**: All events follow the same structure, making them easy to process and analyze
**Extensibility**: New event types can be added without changing the core structure
**Auditability**: Complete trace of all system activities with standardized format
**Flexibility**: JSONB payload allows for event-specific data while maintaining structure
**Monitoring**: Unified format enables centralized monitoring and alerting
**Debugging**: Standardized structure makes troubleshooting easier across services

### Database Schema Overview

#### Vehicle Locations Table
```sql
CREATE TABLE vehicle_locations (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id VARCHAR(50) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    INDEX idx_vehicle_timestamp (vehicle_id, timestamp)
);
```

#### Event Logs Table
```sql
CREATE TABLE event_logs (
    id BIGSERIAL PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    source VARCHAR(100) NOT NULL,
    payload JSONB,
    timestamp TIMESTAMP NOT NULL,
    INDEX idx_event_type_timestamp (event_type, timestamp)
);
```

#### Geofence Events Table
```sql
CREATE TABLE geofence_events (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id VARCHAR(50) NOT NULL,
    geofence_id BIGINT NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    INDEX idx_vehicle_geofence (vehicle_id, geofence_id)
);
```

## Monitoring and Logs

```bash
# View all service logs
docker compose logs -f

# View specific service logs
docker compose logs -f workers
docker compose logs -f rabbitmq-consumer
docker compose logs -f api
```

## Testing the System

### Using the MQTT Publisher

The system includes a **vehicle location publisher** that simulates a moving vehicle and sends location updates via MQTT. This is perfect for testing the geofence functionality.

#### Basic Usage

```bash
# Start the publisher with default settings
docker compose run --rm publisher

# This will:
# - Connect to MQTT broker at mqtt:1883
# - Send location updates every 2 seconds
# - Simulate vehicle TJ001 moving around Bundaran HI
# - Continue until interrupted (Ctrl+C)
```

#### Advanced Usage

The publisher supports various command-line flags for customization:

```bash
# Custom vehicle ID and interval
docker compose run --rm publisher --vehicle-id BUS001 --interval 5

# Simulate a specific number of messages
docker compose run --rm publisher --count 10 --interval 1

# Custom trip parameters
docker compose run --rm publisher \
  --vehicle-id TEST001 \
  --speed 10.0 \
  --trip-length 500 \
  --interval 3

# Use a different MQTT broker
docker compose run --rm publisher --broker tcp://localhost:1883
```

#### Publisher Options

| Flag | Default | Description |
|------|---------|-------------|
| `--vehicle-id` | `TJ001` | Vehicle identifier |
| `--interval` | `2` | Seconds between messages |
| `--count` | `0` | Number of messages (0 = infinite) |
| `--speed` | `5.0` | Vehicle speed in m/s |
| `--trip-length` | `200` | Trip length in meters before turning |
| `--broker` | `mqtt:1883` | MQTT broker URL |

#### Testing Geofence Alerts

1. **Start the publisher**:
   ```bash
   docker compose run --rm publisher --vehicle-id TEST001 --interval 1
   ```

2. **Monitor geofence alerts**:
   ```bash
   # In another terminal
   docker compose logs -f rabbitmq-consumer
   ```

3. **Watch for alerts** when the vehicle enters the geofence area around Bundaran HI.

#### Manual Testing

You can also send individual location messages manually:

```bash
# Vehicle enters geofence (Bundaran HI coordinates)
mosquitto_pub -h localhost -t "fleet/vehicle/TJ001/location" \
  -m '{"vehicle_id":"TJ001","latitude":-6.193125,"longitude":106.820233,"speed":5.0,"timestamp":"2025-07-12T10:00:00Z"}'

# Vehicle outside geofence
mosquitto_pub -h localhost -t "fleet/vehicle/TJ001/location" \
  -m '{"vehicle_id":"TJ001","latitude":-6.200000,"longitude":106.820233,"speed":5.0,"timestamp":"2025-07-12T10:00:00Z"}'
```

#### Query Vehicle Data

```bash
# Get latest location via API
curl http://localhost:8080/api/vehicles/TJ001/location

# Check geofence events in database
docker exec postgres psql -U admin -d vehicle_tracker \
  -c "SELECT * FROM geofence_events ORDER BY timestamp DESC LIMIT 5;"
```

#### Using Makefile for Testing

The project includes a Makefile with convenient testing commands:

```bash
# Test MQTT connection
make test-mqtt

# Test geofence functionality (sends test location)
make test-geofence

# Check service health
make health

# View all service logs
make docker-logs

# View specific service logs
make docker-logs-consumer
```

## Development

### Adding New Features
1. **Domain Models**: Add to `internal/model/`
2. **Business Logic**: Add to `internal/app/service/`
3. **Data Access**: Add to `internal/repository/`
4. **API Endpoints**: Add to `internal/delivery/http/`
5. **New Services**: Add to `cmd/`

### Code Quality
- Follow Go conventions and best practices
- Use meaningful variable and function names
- Add comprehensive error handling
- Include unit tests for business logic
- Document complex algorithms and decisions

## Performance Considerations

- **Redis Caching**: Reduces database load for frequently accessed data
- **Connection Pooling**: Efficient database connection management
- **Asynchronous Processing**: Non-blocking event handling
- **Indexed Queries**: Optimized database queries with proper indexing
- **Container Resource Limits**: Proper resource allocation for containers

## Deployment

### Production Considerations
- **Environment Variables**: Configure production settings
- **SSL/TLS**: Enable HTTPS for API endpoints
- **Monitoring**: Add application monitoring and alerting
- **Backup Strategy**: Implement database backup procedures
- **Load Balancing**: Add load balancer for high availability

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

Built as a portfolio project to demonstrate:
- **Go development** best practices
- **Microservices architecture** design
- **Event-driven systems** implementation
- **DevOps practices** with Docker
- **Clean code** and **maintainable architecture**

## AI Assistance Disclosure

This project was developed with the assistance of AI tools as part of a modern, human-centered software development workflow.

### AI Contributions

- **Code Generation**: Helped scaffold the initial project structure and generate boilerplate.
- **Debugging & Optimization**: Assisted in identifying issues and improving code clarity and performance.
- **Best Practices**: Provided guidance on idiomatic Go patterns, concurrency, and clean architecture.
- **Documentation**: Supported the creation of this README, architectural explanations, and ADR-style rationale.

### Human Contributions

- **Design & Architecture**: All core system and architectural decisions were made by the developer.
- **Business Logic Implementation**: All application-specific logic was written manually based on problem understanding.
- **Customization & Integration**: External tool and library integration was tailored by the developer for this specific use case.
- **Comprehension & Learning**: AI outputs were critically evaluated, refined, and used as part of an active learning process.

### Intent

This project demonstrates how AI tools can enhance developer productivity without replacing critical thinking, ownership, or system design skills.  
It represents a collaborative effort between human developer and machine assistant — with all intellectual accountability retained by the author.

---

*This project showcases modern backend development practices and is suitable for production use with proper configuration and security measures.*
