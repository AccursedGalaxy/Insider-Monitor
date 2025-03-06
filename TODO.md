# TODO

## Upgrade Plan:
**Phase 1:** Create API endpoints in the current application
Add basic CRUD endpoints for configuration
Implement configuration updating while running
Add authentication for sensitive operations
    **=> Completed**

**Phase 2:** Refactor to separate backend and frontend
Split codebase into backend and frontend services
Implement proper API communication
Set up WebSockets for real-time updates

Key considerations:
1. Architecture separation
2. API design
3. Real-time communication
4. Single command startup
5. Performance optimization
6. Data handling concerns</think># Phase 2 Implementation Plan: Refactoring Insider-Monitor into Separate Backend and Frontend

## 1. Architecture Design

### Backend Service
```
backend/
├── cmd/
│   └── server/
│       └── main.go        # Backend entry point
├── internal/
│   ├── api/               # API handlers and routes
│   │   ├── middleware/    # Auth and request middleware
│   │   ├── handlers/      # API endpoint handlers
│   │   └── routes.go      # Route definitions
│   ├── config/            # Existing config package
│   ├── monitor/           # Existing monitor package
│   ├── alerts/            # Existing alerts package
│   ├── storage/           # Existing storage package
│   └── websocket/         # New WebSocket implementation
└── pkg/
    └── models/            # Shared data models
```

### Frontend Service
```
frontend/
├── cmd/
│   └── server/
│       └── main.go        # Frontend server entry point
├── internal/
│   └── server/            # Static file server
├── src/                   # Frontend source code
│   ├── components/        # React components
│   ├── services/          # API client services
│   ├── store/             # State management
│   └── pages/             # Application pages
└── public/                # Static assets
```

### Startup Orchestrator
```
cmd/
└── monitor/
    └── main.go            # Main entry point that starts both services
```

## 2. Backend Implementation

### API Development
1. **Create RESTful API endpoints**:
   - Wallet data endpoints: `/api/v1/wallets`
   - Configuration endpoints: `/api/v1/config`
   - Alert management: `/api/v1/alerts`
   - Status and health: `/api/v1/status`

2. **Authentication and Security**:
   - JWT authentication middleware
   - API key support for machine-to-machine communication
   - Rate limiting to prevent abuse

3. **Request Handling**:
   - Use Goroutines to handle concurrent API requests
   - Implement graceful error handling and response standardization
   - Add request logging and metrics

### WebSocket Implementation
1. **Create WebSocket server**:
   - Use Go's gorilla/websocket or nhooyr.io/websocket library
   - Implement connection management with concurrent-safe structures
   - Add ping/pong heartbeat mechanism

2. **Define Message Protocol**:
   - Create structured message types for different events
   - Implement efficient binary encoding using Protocol Buffers or MessagePack
   - Add compression for large data payloads

3. **Event Broadcasting**:
   - Design a pub/sub system for real-time updates
   - Use channels to efficiently distribute messages to connected clients
   - Implement fan-out pattern for broadcasting to multiple clients

### Performance Optimization
1. **Efficient Data Handling**:
   - Use Go's sync.Pool for frequently allocated objects
   - Implement request batching for high-volume operations
   - Utilize buffered channels for asynchronous processing

2. **Connection Management**:
   - Implement connection pooling for blockchain RPC calls
   - Add circuit breakers to prevent cascading failures
   - Use context.Context for request cancellation and timeouts

## 3. Frontend Implementation

1. **Modern Frontend Stack**:
   - Use React for component-based UI
   - Add TypeScript for type safety
   - Implement responsive design with TailwindCSS

2. **State Management**:
   - Use React Query for API data fetching and caching
   - Implement context or Redux for global state
   - Add optimistic UI updates for better user experience

3. **Real-time Updates**:
   - Create WebSocket client service
   - Implement reconnection logic with exponential backoff
   - Add message queue for offline operation

4. **Performance**:
   - Implement code splitting and lazy loading
   - Add virtualized lists for large data sets
   - Use memoization for expensive computations

## 4. Integration and Orchestration

### Single Command Startup
1. **Main Executable**:
   - Create a master process that spawns both backend and frontend services
   - Implement clean shutdown handling for both services
   - Add health checks to ensure both services are running

