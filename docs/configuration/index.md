# Configuration Guide

Solana Insider Monitor is highly configurable to suit your monitoring needs. This guide explains all configuration options and how to use them effectively.

## Configuration File Format

Solana Insider Monitor uses a JSON configuration file. By default, it looks for `config.json` in the current directory.

Here's the complete configuration file structure with all available options:

```json
{
    "network_url": "https://api.mainnet-beta.solana.com",
    "wallets": [
        "WALLET_ADDRESS_1",
        "WALLET_ADDRESS_2"
    ],
    "scan_interval": "1m",
    "alerts": {
        "minimum_balance": 1000,
        "significant_change": 0.20,
        "ignore_tokens": [
            "TOKEN_MINT_ADDRESS_1"
        ]
    },
    "discord": {
        "enabled": false,
        "webhook_url": "https://discord.com/api/webhooks/your-webhook-url",
        "channel_id": "your-channel-id"
    }
}
```

## Configuration Options

<div class="grid cards" markdown>

-   :material-web:{ .lg .middle } __Network URL__

    ---

    Solana RPC endpoint URL for connecting to the blockchain.

    [:octicons-arrow-right-24: Network Settings](network-settings.md)

-   :material-wallet:{ .lg .middle } __Wallets__

    ---

    List of Solana wallet addresses to monitor.

    [:octicons-arrow-right-24: Wallet Settings](wallet-settings.md)

-   :material-bell-alert:{ .lg .middle } __Alerts__

    ---

    Configure alert thresholds and filter settings.

    [:octicons-arrow-right-24: Alert Settings](alert-settings.md)

-   :fontawesome-brands-discord:{ .lg .middle } __Discord__

    ---

    Discord webhook integration for notifications.

    [:octicons-arrow-right-24: Discord Integration](discord-integration.md)

</div>

## Using Configuration Files

### Loading a Custom Configuration File

You can specify a custom configuration file path:

```bash
insider-monitor -config /path/to/your/config.json
```

### Hot Reloading

Solana Insider Monitor supports hot reloading of configuration changes. When running with the web interface, you can update the configuration through the UI, and changes will be applied without requiring a restart.

### Environment Variables

You can override configuration file settings with environment variables:

| Environment Variable | Maps to | Example |
|----------------------|---------|---------|
| `NETWORK_URL` | `network_url` | `export NETWORK_URL="https://api.mainnet-beta.solana.com"` |
| `SCAN_INTERVAL` | `scan_interval` | `export SCAN_INTERVAL="30s"` |
| `ADMIN_PASSWORD` | Web UI password | `export ADMIN_PASSWORD="secure-password"` |
| `DISCORD_WEBHOOK_URL` | `discord.webhook_url` | `export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/..."` |
| `MIN_BALANCE` | `alerts.minimum_balance` | `export MIN_BALANCE="500"` |
| `SIGNIFICANT_CHANGE` | `alerts.significant_change` | `export SIGNIFICANT_CHANGE="0.1"` |

## Example Configurations

### Basic Monitoring

```json
{
    "network_url": "https://api.mainnet-beta.solana.com",
    "wallets": ["YOUR_WALLET_ADDRESS"],
    "scan_interval": "5m",
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

### Advanced Discord Alerts

```json
{
    "network_url": "https://api.mainnet-beta.solana.com",
    "wallets": [
        "WALLET_ADDRESS_1",
        "WALLET_ADDRESS_2"
    ],
    "scan_interval": "1m",
    "alerts": {
        "minimum_balance": 500,
        "significant_change": 0.10,
        "ignore_tokens": [
            "USDC_MINT_ADDRESS",
            "USDT_MINT_ADDRESS"
        ]
    },
    "discord": {
        "enabled": true,
        "webhook_url": "https://discord.com/api/webhooks/your-webhook-url",
        "channel_id": "your-channel-id"
    }
}
```

## Configuration via Web Interface

When running with the `-web` flag, you can configure the monitor through the web interface:

1. Navigate to `http://localhost:8080/config`
2. Log in with your credentials
3. Modify settings through the UI
4. Save changes, which will be applied immediately

![Configuration Web Interface](../assets/images/config-ui-example.png)

## Configuration Best Practices

1. **Start Small**: Begin with a few important wallets before scaling up
2. **Tune Alert Thresholds**: Adjust `minimum_balance` and `significant_change` to reduce noise
3. **Secure Sensitive Information**: Use environment variables for webhook URLs and passwords
4. **Regular Updates**: Review and update your configuration as monitoring needs change
5. **Use Meaningful Comments**: Add comments to your configuration for team collaboration

## Next Steps

- [Network Settings](network-settings.md) - Configuring network connections
- [Wallet Settings](wallet-settings.md) - Managing wallet monitoring
- [Alert Settings](alert-settings.md) - Customizing alert thresholds
- [Discord Integration](discord-integration.md) - Setting up Discord notifications
