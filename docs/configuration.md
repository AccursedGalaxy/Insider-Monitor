---
layout: default
title: Configuration Guide
description: Detailed configuration guide for Solana Insider Monitor
---

# Configuration Guide

This guide explains all configuration options for Solana Insider Monitor and how to customize them for your needs.

## Configuration File

The application uses a JSON configuration file (`config.json`) to store settings. By default, it looks for this file in the current working directory, but you can specify a different location using the `-config` flag.

### Basic Structure

```json
{
    "network_url": "https://api.mainnet-beta.solana.com",
    "wallets": [
        "WALLET_ADDRESS_1",
        "WALLET_ADDRESS_2"
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

## Configuration Options

### Network Settings

#### `network_url`

The Solana RPC endpoint URL to connect to. This can be a public or private endpoint.

**Options:**
- Mainnet: `https://api.mainnet-beta.solana.com`
- Devnet: `https://api.devnet.solana.com`
- Testnet: `https://api.testnet.solana.com`
- Custom RPC endpoints (e.g., from QuickNode, Alchemy, etc.)

**Example:**
```json
"network_url": "https://api.mainnet-beta.solana.com"
```

### Wallet Settings

#### `wallets`

An array of Solana wallet addresses to monitor. These should be base58-encoded public keys.

**Example:**
```json
"wallets": [
    "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr",
    "DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF"
]
```

### Monitoring Settings

#### `scan_interval`

The interval between wallet scans, specified as a duration string.

**Format:** A string in Go's time.Duration format (e.g., "5s", "1m", "1h")

**Common values:**
- "30s" - 30 seconds
- "1m" - 1 minute
- "5m" - 5 minutes
- "1h" - 1 hour

**Example:**
```json
"scan_interval": "1m"
```

### Alert Settings

The `alerts` section configures when and how alerts are triggered.

#### `minimum_balance`

The minimum token balance required to trigger alerts. Tokens with balances below this threshold will not generate alerts.

**Example:**
```json
"minimum_balance": 1000
```

#### `significant_change`

The percentage change in token balance required to trigger an alert, expressed as a decimal.
- 0.05 = 5%
- 0.20 = 20%
- 1.00 = 100%

**Example:**
```json
"significant_change": 0.20
```

#### `ignore_tokens`

An array of token mint addresses to ignore when monitoring. Useful for excluding tokens you're not interested in.

**Example:**
```json
"ignore_tokens": [
    "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
]
```

### Discord Notification Settings

The `discord` section configures Discord webhook notifications.

#### `enabled`

Whether Discord notifications are enabled.

**Example:**
```json
"enabled": true
```

#### `webhook_url`

The Discord webhook URL to send notifications to. You can create a webhook in your Discord server settings.

**Example:**
```json
"webhook_url": "https://discord.com/api/webhooks/123456789/abcdefghijklmnopqrstuvwxyz"
```

#### `channel_id`

The Discord channel ID to send notifications to. This is optional but can be used to override the channel specified in the webhook.

**Example:**
```json
"channel_id": "123456789012345678"
```

## Environment Variables

The application also supports configuration through environment variables:

| Environment Variable | Description |
|---------------------|-------------|
| `ADMIN_PASSWORD` | Sets the admin password (default: "admin") |
| `JWT_SECRET` | Sets the JWT secret key for authentication |

**Example:**
```bash
export ADMIN_PASSWORD="secure_password"
export JWT_SECRET="custom_secret_key"
```

## Configuration Management

### Via the Web Interface

You can manage the configuration through the web interface at `/config` when running with the `-web` flag.

### Via the API

You can also manage the configuration programmatically through the API.

**Get Configuration:**
```bash
curl -X GET http://localhost:8080/api/admin/config \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Update Configuration:**
```bash
curl -X PUT http://localhost:8080/api/admin/config \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "network_url": "https://api.devnet.solana.com",
    "scan_interval": "5m"
  }'
```

## Configuration Examples

### Minimal Configuration

```json
{
    "network_url": "https://api.mainnet-beta.solana.com",
    "wallets": ["YOUR_WALLET_ADDRESS"],
    "scan_interval": "1m"
}
```

### Production Configuration with Discord Alerts

```json
{
    "network_url": "https://api.mainnet-beta.solana.com",
    "wallets": [
        "WALLET_ADDRESS_1",
        "WALLET_ADDRESS_2"
    ],
    "scan_interval": "5m",
    "alerts": {
        "minimum_balance": 1000,
        "significant_change": 0.10,
        "ignore_tokens": ["TOKEN_MINT_1", "TOKEN_MINT_2"]
    },
    "discord": {
        "enabled": true,
        "webhook_url": "https://discord.com/api/webhooks/...",
        "channel_id": "123456789012345678"
    }
}
```

### Development Configuration

```json
{
    "network_url": "https://api.devnet.solana.com",
    "wallets": [
        "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr"
    ],
    "scan_interval": "30s",
    "alerts": {
        "minimum_balance": 100,
        "significant_change": 0.05
    },
    "discord": {
        "enabled": false
    }
}
```
