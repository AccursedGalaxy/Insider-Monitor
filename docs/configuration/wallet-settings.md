# Wallet Settings

This guide explains how to configure which Solana wallets are monitored by Solana Insider Monitor.

## Wallet Configuration

The `wallets` setting in your configuration specifies which Solana wallet addresses to monitor:

```json
{
    "wallets": [
        "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr",
        "DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF"
    ]
}
```

Each address in the array will be monitored for token balance changes according to your scan interval.

## Adding Wallets

You can add as many wallet addresses as you need to monitor. However, keep in mind that each additional wallet increases:

- API usage (potential for rate limiting)
- Processing time per scan
- Storage requirements for historical data

### Valid Wallet Formats

Solana wallet addresses:
- Are 32-44 characters long
- Use base58 encoding (containing alphanumeric characters)
- Example: `55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr`

!!! warning "Address Validation"
    Solana Insider Monitor validates wallet addresses when loading your configuration. Invalid addresses will cause errors.

## Wallet Management Tips

Here are some practical tips for managing the wallets you monitor:

### Document Your Wallets

Keep a separate document recording:
- Which wallets you're monitoring
- Why each wallet is important
- Where the wallet address was found
- The expected activity level

### Prioritize Important Wallets

When setting up alerts, focus on the most important wallets first:
- Treasury wallets holding significant funds
- Active trading wallets with frequent changes
- Protocol wallets controlling critical functions

### Regular Auditing

Periodically review your wallet list to:
- Remove wallets no longer needed for monitoring
- Add new wallets of interest
- Verify address accuracy

## Per-Wallet Scan Configuration

Solana Insider Monitor supports configuring different scan modes for individual wallets, giving you fine-grained control over what tokens are monitored.

### Scan Modes

Each wallet can use one of three scan modes:

- **all**: Monitor all tokens in the wallet (default)
- **whitelist**: Only monitor specific tokens listed in `include_tokens`
- **blacklist**: Monitor all tokens except those listed in `exclude_tokens`

### Configuration Example

Here's how to configure different scan modes for different wallets:

```json
{
  "wallet_configs": {
    "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr": {
      "scan": {
        "scan_mode": "whitelist",
        "include_tokens": [
          "So11111111111111111111111111111111111111112",  // SOL
          "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"  // USDC
        ],
        "exclude_tokens": []
      }
    },
    "DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF": {
      "scan": {
        "scan_mode": "blacklist",
        "include_tokens": [],
        "exclude_tokens": [
          "7dHbWXmci3dT8UFYWYZweBLXgycu7Y3iL6trKn1Y7ARj",  // Dust token
          "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"   // BONK
        ]
      }
    }
  }
}
```

### Use Cases

- **Whitelist Mode**: Use for wallets where you only care about specific tokens (e.g., stablecoins, major tokens)
- **Blacklist Mode**: Use for wallets with many dust or airdrop tokens you want to ignore
- **All Mode**: Use for comprehensive monitoring of all token activity

### Global vs Per-Wallet Configuration

You can set a global scan mode that applies to all wallets and override it for specific wallets:

```json
{
  "scan": {
    "scan_mode": "all",
    "include_tokens": [],
    "exclude_tokens": []
  },
  "wallet_configs": {
    "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr": {
      "scan": {
        "scan_mode": "whitelist",
        "include_tokens": ["So11111111111111111111111111111111111111112"]
      }
    }
  }
}
```

In this example, all wallets use "all" mode except for the specified wallet, which uses "whitelist" mode.

## Finding Wallet Addresses

You can find Solana wallet addresses through various methods:

1. **Solana Explorer**: Search for a wallet on [Solana Explorer](https://explorer.solana.com/)
2. **Wallet Applications**: Export wallet addresses from Phantom, Solflare, etc.
3. **Project Documentation**: Many projects publish their treasury wallet addresses
4. **Block Explorers**: Search for transactions to discover related addresses

## Monitoring Considerations

### Privacy

Remember that all Solana blockchain data is public. Monitoring wallets:
- Does not require permission from the wallet owner
- Does not provide any private key access
- Only tracks public, on-chain information

### Performance Impact

The number of wallets you monitor affects performance:

| Number of Wallets | Impact | Recommended Scan Interval |
|-------------------|--------|---------------------------|
| 1-10 | Minimal | 30s - 1m |
| 10-50 | Moderate | 1m - 5m |
| 50-100 | Significant | 5m - 15m |
| 100+ | High | 15m+ |

## Configuration Methods

### Via Configuration File

Edit the wallets array in your `config.json` file:

```json
{
    "wallets": [
        "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr",
        "DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF"
    ]
}
```

### Via Web Interface

When running in web mode:

1. Navigate to Settings > Wallets
2. Add or remove wallet addresses
3. Save your changes

The monitor will automatically begin tracking new wallets without requiring a restart.

## Special Wallet Types

### Token Program Wallets

To monitor SPL token program activities, add the Solana token program address:

```json
{
    "wallets": [
        "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
    ]
}
```

### Program ID Monitoring

You can monitor program deployments and interactions by adding program IDs:

```json
{
    "wallets": [
        "YOUR_PROGRAM_ID"
    ]
}
```

## Best Practices

1. **Start Small**: Begin with a few important wallets before scaling up
2. **Review Regularly**: Periodically audit your wallet list to remove unnecessary addresses
3. **Document Sources**: Record why each wallet is being monitored and its significance
4. **Test First**: When adding critical wallets, verify monitoring works as expected
5. **Balance Performance**: Find the right balance between comprehensive monitoring and system performance

## Future Enhancements

In future versions, we plan to add features like:

- Wallet labeling for easier identification
- Grouping wallets by category
- Custom alert thresholds per wallet
- Multiple configuration profiles

## Related Settings

- [Network Settings](network-settings.md) - Configure RPC endpoints and scan intervals
- [Alert Settings](alert-settings.md) - Set up alert thresholds
