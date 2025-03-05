---
layout: default
title: API Reference
nav_order: 5
description: Comprehensive API documentation for Solana Insider Monitor
---

# API Reference
{: .no_toc }

This document provides a comprehensive reference for the Solana Insider Monitor REST API.

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

## Overview

Solana Insider Monitor provides a REST API that allows you to:

- Retrieve wallet data and balance information
- Configure monitoring settings
- Manage wallet addresses
- Refresh data
- Authenticate and secure API access

The API is organized around standard REST principles and uses JSON for request and response bodies.

## Base URL

When running the application with the `-web` flag, the API is available at:

```
http://localhost:8080/api
```

You can specify a different port using the `-port` flag:

```
http://localhost:8080/api  # Default
http://localhost:9090/api  # With -port 9090
```

## Authentication

Most API endpoints require authentication. The API uses JWT (JSON Web Tokens) for authentication.

### Obtaining a Token

To obtain a JWT token, you need to authenticate using the login endpoint:

**Endpoint:**
```
POST /api/login
```

**Request:**
```json
{
  "username": "admin",
  "password": "admin"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2023-05-01T12:00:00Z"
}
```

### Using the Token

Include the token in the `Authorization` header for authenticated requests:

```
Authorization: Bearer YOUR_TOKEN
```

### Token Expiration

Tokens are valid for 24 hours by default. After expiration, you need to obtain a new token.

## Public Endpoints

These endpoints are accessible without authentication.

### Get All Wallets

Retrieves data for all monitored wallets.

**Endpoint:**
```
GET /api/wallets
```

**Response:**
```json
{
  "wallets": [
    {
      "address": "5xLEuN615KQTuGZMqFfiNsR6SMzSG6Sd9PkXHPupYXQL",
      "last_updated": "2023-04-30T15:30:00Z",
      "tokens": [
        {
          "mint": "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
          "symbol": "USDC",
          "balance": 1000000,
          "decimals": 6,
          "usd_value": 1000
        },
        {
          "mint": "So11111111111111111111111111111111111111112",
          "symbol": "SOL",
          "balance": 5000000000,
          "decimals": 9,
          "usd_value": 500
        }
      ]
    }
  ]
}
```

### Get Specific Wallet

Retrieves data for a specific wallet.

**Endpoint:**
```
GET /api/wallets/{address}
```

**Parameters:**
- `address`: The wallet address to retrieve data for

**Response:**
```json
{
  "address": "5xLEuN615KQTuGZMqFfiNsR6SMzSG6Sd9PkXHPupYXQL",
  "last_updated": "2023-04-30T15:30:00Z",
  "tokens": [
    {
      "mint": "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
      "symbol": "USDC",
      "balance": 1000000,
      "decimals": 6,
      "usd_value": 1000
    },
    {
      "mint": "So11111111111111111111111111111111111111112",
      "symbol": "SOL",
      "balance": 5000000000,
      "decimals": 9,
      "usd_value": 500
    }
  ]
}
```

### Refresh Data

Triggers an immediate refresh of wallet data.

**Endpoint:**
```
POST /api/refresh
```

**Response:**
```json
{
  "success": true,
  "message": "Refresh initiated"
}
```

## Protected Endpoints

These endpoints require authentication using a JWT token.

### Get Configuration

Retrieves the current configuration.

**Endpoint:**
```
GET /api/admin/config
```

