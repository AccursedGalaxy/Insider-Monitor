# API Reference

Solana Insider Monitor provides a comprehensive REST API that allows you to programmatically access wallet data, manage configuration, and integrate with other systems.

## API Overview

The API is available when running Solana Insider Monitor with the `-web` flag. All API endpoints are served from the base URL of your instance, for example: `http://localhost:8080/api/`.

## Authentication

Most API endpoints require authentication using JSON Web Tokens (JWT). To obtain a token, you must first authenticate using your credentials.

### Obtaining a JWT Token

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'
```

This will return a JWT token that you should include in the `Authorization` header for subsequent requests:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2023-05-15T12:00:00Z"
}
```

### Using the JWT Token

Include the token in the `Authorization` header for all protected API requests:

```bash
curl -X GET http://localhost:8080/api/admin/config \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## API Endpoints

The API is divided into public and protected endpoints:

### Public Endpoints

These endpoints are accessible without authentication:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Check if the API is running |
| `/api/login` | POST | Authenticate and get a JWT token |
| `/api/wallets` | GET | Get a list of all monitored wallets and their current state |
| `/api/wallets/{address}` | GET | Get detailed information about a specific wallet |

### Protected Endpoints

These endpoints require authentication:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/admin/config` | GET | Get the current configuration |
| `/api/admin/config` | PUT | Update the configuration |
| `/api/admin/wallets` | POST | Add a new wallet to monitor |
| `/api/admin/wallets/{address}` | DELETE | Remove a wallet from monitoring |
| `/api/refresh` | POST | Trigger an immediate refresh of wallet data |

## Response Format

All API responses are returned in JSON format with a consistent structure:

```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

Or in case of an error:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required"
  }
}
```

## Rate Limiting

To prevent abuse, the API implements rate limiting:

- 60 requests per minute for authenticated endpoints
- 20 requests per minute for unauthenticated endpoints

Rate limit headers are included in all responses:

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 58
X-RateLimit-Reset: 1620000000
```

## Detailed API Documentation

For detailed information about specific endpoints, see:

- [Authentication](authentication.md) - Detailed guide on authentication
- [Endpoints](endpoints.md) - Complete endpoint reference
- [Integration Examples](integration-examples.md) - Code examples for common use cases

## Common Use Cases

### Monitoring Multiple Wallets

You can use the API to fetch data for all monitored wallets:

```bash
curl -X GET http://localhost:8080/api/wallets
```

### Adding a New Wallet to Monitor

```bash
curl -X POST http://localhost:8080/api/admin/wallets \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"wallet_address": "NEW_WALLET_ADDRESS"}'
```

### Updating Alert Configuration

```bash
curl -X PUT http://localhost:8080/api/admin/config \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "alerts": {
      "minimum_balance": 500,
      "significant_change": 0.10
    }
  }'
```

## Webhook Integration

You can configure external systems to receive webhook notifications for alerts:

1. Set up an endpoint in your application to receive POST requests
2. Configure the webhook URL in Solana Insider Monitor
3. Receive real-time alerts when significant changes are detected

```json
{
  "webhook_url": "https://your-service.com/webhook",
  "webhook_secret": "your-webhook-secret"
}
```

## API Client Libraries

We provide official client libraries for several programming languages:

- [JavaScript/TypeScript](https://github.com/accursedgalaxy/insider-monitor-js)
- [Python](https://github.com/accursedgalaxy/insider-monitor-py)
- [Go](https://github.com/accursedgalaxy/insider-monitor-go)

## API Versioning

The API follows semantic versioning. The current version is v1.

To ensure backward compatibility, we include the version in the URL:

```
/api/v1/wallets
```

The default `/api/` path points to the latest stable version.

## Error Codes

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Authentication required or invalid token |
| `FORBIDDEN` | 403 | Insufficient permissions for this operation |
| `NOT_FOUND` | 404 | Requested resource not found |
| `VALIDATION_ERROR` | 400 | Invalid request parameters |
| `RATE_LIMITED` | 429 | Too many requests, try again later |
| `SERVER_ERROR` | 500 | Internal server error |

## Next Steps

- [Authentication](authentication.md) - Learn more about authentication
- [Endpoints](endpoints.md) - Explore all available endpoints
- [Integration Examples](integration-examples.md) - See code samples for common use cases
