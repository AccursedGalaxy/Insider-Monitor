# Network Settings

This guide explains how to configure network-related settings for Solana Insider Monitor, including RPC endpoints and scan intervals.

## Network URL

The `network_url` setting specifies the Solana RPC endpoint to connect to. This is a crucial setting as it determines which Solana network your monitor will connect to and how reliable your connection will be.

```json
{
    "network_url": "https://api.mainnet-beta.solana.com"
}
```

### Available Networks

You can configure different Solana networks:

| Network | Example URL | Description |
|---------|-------------|-------------|
| Mainnet | `https://api.mainnet-beta.solana.com` | Production Solana network |
| Devnet | `https://api.devnet.solana.com` | Development test network |
| Testnet | `https://api.testnet.solana.com` | Testing network |

### Private RPC Providers

For production monitoring, we recommend using a private RPC provider for better reliability and performance:

- [QuickNode](https://www.quicknode.com/chains/sol)
- [Alchemy](https://www.alchemy.com/solana)
- [Ankr](https://www.ankr.com/rpc/solana/)
- [Helius](https://helius.xyz/)

Using a private RPC endpoint gives you:
- Higher rate limits
- Better performance
- More reliable connectivity
- Additional features like archive data

## Scan Interval

The `scan_interval` setting determines how frequently the monitor checks wallets for changes:

```json
{
    "scan_interval": "1m"
}
```

### Format

Scan intervals use Go's duration string format:

| Example | Meaning |
|---------|---------|
| `5s` | 5 seconds |
| `30s` | 30 seconds |
| `1m` | 1 minute |
| `5m` | 5 minutes |
| `1h` | 1 hour |

### Considerations

When choosing a scan interval, consider:

- **RPC Rate Limits**: Shorter intervals mean more frequent API calls
- **Number of Wallets**: More wallets require more API calls per scan
- **Detection Speed**: Shorter intervals detect changes sooner
- **Resource Usage**: Shorter intervals increase CPU and network usage

### Recommended Settings

| Use Case | Recommended Interval | Notes |
|----------|---------------------|-------|
| Critical wallet monitoring | `30s` to `1m` | Requires private RPC |
| Standard monitoring | `2m` to `5m` | Good balance for most uses |
| Low-priority monitoring | `15m` to `1h` | Minimal resource usage |

## Advanced Network Configuration

### Fault Tolerance

Solana Insider Monitor includes automatic retry logic for handling temporary network issues:

- Automatic reconnection if the connection drops
- Exponential backoff for failed requests
- Alert notifications for persistent connection problems

### Multiple RPC Endpoints

For high-availability setups, you can configure fallback RPC endpoints by modifying the source code. This allows the monitor to switch to an alternative endpoint if the primary one fails.

## Configuration Methods

### Via Configuration File

Edit the network settings in your `config.json` file:

```json
{
    "network_url": "https://api.mainnet-beta.solana.com",
    "scan_interval": "1m"
}
```

### Via Environment Variables

You can override the configuration using environment variables:

```bash
export NETWORK_URL="https://api.mainnet-beta.solana.com"
export SCAN_INTERVAL="30s"
```

### Via Web Interface

When running in web mode:

1. Navigate to Settings > Network
2. Update the RPC endpoint and scan interval
3. Save your changes

## Troubleshooting

### Connection Issues

If you experience connection problems:

1. **Verify the RPC URL** is correct and accessible
2. **Try a different RPC endpoint** to rule out provider-specific issues
3. **Increase the scan interval** to reduce API call frequency
4. **Check for rate limiting** in the console logs

### Performance Issues

If scanning is slow:

1. **Use a private RPC endpoint** for better performance
2. **Reduce the number of monitored wallets**
3. **Increase the scan interval**
4. **Check your network latency** to the RPC endpoint

## Related Settings

- [Wallet Settings](wallet-settings.md) - Configure which wallets to monitor
- [Alert Settings](alert-settings.md) - Set up alert thresholds
