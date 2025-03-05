# Alert System

Solana Insider Monitor includes a powerful alert system that notifies you about significant changes in wallet balances and activities. This guide explains how the alert system works and how to customize it for your needs.

## Overview

The alert system monitors wallet balances for changes and generates notifications based on configurable thresholds. Alerts can be delivered through multiple channels, including console output and Discord integration.

## Alert Levels

Alerts are categorized into three severity levels:

- 🟢 **Info**: Minor changes that you want to be aware of but don't require immediate attention
- 🟡 **Warning**: Moderate changes that might warrant investigation
- 🔴 **Critical**: Major changes that require immediate attention

## Alert Channels

### Console Alerts

By default, all alerts are displayed in the console output with color-coded severity levels:

```bash
[INFO] Wallet 6VJ...k2F: USDC balance increased by 100.0 (+5%)
[WARNING] Wallet AeW...a3f: SOL balance decreased by 10.5 (-15%)
[CRITICAL] Wallet L3m...8Zq: RAY balance decreased by 5000.0 (-75%)
```

### Discord Alerts

Alerts can be sent to a Discord channel using webhooks. Discord alerts include rich, formatted messages with:

- Wallet address with explorer link
- Token name and symbol
- Balance change amount and percentage
- Timestamp of the change
- Alert severity indicator

![Discord Alert Example](../assets/images/discord-alert-example.png)

## Configuring Alerts

### Alert Thresholds

You can configure the thresholds for each alert level in the configuration file:

```json
"alerts": {
  "info_threshold_percent": 5,
  "warning_threshold_percent": 15,
  "critical_threshold_percent": 25,
  "min_value_threshold": 10
}
```

- **Info threshold**: Percentage change that triggers an info alert (default: 5%)
- **Warning threshold**: Percentage change that triggers a warning alert (default: 15%)
- **Critical threshold**: Percentage change that triggers a critical alert (default: 25%)
- **Minimum value threshold**: Minimum USD value required to trigger any alert (default: $10)

### Customizing Discord Alerts

To enable Discord alerts, configure the Discord webhook in your config file:

```json
"discord": {
  "enabled": true,
  "webhook_url": "https://discord.com/api/webhooks/your-webhook-url",
  "username": "Solana Monitor",
  "avatar_url": "https://solana.com/src/img/branding/solanaLogoMark.svg",
  "mention_role_id": ""
}
```

For more detailed configuration options, see the [Discord Integration](../configuration/discord-integration.md) guide.

## Alert Filtering

You can filter alerts based on:

- **Specific tokens**: Only receive alerts for certain tokens
- **Minimum value**: Ignore changes below a certain USD value
- **Specific wallets**: Prioritize alerts from important wallets

Configure these filters in the alert settings section of your config file.

## Advanced Alert Features

### Cooldown Periods

To prevent alert spam, the system implements cooldown periods for repeated alerts from the same wallet/token combination:

```json
"alerts": {
  "cooldown_minutes": 60
}
```

This setting prevents the same alert from being triggered more than once within the specified time period.

### Aggregate Alerts

When multiple small changes occur within a short time frame, the system can aggregate them into a single alert to reduce noise:

```json
"alerts": {
  "aggregate_changes": true,
  "aggregate_window_minutes": 10
}
```

## Troubleshooting Alerts

If you're not receiving alerts as expected:

1. Check your threshold configurations
2. Verify that the token price information is available
3. Ensure Discord webhook URLs are correct (if using Discord)
4. Check the console output for any error messages related to alert delivery

For more help, see the [Troubleshooting](../troubleshooting.md) guide.

## Alert System API

For programmatic access to alerts, you can use the API endpoints:

- `GET /api/alerts` - Retrieve recent alerts
- `GET /api/alerts/settings` - Get alert configuration
- `PUT /api/alerts/settings` - Update alert configuration

See the [API Reference](../api/index.md) for more details.
