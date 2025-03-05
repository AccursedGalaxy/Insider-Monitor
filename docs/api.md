---
layout: default
title: API Reference
description: Complete API reference for Solana Insider Monitor
---

# API Reference

Solana Insider Monitor provides a REST API for programmatic access and integration with other systems. This document outlines all available endpoints, their parameters, and example responses.

## Authentication

Most API endpoints require authentication. The application uses JWT tokens for authentication.

### Obtaining a JWT Token

To authenticate with the API, you need to obtain a JWT token:

**Endpoint:** `POST /api/login`

**Request Body:**
```json
{
  "username": "admin",
  "password": "your_password"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Using the Token

Include the token in the `Authorization` header of your requests:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Public Endpoints

These endpoints do not require authentication.

### Get All Wallets

Retrieves data for all monitored wallets.

**Endpoint:** `GET /api/wallets`

**Response:**
```json
{
  "wallet_address_1": {
    "wallet_address": "wallet_address_1",
    "token_accounts": {
      "token_mint_1": {
        "mint": "token_mint_1",
        "balance": 1000000,
        "decimals": 6,
        "symbol": "SOL",
        "last_updated": "2023-01-01T12:00:00Z"
      }
    },
    "last_scanned": "2023-01-01T12:00:00Z"
  }
}
```

### Get Wallet Details

Retrieves detailed information for a specific wallet.

**Endpoint:** `GET /api/wallets/{address}`

**Response:**
```json
{
  "wallet_address": "wallet_address_1",
  "token_accounts": {
    "token_mint_1": {
      "mint": "token_mint_1",
      "balance": 1000000,
      "decimals": 6,
      "symbol": "SOL",
      "last_updated": "2023-01-01T12:00:00Z"
    }
  },
  "last_scanned": "2023-01-01T12:00:00Z"
}
```

### Refresh Data

Triggers an immediate scan of all wallets.

**Endpoint:** `POST /api/refresh`

**Response:**
```json
{
  "status": "success",
  "message": "Data refreshed"
}
```

## Protected Endpoints

These endpoints require authentication.

### Get Configuration

Retrieves the current application configuration.

**Endpoint:** `GET /api/admin/config`

**Response:**
```json
{
  "NetworkURL": "https://api.mainnet-beta.solana.com",
  "Wallets": [
    "wallet_address_1",
    "wallet_address_2"
  ],
  "ScanInterval": "1m",
  "Alerts": {
    "MinimumBalance": 1000,
    "SignificantChange": 0.2,
    "IgnoreTokens": []
  },
  "Discord": {
    "Enabled": false,
    "WebhookURL": "",
    "ChannelID": ""
  }
}
```

### Update Configuration

Updates the application configuration.

**Endpoint:** `PUT /api/admin/config`

**Request Body:**
```json
{
  "network_url": "https://api.devnet.solana.com",
  "scan_interval": "5m",
  "alerts": {
    "minimum_balance": 500,
    "significant_change": 0.1
  },
  "discord": {
    "enabled": true,
    "webhook_url": "https://discord.com/api/webhooks/..."
  }
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Configuration updated successfully"
}
```

### Add Wallet

Adds a new wallet to monitor.

**Endpoint:** `POST /api/admin/wallets`

**Request Body:**
```json
{
  "address": "new_wallet_address"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Wallet added successfully"
}
```

### Delete Wallet

Removes a wallet from monitoring.

**Endpoint:** `DELETE /api/admin/wallets/{address}`

**Response:**
```json
{
  "status": "success",
  "message": "Wallet removed successfully"
}
```

## Error Responses

The API returns standard HTTP status codes and error messages.

### Common Error Responses

- **400 Bad Request** - Invalid request format or parameters
- **401 Unauthorized** - Authentication failed
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server error

Example error response:
```json
{
  "error": "Invalid wallet address"
}
```

## Rate Limiting

Currently, there is no rate limiting on the API endpoints. However, be aware that the application limits requests to the Solana RPC endpoint to prevent overloading.
