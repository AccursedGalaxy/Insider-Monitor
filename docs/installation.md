---
layout: default
title: Installation Guide
description: How to install and set up Solana Insider Monitor
---

# Installation Guide

This guide will walk you through the process of installing and setting up Solana Insider Monitor on your system.

## Prerequisites

Before you begin, ensure you have the following prerequisites:

- **Go** - version 1.23.2 or later
- **Git** - for cloning the repository
- **Access to a Solana RPC endpoint** - either mainnet, devnet, or testnet

## Step 1: Clone the Repository

```bash
# Clone the repository
git clone https://github.com/accursedgalaxy/insider-monitor
cd insider-monitor
```

## Step 2: Install Dependencies

```bash
# Install Go dependencies
go mod download
```

## Step 3: Configure the Application

Create a copy of the example configuration file:

```bash
cp config.example.json config.json
```

Edit the `config.json` file to include your Solana RPC endpoint and wallet addresses:

```json
{
    "network_url": "https://api.mainnet-beta.solana.com",
    "wallets": [
        "YOUR_WALLET_ADDRESS_1",
        "YOUR_WALLET_ADDRESS_2"
    ],
    "scan_interval": "1m",
    "alerts": {
        "minimum_balance": 1000,
        "significant_change": 0.20,
        "ignore_tokens": []
    },
    "discord": {
        "enabled": false,
        "webhook_url": "",
        "channel_id": ""
    }
}
```

For more detailed configuration options, see the [Configuration Guide](./configuration.md).

## Step 4: Build the Application

You can either run the application directly with Go or build a binary:

```bash
# Option 1: Run directly
go run cmd/monitor/main.go

# Option 2: Build a binary
make build
```

If you use the `make build` option, the binary will be available in the `bin` directory.

## Step 5: Run the Application

### Console Mode

```bash
go run cmd/monitor/main.go
```

### Web Interface Mode

```bash
go run cmd/monitor/main.go -web
```

### Test Mode (with mock data)

```bash
go run cmd/monitor/main.go -test
```

## Next Steps

- [Configuration Guide](./configuration.md) - Learn about all configuration options
- [Web Interface Guide](./web-interface.md) - Learn how to use the web interface
- [API Reference](./api.md) - Explore the API for programmatic access
