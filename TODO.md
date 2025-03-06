# TODO

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
