---
layout: default
title: Web Interface Guide
nav_order: 7
description: Guide to using the Solana Insider Monitor web interface
---

# Web Interface Guide
{: .no_toc }

This guide explains how to use the web interface of Solana Insider Monitor.

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>

## Getting Started

### Launching the Web Interface

To start the application with the web interface, use the `-web` flag:

```bash
go run cmd/monitor/main.go -web
```

By default, the web interface will be available at `http://localhost:8080`. You can specify a different port using the `-port` flag:

```bash
go run cmd/monitor/main.go -web -port 9090
```

### Login

When you first access the web interface, you'll be prompted to log in:

1. Enter your username (default: `admin`)
2. Enter your password (default: `admin`, or as set in the `ADMIN_PASSWORD` environment variable)
3. Click "Login"

After successful authentication, you'll receive a JWT token that will be stored in your browser and used for subsequent requests.

## Dashboard

The dashboard is the main page of the web interface and provides an overview of all monitored wallets.

### Features

- Quick overview of all monitored wallets
- Token balance summaries
- Last scan time for each wallet
- Quick links to detailed views

### Using the Dashboard

- **Refresh Data**: Click the "Refresh" button to trigger an immediate scan of all wallets
- **Wallet Details**: Click on a wallet address to view detailed information about that wallet

## Wallet Details

The wallet details page provides in-depth information about a specific wallet.

### Features

- Complete list of all tokens in the wallet
- Token balances with decimals applied
- Token symbols (when available)
- Last update time for each token

### Using Wallet Details

- **View Token Information**: See balance, symbol, and last update time for each token
- **Copy Addresses**: Click on addresses to copy them to clipboard
- **Return to Dashboard**: Click "Back to Dashboard" to return to the main view

## Configuration

The configuration page allows you to manage the application's settings.

### Features

- Update network RPC endpoint
- Manage monitored wallets
- Configure scan interval
- Set alert thresholds
- Configure Discord notifications

### Managing Configuration

#### Network Settings

1. Enter the Solana RPC endpoint URL in the "Network URL" field
2. Click "Save Changes" to apply

#### Wallet Management

1. **Add a wallet**: Enter the wallet address in the input field and click "Add Wallet"
2. **Remove a wallet**: Click the "Remove" button next to the wallet address

#### Scan Interval

1. Select the desired scan interval from the dropdown (e.g., "30s", "1m", "5m")
2. Click "Save Changes" to apply

#### Alert Settings

1. **Minimum Balance**: Enter the minimum token balance required to trigger alerts
2. **Significant Change**: Enter the percentage change required to trigger alerts (as a decimal)
3. **Ignore Tokens**: Add token mint addresses to ignore during monitoring
4. Click "Save Changes" to apply

#### Discord Notifications

1. **Enable Discord**: Toggle the switch to enable/disable Discord notifications
2. **Webhook URL**: Enter your Discord webhook URL
3. **Channel ID**: (Optional) Enter a specific Discord channel ID
4. Click "Save Changes" to apply

## Session Management

### Timeout

Your session will expire after 24 hours, after which you'll need to log in again.

### Logout

To log out manually, click the "Logout" button in the navigation bar.

## Browser Compatibility

The web interface is designed to work with modern browsers:

- Chrome (recommended)
- Firefox
- Safari
- Edge

## Troubleshooting

### Common Issues

**Issue**: Unable to log in
**Solution**: Verify you're using the correct credentials. If you've set a custom admin password via the `ADMIN_PASSWORD` environment variable, make sure you're using that password.

**Issue**: Configuration changes not saving
**Solution**: Ensure you're properly authenticated. If your session has expired, you'll need to log in again.

**Issue**: Web interface not loading
**Solution**: Verify the application is running with the `-web` flag and check that you're accessing the correct URL and port.

## Additional Features

### Keyboard Shortcuts

- **R**: Refresh data
- **Esc**: Return to previous page
- **C**: Go to configuration page
- **D**: Go to dashboard

### Mobile Responsiveness

The web interface is designed to be responsive and work well on mobile devices and tablets.

## Related Resources

- [API Reference](./api.html) - Complete API documentation
- [Configuration Guide](./configuration.html) - Detailed configuration options
- [Authentication Guide](./authentication.html) - Authentication and security information
