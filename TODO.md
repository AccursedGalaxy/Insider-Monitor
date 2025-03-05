# TODO

## Current Overview:
- Single process application (cmd/monitor/main.go)
- Web UI and monitor run in the same process
- Configuration is loaded from files only
- No authentication/authorization


## Upgrade Plan:
Phase 1: Create API endpoints in the current application
Add basic CRUD endpoints for configuration
Implement configuration updating while running
Add authentication for sensitive operations

Phase 2: Refactor to separate backend and frontend
Split codebase into backend and frontend services
Implement proper API communication
Set up WebSockets for real-time updates

Phase 3: Add database persistence
Implement database storage for configuration
Add user management and roles
Migrate file-based storage to database

Phase 4: Enhance features
Implement additional dashboard controls
Add monitoring statistics and historical data
Create alert management in the UI
