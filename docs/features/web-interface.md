# Web Interface

Solana Insider Monitor includes a comprehensive web interface that provides a visual dashboard for monitoring wallet activities, managing configuration, and viewing alerts.

## Overview

The web interface provides a user-friendly way to interact with Solana Insider Monitor without relying on command-line tools. It offers real-time data visualization, configuration management, and an overview of monitored wallets.

## Getting Started

### Enabling the Web Interface

To start Solana Insider Monitor with the web interface enabled, use the `-web` flag:

```bash
./bin/monitor -web
```

By default, the web server runs on port 8080. You can specify a different port using the `-port` flag:

```bash
./bin/monitor -web -port 3000
```

### Accessing the Web Interface

Once the monitor is running with the web interface enabled, access it by navigating to:

```
http://localhost:8080
```

(Replace 8080 with your custom port if specified)

### Authentication

The web interface requires authentication to prevent unauthorized access. Default credentials are:

- **Username**: admin
- **Password**: password

!!! warning "Default Credentials"
    It's highly recommended to change the default credentials by setting the `ADMIN_USERNAME` and `ADMIN_PASSWORD` environment variables before running in production.

## Dashboard Features

### Wallet Overview

The dashboard home page displays an overview of all monitored wallets, showing:

- Total number of wallets being monitored
- Summary of recent balance changes
- Alert statistics
- System status information

![Dashboard Overview](../assets/images/dashboard-overview.png)

### Wallet Details

Clicking on a wallet address opens a detailed view showing:

- Complete token balance list
- Historical balance changes
- Recent transactions
- Alert history specific to the wallet

### Real-time Updates

The web interface updates automatically when new data is available, providing a real-time view of wallet activities without requiring manual refreshes.

## Configuration Management

### Editing Configuration

The web interface includes a configuration editor that allows you to:

- Add or remove wallets from monitoring
- Adjust alert thresholds
- Configure Discord integration
- Modify network settings

Changes made through the web interface are applied immediately without requiring a restart of the monitor.

![Configuration Editor](../assets/images/config-editor.png)

### Import/Export

You can import or export configuration files directly through the web interface:

- **Import**: Upload a JSON configuration file to apply settings
- **Export**: Download the current configuration as a JSON file for backup or transferring to another instance

## Alert Management

### Viewing Alerts

The alerts page displays all recent alerts with:

- Severity indicator
- Affected wallet and token
- Change amount and percentage
- Timestamp
- Filtering options by date, wallet, or severity

### Alert Settings

Adjust alert thresholds and notification settings through the web interface:

- Set percentage thresholds for different alert levels
- Configure minimum value thresholds
- Enable or disable notification channels
- Set cooldown periods to prevent alert spam

## Mobile Responsiveness

The web interface is designed to be responsive and works well on mobile devices, allowing you to monitor your wallets on the go.

## API Access

The web interface includes a built-in API explorer where you can:

- Test API endpoints
- Generate API tokens
- View request/response examples

For more information on the API, see the [API Reference](../api/index.md).

## Security Considerations

The web interface implements several security measures:

- JWT-based authentication
- HTTPS support (when configured with SSL certificates)
- Automatic session timeout after period of inactivity
- Brute force protection with request limiting

## Troubleshooting

If you encounter issues with the web interface:

1. Ensure the server is running with the `-web` flag
2. Check network/firewall settings if accessing remotely
3. Clear browser cache if you see outdated information
4. Check console logs for error messages

For more help, see the [Troubleshooting](../troubleshooting.md) guide.
