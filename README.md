# Solana Insider Monitor

A tool for monitoring Solana wallet activities, detecting balance changes, and receiving real-time alerts.

## Community

Join our Discord community to:
- Get help with setup and configuration
- Share feedback and suggestions
- Connect with other users
- Stay updated on new features and releases
- Discuss Solana development

👉 [Join the Discord Server](https://discord.gg/7vY9ZBPdya)

## Features

- 🔍 Monitor multiple Solana wallets simultaneously
- 💰 Track token balance changes
- ⚡ Real-time alerts for significant changes
- 🔔 Discord integration for notifications
- 💾 Persistent storage of wallet data
- 🛡️ Graceful handling of network interruptions

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

2. Edit `config.json` with your settings:
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
- `wallets`: Array of Solana wallet addresses to monitor
- `scan_interval`: Time between scans (e.g., "30s", "1m", "5m")
- `alerts`:
  - `minimum_balance`: Minimum token balance to trigger alerts
  - `significant_change`: Percentage change to trigger alerts (0.20 = 20%)
  - `ignore_tokens`: Array of token addresses to ignore
- `discord`:
  - `enabled`: Set to true to enable Discord notifications
  - `webhook_url`: Discord webhook URL
  - `channel_id`: Discord channel ID
- `scan`:
  - `scan_mode`: Token scanning mode
    - `"all"`: Monitor all tokens (default)
    - `"whitelist"`: Only monitor tokens in `include_tokens`
    - `"blacklist"`: Monitor all tokens except those in `exclude_tokens`
  - `include_tokens`: Array of token addresses to specifically monitor (used with `whitelist` mode)
  - `exclude_tokens`: Array of token addresses to ignore (used with `blacklist` mode)

### Scan Mode Examples

Here are examples of different scan configurations:

1. Monitor all tokens:
```json
{
    "scan": {
        "scan_mode": "all",
        "include_tokens": [],
        "exclude_tokens": []
    }
}
```

2. Monitor only specific tokens (whitelist):
```json
{
    "scan": {
        "scan_mode": "whitelist",
        "include_tokens": [
            "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",  // USDC
            "So11111111111111111111111111111111111111112"     // SOL
        ],
        "exclude_tokens": []
    }
}
```

3. Monitor all tokens except specific ones (blacklist):
```json
{
    "scan": {
        "scan_mode": "blacklist",
        "include_tokens": [],
        "exclude_tokens": [
            "TokenAddressToIgnore1",
            "TokenAddressToIgnore2"
        ]
    }
}
```

### Running the Monitor

```bash
go run cmd/monitor/main.go
```

#### Custom Config File
```bash
go run cmd/monitor/main.go -config path/to/config.json
```

### Alert Levels

The monitor uses three alert levels based on the configured `significant_change`:
- 🔴 **Critical**: Changes >= 5x the threshold
- 🟡 **Warning**: Changes >= 2x the threshold
- 🟢 **Info**: Changes below 2x the threshold

### Data Storage

The monitor stores wallet data in the `./data` directory to:
- Prevent false alerts after restarts
- Track historical changes
- Handle network interruptions gracefully

### Building from Source

```bash
make build
```

The binary will be available in the `bin` directory.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