**Response:**
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
    "webhook_url": "https://discord.com/api/webhooks/...",
    "channel_id": "123456789012345678"
  }
}
```

### Update Configuration

Updates the configuration. You can provide partial updates - only the specified fields will be updated.

**Endpoint:**
```
PUT /api/admin/config
```

**Request:**
```json
{
  "network_url": "https://api.devnet.solana.com",
  "scan_interval": "5m",
  "alerts": {
    "minimum_balance": 500
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "Configuration updated"
}
```

### Add Wallet

Adds a new wallet address to monitor.

**Endpoint:**
```
POST /api/admin/wallets
```

**Request:**
```json
{
  "address": "DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Wallet added"
}
```

### Remove Wallet

Removes a wallet address from monitoring.

**Endpoint:**
```
DELETE /api/admin/wallets/{address}
```

**Parameters:**
- `address`: The wallet address to remove

**Response:**
```json
{
  "success": true,
  "message": "Wallet removed"
}
```

## Error Handling

The API uses standard HTTP status codes to indicate success or failure:

- `200 OK`: The request was successful
- `400 Bad Request`: The request was invalid
- `401 Unauthorized`: The request requires authentication
- `403 Forbidden`: The authenticated user does not have permission
- `404 Not Found`: The requested resource was not found
- `500 Internal Server Error`: An error occurred on the server

Error responses include a JSON body with details:

```json
{
  "error": true,
  "message": "Detailed error message",
  "code": "ERROR_CODE"
}
```

## Rate Limiting

To prevent abuse, the API implements rate limiting. If you exceed the rate limit, you'll receive a `429 Too Many Requests` response.

## Webhooks

In addition to the REST API, Solana Insider Monitor can send webhook notifications when wallet balances change significantly.

### Discord Webhooks

To set up Discord notifications, configure the discord section in your configuration:

```json
"discord": {
  "enabled": true,
  "webhook_url": "https://discord.com/api/webhooks/...",
  "channel_id": "123456789012345678"
}
```

For more information on setting up Discord webhooks, see the [Configuration Guide](./configuration.html).

## API Client Examples

### cURL

**Authentication:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'
```

**Get Wallets:**
```bash
curl -X GET http://localhost:8080/api/wallets
```

**Update Configuration:**
```bash
curl -X PUT http://localhost:8080/api/admin/config \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "scan_interval": "5m"
  }'
```

### JavaScript (Node.js)

```javascript
const axios = require('axios');

// Authentication
const authenticate = async () => {
  const response = await axios.post('http://localhost:8080/api/login', {
    username: 'admin',
    password: 'admin'
  });

  return response.data.token;
};

// Get Wallets
const getWallets = async () => {
  const response = await axios.get('http://localhost:8080/api/wallets');
  return response.data;
};

// Update Configuration
const updateConfig = async (token) => {
  const response = await axios.put(
    'http://localhost:8080/api/admin/config',
    {
      scan_interval: '5m'
    },
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );

  return response.data;
};

// Example usage
(async () => {
  try {
    const token = await authenticate();
    const wallets = await getWallets();
    const updateResult = await updateConfig(token);

    console.log('Wallets:', wallets);
    console.log('Update Result:', updateResult);
  } catch (error) {
    console.error('Error:', error.response?.data || error.message);
  }
})();
```

### Python

```python
import requests

# Authentication
def authenticate():
    response = requests.post(
        'http://localhost:8080/api/login',
        json={'username': 'admin', 'password': 'admin'}
    )
    return response.json()['token']

# Get Wallets
def get_wallets():
    response = requests.get('http://localhost:8080/api/wallets')
    return response.json()

# Update Configuration
def update_config(token):
    headers = {'Authorization': f'Bearer {token}'}
    response = requests.put(
        'http://localhost:8080/api/admin/config',
        json={'scan_interval': '5m'},
        headers=headers
    )
    return response.json()

# Example usage
try:
    token = authenticate()
    wallets = get_wallets()
    update_result = update_config(token)

    print('Wallets:', wallets)
    print('Update Result:', update_result)
except Exception as e:
    print('Error:', str(e))
```

## API Versioning

The current API version is v1, which is implicit in the API paths. Future API versions may be introduced with explicit version prefixes (e.g., `/api/v2/...`).

## Further Resources

- [Configuration Guide](./configuration.html)
- [Authentication Guide](./authentication.html)
- [Web Interface Guide](./web-interface.html)
