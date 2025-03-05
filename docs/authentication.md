---
layout: default
title: Authentication
nav_order: 6
description: Authentication and security features for Solana Insider Monitor
---

# Authentication Guide
{: .no_toc }

This guide explains the authentication and security features in Solana Insider Monitor.

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

## Overview

Solana Insider Monitor uses JWT (JSON Web Tokens) for authentication. This secure method allows the API to verify the identity of users without storing session information on the server.

## Authentication Flow

1. **Login**: The user provides credentials
2. **Token Generation**: The server validates credentials and generates a JWT token
3. **Token Usage**: The client includes the token in subsequent requests
4. **Verification**: The server verifies the token for each protected request

## Login

To authenticate with the API, you need to send a POST request to the login endpoint with your credentials:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'
```

### Request

**Endpoint:** `POST /api/login`

**Headers:**
- `Content-Type: application/json`

**Body:**
```json
{
  "username": "admin",
  "password": "admin"
}
```

### Response

On successful authentication, the server returns a JWT token:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2023-05-01T12:00:00Z"
}
```

## Using the Token

Once you have obtained a token, include it in the `Authorization` header of your requests to protected endpoints:

```bash
curl -X GET http://localhost:8080/api/admin/config \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Headers

For all protected endpoints, include the following header:

```
Authorization: Bearer YOUR_TOKEN
```

Replace `YOUR_TOKEN` with the actual token received from the login endpoint.

## Token Expiration

JWT tokens have an expiration time (typically 24 hours). After expiration, you need to request a new token by authenticating again.

## Protected Endpoints

The following endpoints require authentication:

- `GET /api/admin/config` - Get configuration
- `PUT /api/admin/config` - Update configuration
- `POST /api/admin/wallets` - Add wallet
- `DELETE /api/admin/wallets/{address}` - Remove wallet

## Public Endpoints

These endpoints do not require authentication:

- `POST /api/login` - Authentication
- `GET /api/wallets` - Get all wallets
- `GET /api/wallets/{address}` - Get specific wallet
- `POST /api/refresh` - Refresh data

## Default Credentials

The default credentials for the web interface and API are:

- **Username:** `admin`
- **Password:** `admin`

For security reasons, it's strongly recommended to change the default password in production environments.

## Changing the Admin Password

You can change the admin password by setting the `ADMIN_PASSWORD` environment variable before starting the application:

```bash
export ADMIN_PASSWORD="your-secure-password"
go run cmd/monitor/main.go -web
```

## Security Best Practices

### 1. Change Default Credentials

Always change the default admin password in production environments.

### 2. Use HTTPS

When deploying in production, use HTTPS to encrypt the communication between clients and the server, protecting tokens and credentials in transit.

### 3. Store Tokens Securely

Client applications should store JWT tokens securely:

- Web applications: Use HttpOnly cookies or secure local storage
- Mobile applications: Use secure storage options like Keychain (iOS) or KeyStore (Android)
- Scripts: Store tokens in environment variables instead of hardcoding them

### 4. Token Refresh Strategy

If your application uses the API for extended periods, implement a token refresh strategy to handle token expiration gracefully.

### 5. Firewall Rules

Configure firewall rules to restrict access to the admin API endpoints, allowing connections only from trusted IP addresses.

## Authentication Errors

Common authentication errors include:

### 1. Invalid Credentials

```json
{
  "error": true,
  "message": "Invalid username or password",
  "code": "AUTH_INVALID_CREDENTIALS"
}
```

### 2. Missing Authorization Header

```json
{
  "error": true,
  "message": "Authorization header required",
  "code": "AUTH_HEADER_MISSING"
}
```

### 3. Invalid Token Format

```json
{
  "error": true,
  "message": "Invalid token format",
  "code": "AUTH_TOKEN_INVALID"
}
```

### 4. Expired Token

```json
{
  "error": true,
  "message": "Token has expired",
  "code": "AUTH_TOKEN_EXPIRED"
}
```

## Authentication in Code Examples

### JavaScript (Browser)

```javascript
// Login function
async function login(username, password) {
  const response = await fetch('http://localhost:8080/api/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ username, password })
  });

  const data = await response.json();

  if (response.ok) {
    // Store token securely
    localStorage.setItem('auth_token', data.token);
    return data.token;
  } else {
    throw new Error(data.message || 'Authentication failed');
  }
}

// Function to make authenticated requests
async function authenticatedRequest(url, method = 'GET', body = null) {
  const token = localStorage.getItem('auth_token');

  if (!token) {
    throw new Error('Not authenticated');
  }

  const options = {
    method,
    headers: {
      'Authorization': `Bearer ${token}`
    }
  };

  if (body) {
    options.headers['Content-Type'] = 'application/json';
    options.body = JSON.stringify(body);
  }

  const response = await fetch(url, options);
  const data = await response.json();

  if (response.ok) {
    return data;
  } else if (response.status === 401) {
    // Token expired or invalid
    localStorage.removeItem('auth_token');
    throw new Error('Authentication expired. Please login again.');
  } else {
    throw new Error(data.message || 'Request failed');
  }
}

// Example usage
async function updateConfiguration(config) {
  try {
    await login('admin', 'admin');
    const result = await authenticatedRequest(
      'http://localhost:8080/api/admin/config',
      'PUT',
      config
    );
    console.log('Configuration updated:', result);
  } catch (error) {
    console.error('Error:', error.message);
  }
}
```

### Python

```python
import requests
import time

class InsiderMonitorClient:
    def __init__(self, base_url='http://localhost:8080'):
        self.base_url = base_url
        self.token = None
        self.token_expiry = 0

    def login(self, username='admin', password='admin'):
        response = requests.post(
            f'{self.base_url}/api/login',
            json={'username': username, 'password': password}
        )

        if response.status_code == 200:
            data = response.json()
            self.token = data['token']
            # Set expiry time to 5 minutes before actual expiry
            self.token_expiry = time.time() + (24 * 60 * 60) - 300
            return True
        else:
            raise Exception(f"Authentication failed: {response.text}")

    def get_auth_headers(self):
        if not self.token or time.time() > self.token_expiry:
            self.login()

        return {
            'Authorization': f'Bearer {self.token}'
        }

    def get_wallets(self):
        response = requests.get(f'{self.base_url}/api/wallets')
        return response.json()

    def get_config(self):
        headers = self.get_auth_headers()
        response = requests.get(
            f'{self.base_url}/api/admin/config',
            headers=headers
        )

        if response.status_code == 401:
            # Token expired, try login again and retry
            self.login()
            headers = self.get_auth_headers()
            response = requests.get(
                f'{self.base_url}/api/admin/config',
                headers=headers
            )

        return response.json()

    def update_config(self, config_update):
        headers = self.get_auth_headers()
        response = requests.put(
            f'{self.base_url}/api/admin/config',
            headers=headers,
            json=config_update
        )

        return response.json()

# Example usage
client = InsiderMonitorClient()
try:
    client.login()
    wallets = client.get_wallets()
    print("Wallets:", wallets)

    config = client.get_config()
    print("Current config:", config)

    result = client.update_config({'scan_interval': '5m'})
    print("Update result:", result)
except Exception as e:
    print(f"Error: {str(e)}")
```

## Related Topics

- [API Reference](./api.html) - Complete API documentation
- [Configuration Guide](./configuration.html) - Learn about configuration options
- [Web Interface Guide](./web-interface.html) - Learn about the web interface
