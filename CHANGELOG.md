# Changelog

## [Unreleased]

### Added
- Token filtering functionality
  - New `scan` configuration section in `config.json`
  - Three scanning modes:
    - `all`: Monitor all tokens (default)
    - `whitelist`: Only monitor specific tokens listed in `include_tokens`
    - `blacklist`: Monitor all tokens except those listed in `exclude_tokens`
  - Example configuration:
    ```json
    "scan": {
        "scan_mode": "whitelist",
        "include_tokens": ["token1", "token2"],
        "exclude_tokens": []
    }
    ```

## Usage
To filter tokens, update your `config.json` with the new `scan` section:

1. To monitor all tokens (default behavior):
```json
"scan": {
    "scan_mode": "all",
    "include_tokens": [],
    "exclude_tokens": []
}
```

2. To monitor only specific tokens:
```json
"scan": {
    "scan_mode": "whitelist",
    "include_tokens": [
        "token_address_1",
        "token_address_2"
    ],
    "exclude_tokens": []
}
```

3. To exclude specific tokens:
```json
"scan": {
    "scan_mode": "blacklist",
    "include_tokens": [],
    "exclude_tokens": [
        "token_address_to_ignore_1",
        "token_address_to_ignore_2"
    ]
}
```

### Notes
- Token addresses should be in Solana's base58 format
- Changes take effect immediately after updating the configuration
- Invalid token addresses will be logged but won't cause the monitor to fail 