2. **Orchestration Logic**:
```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Set up signal handling for graceful shutdown
    signalCh := make(chan os.Signal, 1)
    signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

    // Configuration shared between services
    cfg, err := config.LoadConfig("config.json")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Start backend service
    backendCtx, backendCancel := context.WithCancel(ctx)
    backendErrCh := make(chan error, 1)
    go func() {
        backendErrCh <- startBackendService(backendCtx, cfg)
    }()

    // Start frontend service
    frontendCtx, frontendCancel := context.WithCancel(ctx)
    frontendErrCh := make(chan error, 1)
    go func() {
        frontendErrCh <- startFrontendService(frontendCtx, cfg)
    }()

    // Wait for interrupt or error
    select {
    case <-signalCh:
        log.Println("Received shutdown signal")
    case err := <-backendErrCh:
        log.Printf("Backend service error: %v", err)
    case err := <-frontendErrCh:
        log.Printf("Frontend service error: %v", err)
    }

    // Graceful shutdown
    log.Println("Shutting down services...")
    backendCancel()
    frontendCancel()

    // Wait with timeout for services to shut down
    shutdownTimeout := time.Second * 10
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
    defer shutdownCancel()

    // Wait for both services to shut down or timeout
    select {
    case <-shutdownCtx.Done():
        log.Println("Shutdown timed out, forcing exit")
    case <-waitForShutdown(backendErrCh, frontendErrCh):
        log.Println("All services shut down successfully")
    }
}
```

### Communication Between Services
1. **API Gateway**:
   - Configure reverse proxy to route requests to appropriate services
   - Implement CORS handling for cross-domain requests
   - Add response caching for frequently accessed data

## 5. Development and Deployment Workflow

1. **Development Setup**:
   - Create a docker-compose.yml for local development
   - Implement hot reloading for both backend and frontend
   - Add dev/prod environment configuration

2. **Build Process**:
   - Configure build scripts for both services
   - Implement asset optimization for frontend
   - Create single binary distribution

## 6. Implementation Steps and Timeline

1. **Week 1: Infrastructure Setup**
   - Set up project structure
   - Create API boilerplate
   - Implement authentication system

2. **Week 2: Backend API Development**
   - Migrate existing API endpoints
   - Implement new RESTful API
   - Add API tests

3. **Week 3: WebSocket Implementation**
   - Create WebSocket server
   - Implement message protocol
   - Add real-time update system

4. **Week 4: Frontend Framework**
   - Set up React application
   - Create core components
   - Implement API client services

5. **Week 5: Frontend UI Development**
   - Build dashboard interface
   - Implement wallet monitoring views
   - Create configuration management UI

6. **Week 6: Integration and Testing**
   - Connect frontend to backend API
   - Implement WebSocket client
   - Add end-to-end tests

7. **Week 7: Performance Optimization**
   - Optimize API performance
   - Add caching mechanisms
   - Implement data virtualization

8. **Week 8: Finalization**
   - Complete orchestration logic
   - Add documentation
   - Finalize deployment scripts


**Phase 3:** Add database persistence
Implement database storage for configuration
Add user management and roles
Migrate file-based storage to database

**Phase 4:** Enhance features
Implement additional dashboard controls
Add monitoring statistics and historical data
Create alert management in the UI

## Alerts and Shizz:
- [ ] Test Alerts To Console, Web and Discord
- [ ] Make sure Different Alert Types work
- [ ] Configurable Conditions for Alert Types
- [ ] Send different types to different dicsord channels

## Web Dashboard Stuff:
- [x] Order Tables by Values

- [ ] Show USD Values
        - Think about how to best integrate a price stream (possible through QuickNode RPC Endpoint we arealdy use)

- [ ] Show Token Names
- [ ] Show Wallet Names (Need to setup naming and grouping via config)
        - Provide Search for Names, Addresses, (Token Names?)

- [x] Update Configuration Page to properly show the new settings
        - may need a re-design


## Documentation Stuff:
- [x] Better Align "Community" and "Open Source" sections on the hompage.
        - Especially the "Join Discord" and "GitHub" buttons.

- [x] /discord-integration.html seems to have a missing image.

- [ ] /api/authentication.md has not yet been created.
- [ ] /api/endpoints.md has not yet been created.
- [ ] /api/integration-examples.md has not yet been created.
- [ ] /developer/building.md has not yet been created.
- [ ] /developer/testing.md has not yet been created.

- [x] better seo for the homepage

-> Consider changing the name form "Insider Monitor" to "Solana Monitor"
        - But I don't konw if this is a good idea and if I want to add other chains in the future.
