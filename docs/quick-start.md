---
layout: default
title: Quick Start Guide
description: Get started quickly with Solana Insider Monitor
---

# Quick Start Guide

This guide will help you get up and running with Solana Insider Monitor in just a few minutes.

## Prerequisites

- **Go** - version 1.23.2 or later
- **Access to a Solana RPC endpoint** - either mainnet, devnet, or testnet

## 1. Clone and Setup

```bash
# Clone the repository
git clone https://github.com/accursedgalaxy/insider-monitor
cd insider-monitor

# Install dependencies
go mod download
```

## 2. Configure

Create a minimal configuration file:

```bash
cat > config.json << EOF
{
    "network_url": "https://api.devnet.solana.com",
    "wallets": [
        "YOUR_WALLET_ADDRESS"
    ],
    "scan_interval": "1m"
}
EOF
```

> **Note:** Replace `YOUR_WALLET_ADDRESS` with an actual Solana wallet address you want to monitor.

## 3. Run in Test Mode

For a quick demonstration with mock data:

```bash
go run cmd/monitor/main.go -test
```

This will run the application with simulated wallet data, perfect for testing and demonstrations.

## 4. Run with Web Interface

To monitor actual wallets with a web interface:

```bash
go run cmd/monitor/main.go -web
```

Then open your browser and navigate to: [http://localhost:8080](http://localhost:8080)

### Default Login Credentials

- **Username:** admin
- **Password:** admin

## 5. Use the API

Once running, you can access the API endpoints:

```bash
# Get a JWT token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Use the token to access protected endpoints
curl -X GET http://localhost:8080/api/admin/config \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Next Steps

- [Configuration Guide](./configuration.md) - Learn about all configuration options
- [Web Interface Guide](./web-interface.md) - Learn how to use the web interface
- [API Reference](./api.md) - Explore the API for programmatic access
- [Authentication](./authentication.md) - Learn about authentication and security
