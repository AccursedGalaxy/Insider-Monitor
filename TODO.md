# TODO

## Current Overview:
- Single process application (cmd/monitor/main.go)
- Web UI and monitor run in the same process
- Configuration is loaded from files only
- No authentication/authorization

## Proposed Architecture:
- Split into separate backend service and frontend application
- RESTful API or gRPC communication between frontend and backend
- Configuration stored in a database with file fallback
- Authentication and authorization for admin actions (e.g. configuration changes)
