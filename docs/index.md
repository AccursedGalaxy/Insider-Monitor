---
layout: home
title: Home
nav_order: 1
permalink: /
---

# Solana Insider Monitor

{: .fs-9 }
A tool for monitoring Solana wallet activities and balance changes
{: .fs-6 .fw-300 }

[Get Started](./quick-start.html){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View on GitHub](https://github.com/accursedgalaxy/Insider-Monitor){: .btn .fs-5 .mb-4 .mb-md-0 }

---

## Features

- **Monitor multiple Solana wallets** simultaneously
- **Track token balance changes** in real-time
- **Real-time alerts** for significant changes
- **Discord integration** for instant notifications
- **Persistent storage** of wallet data
- **Web interface** for monitoring and configuration
- **REST API** for programmatic access and integration
{: .mb-6 }

## Getting Started

### Prerequisites

- Go 1.23.2 or later
- Access to a Solana RPC endpoint (mainnet, devnet, or testnet)

### Installation

```bash
# Clone the repository
git clone https://github.com/accursedgalaxy/insider-monitor
cd insider-monitor

# Install dependencies
go mod download
```

For detailed installation instructions, see the [Installation Guide](./installation.html).

---

## Documentation

### Guides
- [Installation Guide](./installation.html)
- [Quick Start Guide](./quick-start.html)
- [Configuration Guide](./configuration.html)

### Reference
- [API Reference](./api.html)
- [Authentication](./authentication.html)
- [Web Interface](./web-interface.html)

## Current Features and Roadmap

Solana Insider Monitor is under active development. Here's the current status:

### Phase 1: API and Configuration Management (Complete)
- ✅ REST API endpoints
- ✅ Configuration management via API
- ✅ Authentication for sensitive operations
- ✅ Web interface for monitoring

### Phase 2: Refactor Backend/Frontend (Planned)
- 🔄 Separate backend and frontend services
- 🔄 Enhanced API communication
- 🔄 WebSockets for real-time updates

### Phase 3: Database Persistence (Planned)
- 🔄 Database storage for configuration
- 🔄 User management and roles
- 🔄 Migration from file-based storage

### Phase 4: Enhanced Features (Planned)
- 🔄 Advanced dashboard controls
- 🔄 Monitoring statistics and historical data
- 🔄 Alert management in UI

## Community

Join our [Discord community](https://discord.gg/7vY9ZBPdya) to:
- Get help with setup and configuration
- Share feedback and suggestions
- Connect with other users
- Stay updated on new features and releases
- Discuss Solana development
