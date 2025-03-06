# Solana Wallet Monitor - Track Insider Trading & Token Balance Changes

Real-time monitoring tool for Solana blockchain wallets. Track wallet activities, detect potential insider trading patterns, monitor SOL & SPL token balance changes, and receive instant alerts.

[![Documentation](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://accursedgalaxy.github.io/Insider-Monitor/)
[![Discord](https://img.shields.io/discord/YOUR_DISCORD_ID?color=7289DA&label=Discord&logo=discord&logoColor=white)](https://discord.gg/7vY9ZBPdya)
[![License](https://img.shields.io/github/license/accursedgalaxy/insider-monitor)](LICENSE)

## What is Solana Insider Monitor?

Solana Insider Monitor is an open-source tool designed to help traders, developers, and blockchain enthusiasts track Solana wallet activities in real-time. Whether you're monitoring for potential insider trading signals, tracking project wallets, or analyzing token flows, this tool provides comprehensive wallet monitoring capabilities with minimal setup.

## Community

Join our growing Solana wallet monitoring community on Discord:
- Get help with setup and configuration
- Share feedback and suggestions
- Connect with other Solana wallet trackers
- Stay updated on new features and releases
- Discuss Solana development and wallet monitoring strategies

👉 [Join the Discord Server](https://discord.gg/7vY9ZBPdya)

## Key Features

- 🔍 **Multi-Wallet Monitoring** - Track multiple Solana wallets simultaneously for any activity
- 💰 **SOL & SPL Token Tracking** - Monitor all token balance changes with precision
- 🕵️ **Insider Trading Detection** - Identify suspicious wallet activity before major announcements
- ⚡ **Real-Time Alerts** - Get notified immediately of significant wallet changes
- 🔔 **Discord Integration** - Receive instant notifications in your Discord server
- 💾 **Historical Data Storage** - Maintain records of all wallet activity over time
- 🛡️ **Network Interruption Handling** - Continue monitoring even through connection issues
- 🌐 **Web Dashboard Interface** - Manage everything through an intuitive UI
- 🔑 **Secure Authentication** - Protect your monitoring setup with JWT authentication
- 🔄 **Comprehensive REST API** - Integrate with other systems through our well-documented API

## Quick Start

### Prerequisites

- Go 1.23.2 or later
- Access to a Solana RPC endpoint (mainnet, devnet, or testnet)

### Installation

```bash
# Clone the repository
git clone https://github.com/accursedgalaxy/insider-monitor
cd insider-monitor

# Install dependencies
go mod download
```

### Configuration

1. Copy the example configuration:
```bash
cp config.example.json config.json
```

2. Edit `config.json` with your wallet monitoring settings:
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

### Configuration Options

- `network_url`: Solana RPC endpoint URL
  - Mainnet: "https://api.mainnet-beta.solana.com"
  - Devnet: "https://api.devnet.solana.com"
  - Custom RPC endpoints are supported
- `wallets`: Array of Solana wallet addresses to monitor for activity
- `scan_interval`: Time between wallet scans (e.g., "30s", "1m", "5m")
- `alerts`:
  - `minimum_balance`: Minimum token balance to trigger alerts
  - `significant_change`: Percentage change to trigger alerts (0.20 = 20%)
  - `ignore_tokens`: Array of token addresses to ignore from monitoring
- `discord`:
  - `enabled`: Set to true to enable Discord wallet notifications
  - `webhook_url`: Discord webhook URL
  - `channel_id`: Discord channel ID for alerts

### Running the Solana Wallet Monitor

#### Console Mode
```bash
go run cmd/monitor/main.go
```

#### Test Mode (with mock wallet data)
```bash
go run cmd/monitor/main.go -test
```

#### Web Interface Mode
```bash
go run cmd/monitor/main.go -web
```

#### Custom Config File
```bash
go run cmd/monitor/main.go -config path/to/config.json
```

#### Custom Web Port
```bash
go run cmd/monitor/main.go -web -port 9090
```

### Alert Levels for Wallet Activity

The monitor uses three alert levels based on the configured `significant_change`:
- 🔴 **Critical**: Token balance changes >= 5x the threshold
- 🟡 **Warning**: Token balance changes >= 2x the threshold
- 🟢 **Info**: Token balance changes below 2x the threshold

### Historical Data Storage

The monitor stores Solana wallet data in the `./data` directory to:
- Prevent false alerts after monitor restarts
- Track historical wallet balance changes
- Handle Solana network interruptions gracefully
- Enable analysis of wallet activity patterns over time

## Web Dashboard Interface

The application includes a comprehensive web dashboard for easy monitoring and configuration of Solana wallets. When running with the `-web` flag, you can access the interface at `http://localhost:8080` (or your custom port).

### Dashboard Features:
- Real-time overview of all monitored Solana wallets
- Detailed view of each wallet's SOL and SPL token balances
- Historical charts of wallet activity and balance changes
- Configuration management through a user-friendly interface
- Instant alerts for significant wallet activity
- Real-time data refresh from the Solana blockchain

### Authentication:
The web interface and API use JWT authentication for secure access. The default credentials are:
- Username: `admin`
- Password: `admin` (can be changed via `ADMIN_PASSWORD` environment variable)

## API Documentation

The application provides a comprehensive REST API for programmatic access and integration with other systems for monitoring Solana wallets.

### Public Endpoints:

- `GET /api/wallets` - Get all monitored Solana wallet data
- `GET /api/wallets/{address}` - Get details for a specific Solana wallet
- `POST /api/refresh` - Trigger an immediate refresh of wallet data

### Protected Endpoints (require authentication):

- `GET /api/admin/config` - Get current wallet monitoring configuration
- `PUT /api/admin/config` - Update wallet monitoring configuration
- `POST /api/admin/wallets` - Add a new Solana wallet to monitor
- `DELETE /api/admin/wallets/{address}` - Remove a Solana wallet from monitoring

### Authentication:

Authentication is handled via JWT tokens. To obtain a token:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'
```

Use the returned token in subsequent requests:

```bash
curl -X GET http://localhost:8080/api/admin/config \
  -H "Authorization: Bearer YOUR_TOKEN"
```

For more detailed API documentation, see the [API Documentation](./docs/api.md).

### Building from Source

```bash
make build
```

The binary will be available in the `bin` directory.

## Comprehensive Documentation

📚 **[View the full documentation](https://accursedgalaxy.github.io/Insider-Monitor/)** for detailed guides, API reference, and more.

The documentation covers:
- [Installation Guide](https://accursedgalaxy.github.io/Insider-Monitor/installation)
- [Configuration Guide](https://accursedgalaxy.github.io/Insider-Monitor/configuration)
- [API Reference](https://accursedgalaxy.github.io/Insider-Monitor/api)
- [Authentication](https://accursedgalaxy.github.io/Insider-Monitor/authentication)
- [Web Interface Guide](https://accursedgalaxy.github.io/Insider-Monitor/web-interface)

## Use Cases

The Solana Insider Monitor is ideal for:

- **Traders**: Track wallets of project insiders for potential trading signals
- **Developers**: Monitor project treasury wallets for fund movements
- **Security Teams**: Detect unusual activity on critical wallets
- **Analysts**: Study token flow patterns and large holder behaviors
- **Project Owners**: Keep track of team and investor wallets

## Contributing

Contributions to improve Solana Insider Monitor are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
