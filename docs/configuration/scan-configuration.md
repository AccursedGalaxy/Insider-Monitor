# Scan Configuration

This guide explains how to use Solana Insider Monitor's scan configuration system to control exactly which tokens are monitored for each wallet.

## Overview

The scan configuration system allows you to:

1. Set a global scan mode that applies to all wallets
2. Configure different scan modes for individual wallets
3. Control precisely which tokens are monitored or ignored

This provides a powerful way to focus your monitoring on tokens that matter to you while reducing noise from irrelevant tokens.

## Scan Modes

Solana Insider Monitor supports three scan modes:

| Mode | Description |
|------|-------------|
| `all` | Monitor all tokens in the wallet (default) |
| `whitelist` | Only monitor tokens explicitly included in the `include_tokens` list |
| `blacklist` | Monitor all tokens except those excluded in the `exclude_tokens` list |

## Global Scan Configuration

To set a scan mode that applies to all wallets, use the `scan` section in your config:

```json
{
  "scan": {
    "scan_mode": "all",              // "all", "whitelist", or "blacklist"
    "include_tokens": [],            // Used with "whitelist" mode
    "exclude_tokens": []             // Used with "blacklist" mode
  }
}
```

### Examples

#### Monitor All Tokens (Default)

```json
{
  "scan": {
    "scan_mode": "all",
    "include_tokens": [],
    "exclude_tokens": []
  }
}
```

#### Whitelist Mode (Only Monitor Specific Tokens)

```json
{
  "scan": {
    "scan_mode": "whitelist",
    "include_tokens": [
      "So11111111111111111111111111111111111111112",  // SOL
      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"  // USDC
    ],
    "exclude_tokens": []
  }
}
```

#### Blacklist Mode (Exclude Specific Tokens)

```json
{
  "scan": {
    "scan_mode": "blacklist",
    "include_tokens": [],
    "exclude_tokens": [
      "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263",  // BONK
      "7dHbWXmci3dT8UFYWYZweBLXgycu7Y3iL6trKn1Y7ARj"   // Some dust token
    ]
  }
}
```

## Per-Wallet Scan Configuration

You can override the global scan settings for specific wallets using the `wallet_configs` section:

```json
{
  "wallet_configs": {
    "YOUR_WALLET_ADDRESS": {
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

### Mixed Configuration Example

```json
{
  "scan": {
    "scan_mode": "all",              // Default mode for all wallets
    "include_tokens": [],
    "exclude_tokens": []
  },
  "wallet_configs": {
    "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr": {
      "scan": {
        "scan_mode": "whitelist",    // This wallet uses whitelist mode
        "include_tokens": [
          "So11111111111111111111111111111111111111112",  // SOL
          "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"  // USDC
        ],
        "exclude_tokens": []
      }
    },
    "DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF": {
      "scan": {
        "scan_mode": "blacklist",    // This wallet uses blacklist mode
        "include_tokens": [],
        "exclude_tokens": [
          "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"   // BONK
        ]
      }
    }
  }
}
```

In this example:
- The global scan mode is "all"
- The first wallet only monitors SOL and USDC
- The second wallet monitors all tokens except BONK
- Any other wallets use the default "all" mode

## Use Cases

### 1. Focus on Important Tokens

Whitelist mode is perfect for treasury wallets or investment accounts where you only care about major tokens:

```json
{
  "scan": {
    "scan_mode": "whitelist",
    "include_tokens": [
      "So11111111111111111111111111111111111111112",  // SOL
      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",  // USDC
      "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"   // USDT
    ],
    "exclude_tokens": []
  }
}
```

### 2. Filtering Out Noise

Many wallets receive airdrops or dust tokens that create noise. Blacklist mode helps filter these out:

```json
{
  "scan": {
    "scan_mode": "blacklist",
    "include_tokens": [],
    "exclude_tokens": [
      "TokenAddress1",  // Airdrop 1
      "TokenAddress2",  // Airdrop 2
      "TokenAddress3"   // Dust token
    ]
  }
}
```

### 3. Different Strategies for Different Wallets

Configure scan modes based on each wallet's purpose:

```json
{
  "wallet_configs": {
    "TreasuryWallet": {
      "scan": {
        "scan_mode": "whitelist",
        "include_tokens": ["SOL", "USDC", "USDT"]
      }
    },
    "TradingWallet": {
      "scan": {
        "scan_mode": "all"           // Monitor everything in trading wallet
      }
    },
    "PersonalWallet": {
      "scan": {
        "scan_mode": "blacklist",
        "exclude_tokens": ["Dust1", "Dust2", "Airdrop1"]
      }
    }
  }
}
```

## Alert Filtering vs. Scan Filtering

Note the difference between scan configuration and alert filtering:

1. **Scan Configuration**: Controls which tokens are monitored at all. Tokens filtered out here won't appear in the monitor's data.

2. **Alert Filtering**: Controls which tokens generate alerts when they change. Tokens in the ignore list are still monitored but don't trigger alerts:

```json
{
  "alerts": {
    "ignore_tokens": ["TokenAddress1", "TokenAddress2"]
  }
}
```

The main difference:
- Tokens excluded by scan configuration aren't monitored or stored
- Tokens in the alert ignore list are monitored and stored but don't trigger alerts

## Performance Considerations

Using whitelist mode can improve performance, especially for wallets with many tokens, as it reduces the amount of data that needs to be processed and stored. This can be helpful when:

- Monitoring wallets with hundreds of tokens
- Running on systems with limited resources
- Reducing RPC usage on public endpoints

## Examples

See [examples/scan_config_examples.json](https://github.com/accursedgalaxy/insider-monitor/blob/main/examples/scan_config_examples.json) for a complete example configuration.
