# Running Modes

Solana Insider Monitor can be run in several different modes to suit your specific needs. This guide explains each mode and how to use it.

## Console Mode

Console mode is the most basic way to run Solana Insider Monitor. In this mode, the application runs in your terminal, scanning wallets according to your configuration, and outputting alerts directly to the console.

### When to Use Console Mode

- For quick monitoring sessions
- In server environments without a GUI
- For debugging and testing configurations
- When integrating with other command-line tools

### How to Run in Console Mode

```bash
go run cmd/monitor/main.go
# or if you've built the binary
./bin/insider-monitor
```

### Example Output

```
[INFO] 2023-05-01T12:34:56Z: Scanning wallets...
[INFO] 2023-05-01T12:35:01Z: Found token balance change for wallet 55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr
[WARNING] 2023-05-01T12:35:01Z: Token SOL decreased by 25% in wallet 55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr
```

## Web Interface Mode

Web interface mode starts the monitor with a web-based UI, allowing you to visualize wallet data, manage configuration, and view alerts through a browser.

### When to Use Web Interface Mode

- For visual monitoring of multiple wallets
- When you need a user-friendly interface
- For team environments where multiple people need access
- When you want to view historical data and trends

### How to Run with Web Interface

```bash
go run cmd/monitor/main.go -web
# or if you've built the binary
./bin/insider-monitor -web
```

By default, the web interface runs on port 8080. You can specify a different port:

```bash
go run cmd/monitor/main.go -web -port 9090
```

### Accessing the Web Interface

Once running, access the web interface by navigating to:
```
http://localhost:8080
```

!!! info "Authentication Required"
    The web interface requires authentication. Default credentials are:
    - Username: `admin`
    - Password: `admin`

    Change the default password by setting the `ADMIN_PASSWORD` environment variable.

### Web Interface Features

The web interface includes:

- Dashboard with overview of all monitored wallets
- Detailed view of each wallet's token balances
- Configuration management
- Alert history and filtering
- Real-time data refresh

## Test Mode

Test mode runs the monitor with mock data and accelerated scanning, useful for testing your setup without connecting to a real network.

### When to Use Test Mode

- During initial setup to verify configuration
- When developing new features
- For demonstrations and presentations
- Testing alert thresholds and notification systems

### How to Run in Test Mode

```bash
go run cmd/monitor/main.go -test
# or if you've built the binary
./bin/insider-monitor -test
```

You can combine test mode with web interface mode:

```bash
go run cmd/monitor/main.go -test -web
```

## Custom Configuration File

By default, Solana Insider Monitor looks for `config.json` in the current directory. You can specify a different configuration file:

```bash
go run cmd/monitor/main.go -config /path/to/custom-config.json
```

## Scan Configuration Modes

Solana Insider Monitor supports different scanning modes to control which tokens are monitored:

### Global Scan Settings

You can set a global scan mode that applies to all wallets:

```json
{
  "scan": {
    "scan_mode": "all",              // "all", "whitelist", or "blacklist"
    "include_tokens": [],            // Used with "whitelist" mode
    "exclude_tokens": []             // Used with "blacklist" mode
  }
}
```

### Per-Wallet Scan Settings

You can also configure different scan modes for individual wallets:

```json
{
  "wallet_configs": {
    "52C9T2T7JRojtxumYnYZhyUmrN7kqzvCLc4Ksvjk7TxD": {
      "scan": {
        "scan_mode": "whitelist",
        "include_tokens": [
          "So11111111111111111111111111111111111111112",  // SOL
          "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"  // USDC
        ],
        "exclude_tokens": []
      }
    }
  }
}
```

### Scan Modes Explained

- **all**: Monitor all tokens in the wallet (default)
- **whitelist**: Only monitor tokens specified in the `include_tokens` list
- **blacklist**: Monitor all tokens except those in the `exclude_tokens` list

This system allows you to precisely control which tokens trigger alerts, whether you want to focus on specific tokens or exclude known noisy ones.

## Environment Variables

Solana Insider Monitor supports configuration via environment variables, which can be useful in containerized environments:

| Environment Variable | Description |
| ------------------- | ----------- |
| `NETWORK_URL` | Solana RPC endpoint URL |
| `SCAN_INTERVAL` | Time between scans (e.g., "30s", "1m") |
| `ADMIN_PASSWORD` | Password for web interface admin access |
| `DISCORD_WEBHOOK_URL` | Discord webhook for notifications |
| `MIN_BALANCE` | Minimum token balance to trigger alerts |
| `SIGNIFICANT_CHANGE` | Percentage change to trigger alerts (e.g., 0.20 = 20%) |

### Using Environment Variables

```bash
export NETWORK_URL="https://api.mainnet-beta.solana.com"
export ADMIN_PASSWORD="your-secure-password"
go run cmd/monitor/main.go -web
```

Or in one line:

```bash
NETWORK_URL="https://api.mainnet-beta.solana.com" ADMIN_PASSWORD="your-secure-password" go run cmd/monitor/main.go -web
```

## Running as a Service

For production environments, you might want to run Solana Insider Monitor as a system service.

### Systemd Service Example

Create a file at `/etc/systemd/system/solana-monitor.service`:

```ini
[Unit]
Description=Solana Insider Monitor
After=network.target

[Service]
User=solana
WorkingDirectory=/opt/solana-monitor
ExecStart=/opt/solana-monitor/bin/insider-monitor -web
Restart=on-failure
Environment=ADMIN_PASSWORD=your-secure-password

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable solana-monitor
sudo systemctl start solana-monitor
```

Check status:

```bash
sudo systemctl status solana-monitor
```

## Logging

By default, Solana Insider Monitor logs to stdout. You can redirect logs to a file:

```bash
go run cmd/monitor/main.go > monitor.log 2>&1
```

For more advanced logging configurations, consider using a log management solution like logrotate.
