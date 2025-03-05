---
layout: default
title: Configuration Guide
nav_order: 4
description: Configuration options and examples for Solana Insider Monitor
---

# Configuration Guide
{: .no_toc }

This guide explains all configuration options available in Solana Insider Monitor.

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

## Configuration File

Solana Insider Monitor uses a JSON configuration file to store settings. By default, the application looks for a `config.json` file in the current directory, but you can specify a different path using the `-config` flag.

### Example Configuration

```json
{
    "network_url": "https://api.mainnet-beta.solana.com",
    "wallets": [
        "5xLEuN615KQTuGZMqFfiNsR6SMzSG6Sd9PkXHPupYXQL",
        "AaXs7cLGcSVAsEt8QxstVrqhLhYN2iGGEwRLMmrAjHkN"
    ],
    "scan_interval": "1m",
    "alerts": {
        "minimum_balance": 1000,
        "significant_change": 0.20,
        "ignore_tokens": []
    },
    "discord": {
        "enabled": true,
        "webhook_url": "https://discord.com/api/webhooks/your-webhook-url",
        "channel_id": "your-channel-id"
    },
    "web": {
        "port": 8080,
        "enable_cors": true,
        "allowed_origins": ["http://localhost:3000"]
    }
}
```

## Network Configuration

### network_url

The Solana RPC endpoint URL to use for querying the blockchain.

**Type:** String
**Required:** Yes
**Default:** None

**Common options:**
- Mainnet: `https://api.mainnet-beta.solana.com`
- Devnet: `https://api.devnet.solana.com`
- Testnet: `https://api.testnet.solana.com`

**Example:**
```json
"network_url": "https://api.mainnet-beta.solana.com"
```

### use_rate_limit

Whether to apply rate limiting to RPC requests to avoid hitting rate limits.

**Type:** Boolean
**Required:** No
**Default:** true

**Example:**
```json
"use_rate_limit": true
```

## Wallet Configuration

### wallets

An array of Solana wallet addresses to monitor.

**Type:** Array of Strings
**Required:** Yes
**Default:** None

**Example:**
```json
"wallets": [
    "5xLEuN615KQTuGZMqFfiNsR6SMzSG6Sd9PkXHPupYXQL",
    "AaXs7cLGcSVAsEt8QxstVrqhLhYN2iGGEwRLMmrAjHkN"
]
```

### scan_interval

The time interval between scanning wallets for changes.

**Type:** String (Duration format)
**Required:** No
**Default:** "1m" (1 minute)

**Format examples:**
- "30s" (30 seconds)
- "1m" (1 minute)
- "5m" (5 minutes)
- "1h" (1 hour)

**Example:**
```json
"scan_interval": "1m"
```

## Alert Configuration

### alerts.minimum_balance

The minimum token balance required to trigger alerts. Balances below this threshold will not generate alerts.

**Type:** Number
**Required:** No
**Default:** 1000

**Example:**
```json
"minimum_balance": 1000
```

### alerts.significant_change

The percentage change required to trigger alerts. Changes below this threshold will not generate alerts.

**Type:** Number (0.0 to 1.0, representing 0% to 100%)
**Required:** No
**Default:** 0.20 (20%)

**Example:**
```json
"significant_change": 0.20
```

### alerts.ignore_tokens

An array of token addresses to ignore when monitoring for changes.

**Type:** Array of Strings
**Required:** No
**Default:** [] (empty array)

**Example:**
```json
"ignore_tokens": [
    "So11111111111111111111111111111111111111112",
    "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
]
```

## Discord Integration

### discord.enabled

Whether to enable Discord notifications.

**Type:** Boolean
**Required:** No
**Default:** false

**Example:**
```json
"enabled": true
```

### discord.webhook_url

The Discord webhook URL to use for sending notifications.

**Type:** String
**Required:** If discord.enabled is true
**Default:** None

