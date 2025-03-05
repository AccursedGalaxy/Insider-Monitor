---
layout: home
title: Home
nav_order: 1
description: "Solana Insider Monitor - A tool for monitoring Solana wallet activities and balance changes"
permalink: /
---

# Solana Insider Monitor
{: .fs-9 }

A powerful tool for monitoring Solana wallet activities, detecting balance changes, and receiving real-time alerts.
{: .fs-6 .fw-300 }

[Get Started](./quick-start.html){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View on GitHub](https://github.com/accursedgalaxy/Insider-Monitor){: .btn .fs-5 .mb-4 .mb-md-0 }

---

## Features

- Monitor multiple Solana wallets simultaneously
- Track token balance changes
- Real-time alerts for significant changes
- Discord integration for notifications
- Persistent storage of wallet data
- Web interface for monitoring and configuration
- REST API for programmatic access and integration
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

## Guides

<div class="card-grid">
  <div class="card">
    <h3>Installation Guide</h3>
    <p>Step-by-step instructions for installing and setting up Solana Insider Monitor.</p>
    <div class="card-footer">
      <a href="./installation.html" class="btn btn-primary">Read More</a>
    </div>
  </div>

  <div class="card">
    <h3>Quick Start Guide</h3>
    <p>Get up and running quickly with a basic configuration.</p>
    <div class="card-footer">
      <a href="./quick-start.html" class="btn btn-primary">Read More</a>
    </div>
  </div>

  <div class="card">
    <h3>Configuration Guide</h3>
    <p>Learn about all configuration options and how to customize your setup.</p>
    <div class="card-footer">
      <a href="./configuration.html" class="btn btn-primary">Read More</a>
    </div>
  </div>
</div>

## Documentation

<div class="card-grid">
  <div class="card">
    <h3>API Reference</h3>
    <p>Comprehensive API documentation for developers.</p>
    <div class="card-footer">
      <a href="./api.html" class="btn btn-primary">View Reference</a>
    </div>
  </div>

  <div class="card">
    <h3>Authentication</h3>
    <p>Learn about authentication and security features.</p>
    <div class="card-footer">
      <a href="./authentication.html" class="btn btn-primary">View Guide</a>
    </div>
  </div>

  <div class="card">
    <h3>Web Interface</h3>
    <p>Learn how to use the web interface to monitor wallets and manage settings.</p>
    <div class="card-footer">
      <a href="./web-interface.html" class="btn btn-primary">View Guide</a>
    </div>
  </div>
</div>

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
