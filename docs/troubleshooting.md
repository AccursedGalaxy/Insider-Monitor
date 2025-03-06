# Troubleshooting

This guide helps you diagnose and fix common issues that you might encounter when using Solana Insider Monitor.

## Diagnostic Steps

Before diving into specific issues, try these general diagnostic steps:

1. **Check Logs**: Run the monitor with verbose logging enabled:
   ```bash
   RUST_LOG=debug insider-monitor
   ```

2. **Verify Configuration**: Ensure your `config.json` file is correctly formatted and contains valid values:
   ```bash
   cat config.json | jq
   ```

3. **Check Connectivity**: Verify that you can connect to the Solana RPC endpoint:
   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","id":1,"method":"getHealth"}' https://api.mainnet-beta.solana.com
   ```

4. **Restart the Monitor**: Sometimes simply restarting the monitor resolves issues.

## Common Issues and Solutions

### Connection Issues

#### "Failed to connect to RPC endpoint"

**Symptoms:**
- Error messages about connection failures
- No wallet data is retrieved
- Log messages indicating RPC timeouts

**Possible causes:**
1. Invalid RPC URL in configuration
2. Network connectivity problems
3. RPC endpoint is down or under maintenance
4. Rate limiting by the RPC provider

**Solutions:**

1. **Verify your RPC URL**:
   Check that your `network_url` setting is correct:
   ```json
   {
     "network_url": "https://api.mainnet-beta.solana.com"
   }
   ```

2. **Try an alternative RPC endpoint**:
   ```json
   {
     "network_url": "https://solana-api.projectserum.com"
   }
   ```

3. **Check your network connection**:
   ```bash
   ping api.mainnet-beta.solana.com
   ```

4. **Increase retry settings**:
   If you're experiencing intermittent connectivity issues, adjust your scan interval:
   ```json
   {
     "scan_interval": "5m"
   }
   ```

### Alert Configuration Issues

#### "Not receiving alerts for significant changes"

**Symptoms:**
- You can see balance changes in the web interface
- No alerts are being sent to Discord or shown in the console

**Possible causes:**
1. Alert thresholds are set too high
2. Discord integration is misconfigured
3. Token is in the ignore list
4. Balance is below minimum threshold

**Solutions:**

1. **Check alert thresholds**:
   Lower the `significant_change` value to catch smaller changes:
   ```json
   {
     "alerts": {
       "significant_change": 0.05
     }
   }
   ```

2. **Verify Discord configuration**:
   Ensure Discord integration is properly configured:
   ```json
   {
     "discord": {
       "enabled": true,
       "webhook_url": "https://discord.com/api/webhooks/your-webhook-url",
       "channel_id": "your-channel-id"
     }
   }
   ```

3. **Check the token ignore list**:
   Remove any tokens that you want to monitor from the ignore list:
   ```json
   {
     "alerts": {
       "ignore_tokens": []
     }
   }
   ```

4. **Lower minimum balance threshold**:
   Adjust the minimum balance to catch changes in smaller token holdings:
   ```json
   {
     "alerts": {
       "minimum_balance": 100
     }
   }
   ```

### Web Interface Issues

#### "Cannot access web interface"

**Symptoms:**
- Web interface is not loading at http://localhost:8080
- Connection refused errors

**Possible causes:**
1. Monitor not running with web mode enabled
2. Firewall blocking the port
3. Port conflict with another application
4. Permission issues

**Solutions:**

1. **Start with web flag enabled**:
   ```bash
   insider-monitor -web
   ```

2. **Check for port conflicts**:
   ```bash
   lsof -i :8080
   ```
   If another application is using the port, use a different port:
   ```bash
   insider-monitor -web -port 9090
   ```

3. **Check firewall settings**:
   ```bash
   sudo ufw status
   # For Ubuntu/Debian

   # or

   sudo firewall-cmd --list-all
   # For Fedora/CentOS
   ```

4. **Run with elevated permissions if needed**:
   ```bash
   sudo insider-monitor -web
   ```

#### "Authentication issues with web interface"

**Symptoms:**
- Unable to log in to the web interface
- "Invalid credentials" error

**Solutions:**

1. **Use default credentials**:
   Username: `admin`
   Password: `admin`

2. **Set a custom password**:
   ```bash
   export ADMIN_PASSWORD="your-secure-password"
   insider-monitor -web
   ```

3. **Reset authentication**:
   Stop the monitor, delete any existing auth files, and restart:
   ```bash
   rm -f data/auth.json
   insider-monitor -web
   ```

### Data Storage Issues

#### "Missing or corrupt data after restart"

**Symptoms:**
- No historical data available after restarting the monitor
- Errors about corrupted data files

**Possible causes:**
1. Permission issues with data directory
2. Disk space issues
3. Corrupted data files

**Solutions:**

1. **Check permissions**:
   ```bash
   ls -la data/
   chmod -R 755 data/
   ```

2. **Check disk space**:
   ```bash
   df -h
   ```

3. **Backup and reset data**:
   ```bash
   cp -r data data_backup
   rm -f data/wallet_data.json
   insider-monitor
   ```

### Performance Issues

#### "High CPU or memory usage"

**Symptoms:**
- Monitor process using excessive CPU or memory
- System becoming slow or unresponsive

**Possible causes:**
1. Too many wallets being monitored
2. Scan interval too short
3. Memory leak in application

**Solutions:**

1. **Reduce number of monitored wallets**:
   Start with fewer wallets and gradually add more.

2. **Increase scan interval**:
   ```json
   {
     "scan_interval": "15m"
   }
   ```

3. **Limit token scan depth**:
   ```json
   {
     "advanced": {
       "max_tokens_per_wallet": 100
     }
   }
   ```

4. **Update to latest version**:
   ```bash
   git pull
   make build
   ```

## Debugging Tools

### Log Analysis

Run with debug logging and save to a file for analysis:

```bash
RUST_LOG=debug insider-monitor > monitor.log 2>&1
```

Then analyze the log file:

```bash
grep ERROR monitor.log
grep "wallet scan" monitor.log
```

### Network Debugging

Test connectivity to the Solana RPC endpoint:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","id":1,"method":"getHealth"}' https://api.mainnet-beta.solana.com
```

### Configuration Validation

Validate your configuration file:

```bash
cat config.json | jq
```

## Getting Help

If you're still experiencing issues after trying these troubleshooting steps:

1. **Check GitHub Issues**:
   Search the [GitHub Issues](https://github.com/accursedgalaxy/insider-monitor/issues) to see if others have reported similar problems.

2. **Discord Community**:
   Join our [Discord server](https://discord.gg/7vY9ZBPdya) to get help from the community and developers.

3. **Create a Bug Report**:
   If you've found a bug, please [create an issue](https://github.com/accursedgalaxy/insider-monitor/issues/new) with:
   - Detailed description of the problem
   - Steps to reproduce
   - Your configuration (with sensitive information removed)
   - Log output
   - System information (OS, Go version)