**Example:**
```json
"webhook_url": "https://discord.com/api/webhooks/your-webhook-url"
```

### discord.channel_id

The Discord channel ID to send notifications to.

**Type:** String
**Required:** If discord.enabled is true
**Default:** None

**Example:**
```json
"channel_id": "your-channel-id"
```

## Web Interface Configuration

### web.port

The port to use for the web interface.

**Type:** Number
**Required:** No
**Default:** 8080

**Example:**
```json
"port": 8080
```

### web.enable_cors

Whether to enable Cross-Origin Resource Sharing (CORS) for the web interface.

**Type:** Boolean
**Required:** No
**Default:** false

**Example:**
```json
"enable_cors": true
```

### web.allowed_origins

An array of allowed origins for CORS requests.

**Type:** Array of Strings
**Required:** If web.enable_cors is true
**Default:** [] (empty array)

**Example:**
```json
"allowed_origins": ["http://localhost:3000", "https://yourdomain.com"]
```

## Data Storage Configuration

### data.directory

The directory to use for storing data.

**Type:** String
**Required:** No
**Default:** "./data"

**Example:**
```json
"data.directory": "./storage/wallet-data"
```

## Command Line Flags

In addition to the configuration file, Solana Insider Monitor also supports several command line flags:

| Flag | Description | Default |
|------|-------------|---------|
| `-config` | Path to the configuration file | `./config.json` |
| `-web` | Enable the web interface | `false` |
| `-port` | Port to use for the web interface | `8080` |
| `-test` | Run in test mode with mock data | `false` |
| `-debug` | Enable debug logging | `false` |

**Example:**
```bash
go run cmd/monitor/main.go -web -port 9090 -config custom-config.json
```

## Environment Variables

Solana Insider Monitor also supports configuration via environment variables, which take precedence over the configuration file:

| Variable | Description | Example |
|----------|-------------|---------|
| `NETWORK_URL` | Solana RPC endpoint URL | `https://api.mainnet-beta.solana.com` |
| `SCAN_INTERVAL` | Time between scans | `1m` |
| `MINIMUM_BALANCE` | Minimum balance for alerts | `1000` |
| `SIGNIFICANT_CHANGE` | Percentage for significant change alerts | `0.20` |
| `DISCORD_WEBHOOK_URL` | Discord webhook URL | `https://discord.com/api/webhooks/...` |
| `DISCORD_CHANNEL_ID` | Discord channel ID | `123456789012345678` |
| `ADMIN_PASSWORD` | Admin password for web interface | `your-secure-password` |
| `WEB_PORT` | Port for web interface | `8080` |
| `DATA_DIRECTORY` | Directory for data storage | `./data` |

**Example:**
```bash
export NETWORK_URL="https://api.mainnet-beta.solana.com"
export SCAN_INTERVAL="1m"
export ADMIN_PASSWORD="secure-password"
go run cmd/monitor/main.go -web
```

## Managing Configuration via API

Solana Insider Monitor provides API endpoints for managing configuration at runtime:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/admin/config` | GET | Get current configuration |
| `/api/admin/config` | PUT | Update configuration |
| `/api/admin/wallets` | POST | Add a new wallet to monitor |
| `/api/admin/wallets/{address}` | DELETE | Remove a wallet from monitoring |

For more information on using the API, see the [API Reference](./api.html).

## Configuration Best Practices

1. **Security**:
   - Store sensitive configuration (like webhook URLs) as environment variables rather than in the configuration file
   - Use a strong admin password for the web interface
   - Enable CORS only if necessary, and restrict to specific origins

2. **Performance**:
   - Set a reasonable scan interval based on your needs (1-5 minutes is usually sufficient)
   - Enable rate limiting to avoid hitting RPC endpoint limits
   - Monitor CPU and memory usage if tracking many wallets

3. **Storage**:
   - Regularly backup your data directory
   - Consider mounting a separate volume for data storage in production environments
