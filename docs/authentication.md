---
layout: default
title: Authentication
description: Authentication and security documentation for Solana Insider Monitor
---

# Authentication

Solana Insider Monitor implements a secure authentication system to protect sensitive operations and data. This document details the authentication mechanisms and how to use them.

## Authentication Methods

The application supports two authentication methods:

1. **JWT (JSON Web Token)** - Used for API authentication
2. **Basic Authentication** - Alternative method for simple integrations

## JWT Authentication

JWT is the primary authentication method used by the web interface and recommended for API integrations.

### How JWT Authentication Works

1. Client sends credentials to the login endpoint
2. Server validates credentials and returns a signed JWT token
3. Client includes the token in subsequent requests
4. Server validates the token and processes the request

### Obtaining a JWT Token

To obtain a token, send a POST request to the login endpoint:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_password"}'
```

The response will include a token:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Using the JWT Token

Include the token in the `Authorization` header of your requests:

```bash
curl -X GET http://localhost:8080/api/admin/config \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Token Expiration

JWT tokens expire after 24 hours. After that, you'll need to obtain a new token.

## Basic Authentication

Basic authentication is supported as an alternative method for simple integrations. Note that this method is less secure than JWT and should be used only over HTTPS.

### Using Basic Authentication

Include the `Authorization` header with Base64-encoded credentials:

```bash
curl -X GET http://localhost:8080/protected/endpoint \
  -H "Authorization: Basic $(echo -n 'admin:your_password' | base64)"
```

## Default Credentials

The default credentials are:

- Username: `admin`
- Password: `admin`

It's strongly recommended to change the default password in production environments.

### Changing the Default Password

You can change the default password by setting the `ADMIN_PASSWORD` environment variable before starting the application:

```bash
export ADMIN_PASSWORD="your_secure_password"
go run cmd/monitor/main.go
```

## JWT Secret Key

By default, the application uses a predefined JWT secret key. For production environments, it's recommended to set a custom secret key.

### Setting a Custom JWT Secret

Set the `JWT_SECRET` environment variable before starting the application:

```bash
export JWT_SECRET="your_custom_jwt_secret"
go run cmd/monitor/main.go
```

## Security Best Practices

1. **Always use HTTPS** in production environments
2. **Change default credentials** before deploying to production
3. **Set a custom JWT secret** for production environments
4. **Limit access** to the API endpoints to trusted clients
5. **Regularly rotate** JWT secrets in production environments

## Future Enhancements

As part of the project roadmap, the following authentication enhancements are planned:

1. User management with different roles and permissions
2. OAuth integration for third-party authentication
3. API key-based authentication for machine-to-machine communication
