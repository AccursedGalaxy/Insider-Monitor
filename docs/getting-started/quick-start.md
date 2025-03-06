# Quick Start Guide

This guide will help you get Solana Insider Monitor up and running quickly. For a more detailed setup, refer to the [Installation Guide](installation.md).

## 5-Minute Setup

Follow these steps to have Solana Insider Monitor running in less than 5 minutes.

### Step 1: Install

Clone the repository and install dependencies:

```bash
git clone https://github.com/accursedgalaxy/insider-monitor.git
cd insider-monitor
go mod download
```

### Step 2: Configure

Create a basic configuration file:

```bash
cp config.example.json config.json
```

Edit `config.json` with your favorite text editor:

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

!!! warning "Replace Wallet Addresses"
    Make sure to replace `YOUR_WALLET_ADDRESS_1` and `YOUR_WALLET_ADDRESS_2` with actual Solana wallet addresses you want to monitor.

### Step 3: Run

Start the monitor in console mode:

```bash
go run cmd/monitor/main.go
```

Or with the web interface:

```bash
go run cmd/monitor/main.go -web
```

If you've built the binary:

```bash
./bin/insider-monitor -web
```

### Step 4: Access the Web Interface

If you started with the `-web` flag, open your browser and navigate to:

```
http://localhost:8080
```

Default login credentials:
- Username: `admin`
- Password: `admin`

!!! danger "Security Notice"
    Change the default password after your first login by setting the `ADMIN_PASSWORD` environment variable.

## What's Next?

Now that you have the basic setup running, you can:

- [Configure alert thresholds](../configuration/alert-settings.md)
- [Set up Discord notifications](../configuration/discord-integration.md)
- [Learn about different running modes](running-modes.md)
- [Explore the API](../api/index.md) for programmatic access

## Running in Test Mode

For testing purposes, you can run the monitor in test mode, which uses mock data and accelerated scanning:

```bash
go run cmd/monitor/main.go -test
```

This is useful for:
- Testing your setup without connecting to a real network
- Developing new features
- Verifying alerts are working properly

## Common Issues

### Connection Problems

If you're having trouble connecting to the Solana network:

- Verify your network URL is correct
- Check if the RPC endpoint has rate limits
- Consider using a paid RPC service for better reliability

### No Alerts

If you're not receiving alerts:

- Ensure your alert thresholds are set appropriately
- Check that the wallets have activity
- Verify the `minimum_balance` isn't set too high

See the [Troubleshooting](../troubleshooting.md) section for more help.
