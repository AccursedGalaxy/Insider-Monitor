# Alert Settings

This guide explains how to configure alert thresholds and settings for Solana Insider Monitor.

## Alert Configuration Options

Alerts in Solana Insider Monitor are highly configurable. You can adjust the following settings:

### Minimum Balance

The `minimum_balance` setting specifies the minimum token balance (in token units) required to trigger alerts. Tokens with balances below this threshold won't generate alerts, even if they experience significant changes.

```json
{
    "alerts": {
        "minimum_balance": 1000
    }
}
```

This is useful for:
- Ignoring dust amounts
- Reducing alert noise
- Focusing on more significant holdings

### Significant Change Threshold

The `significant_change` setting defines the percentage change in token balance required to trigger an alert. This is expressed as a decimal (e.g., 0.20 = 20%).

```json
{
    "alerts": {
        "significant_change": 0.20
    }
}
```

Lower values make the monitor more sensitive to small changes:
- `0.05` (5%) - Very sensitive, good for critical wallets
- `0.20` (20%) - Moderate sensitivity, good for general monitoring
- `0.50` (50%) - Low sensitivity, only major changes

### Ignored Tokens

You can specify token mint addresses to exclude from monitoring using the `ignore_tokens` array:

```json
{
    "alerts": {
        "ignore_tokens": [
            "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",  // USDC
            "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"   // USDT
        ]
    }
}
```

This is useful for:
- Ignoring stablecoins
- Excluding airdrops or spam tokens
- Focusing on specific assets of interest

## Alert Levels

Solana Insider Monitor categorizes alerts into three severity levels:

1. **Info** (🔵): Minor changes slightly above the threshold
2. **Warning** (🟡): Moderate changes (typically 1.5x the threshold)
3. **Critical** (🔴): Major changes (typically 2.5x the threshold)

## Configuration Methods

### Via Configuration File

Edit the `alerts` section in your `config.json` file:

```json
{
    "alerts": {
        "minimum_balance": 1000,
        "significant_change": 0.20,
        "ignore_tokens": [
            "TOKEN_MINT_ADDRESS_1"
        ]
    }
}
```

### Via Web Interface

When running in web mode:

1. Navigate to Settings > Alerts
2. Adjust the threshold values
3. Add or remove ignored tokens
4. Save your changes

Alert settings will be applied immediately without requiring a restart.

### Via Environment Variables

You can override the configuration using environment variables:

```bash
export MIN_BALANCE=500
export SIGNIFICANT_CHANGE=0.10
```

## Best Practices

1. **Start Conservative**: Begin with higher thresholds and gradually lower them
2. **Regular Review**: Periodically audit your alert settings as your needs change
3. **Custom Thresholds**: Consider using different thresholds for different wallets
4. **Combine with Filtering**: Use token filtering alongside thresholds for precision

## Example Configurations

### High Sensitivity Monitoring

```json
{
    "alerts": {
        "minimum_balance": 100,
        "significant_change": 0.05,
        "ignore_tokens": []
    }
}
```

### Balanced Monitoring

```json
{
    "alerts": {
        "minimum_balance": 1000,
        "significant_change": 0.20,
        "ignore_tokens": [
            "USDC_MINT_ADDRESS",
            "USDT_MINT_ADDRESS"
        ]
    }
}
```

### Low Noise Monitoring

```json
{
    "alerts": {
        "minimum_balance": 5000,
        "significant_change": 0.40,
        "ignore_tokens": [
            "USDC_MINT_ADDRESS",
            "USDT_MINT_ADDRESS",
            "OTHER_STABLECOINS"
        ]
    }
}
```

## Related Settings

- [Discord Integration](discord-integration.md) - Configure alert delivery via Discord
- [Network Settings](network-settings.md) - Adjust scanning frequency
- [Wallet Settings](wallet-settings.md) - Manage monitored wallets
