# Installation Guide

This guide walks you through the process of installing the Solana Insider Monitor on your system.

## Prerequisites

Before installing Solana Insider Monitor, ensure you have the following:

- **Go** version 1.23.2 or later
- **Git** for cloning the repository
- **Access to a Solana RPC endpoint** (mainnet, devnet, or testnet)

!!! tip "RPC Endpoints"
    You can use public RPC endpoints, but for better reliability and performance, consider using a dedicated service like:

    - [QuickNode](https://www.quicknode.com/chains/sol)
    - [Alchemy](https://www.alchemy.com/solana)
    - [Ankr](https://www.ankr.com/rpc/solana/)

## Installation Methods

=== "From Source"

    ### Clone the Repository

    ```bash
    git clone https://github.com/accursedgalaxy/insider-monitor.git
    cd insider-monitor
    ```

    ### Install Dependencies

    ```bash
    go mod download
    ```

    ### Build the Binary (Optional)

    ```bash
    make build
    ```

    This will compile the binary to the `bin` directory.

=== "Using Docker"

    ### Pull the Docker Image

    ```bash
    docker pull accursedgalaxy/insider-monitor:latest
    ```

    ### Run with Docker

    ```bash
    docker run -d \
      --name solana-monitor \
      -p 8080:8080 \
      -v $(pwd)/config.json:/app/config.json \
      -v $(pwd)/data:/app/data \
      accursedgalaxy/insider-monitor
    ```

    ### Using Docker Compose

    Create a `docker-compose.yml` file:

    ```yaml
    version: '3'
    services:
      monitor:
        image: accursedgalaxy/insider-monitor:latest
        ports:
          - "8080:8080"
        volumes:
          - ./config.json:/app/config.json
          - ./data:/app/data
        restart: unless-stopped
    ```

    Then run:

    ```bash
    docker-compose up -d
    ```

=== "Prebuilt Binaries"

    ### Download Latest Release

    Go to the [Releases page](https://github.com/accursedgalaxy/insider-monitor/releases) and download the appropriate binary for your operating system.

    ### Extract and Install

    ```bash
    # Extract the archive
    tar -xzf insider-monitor-v1.0.0-linux-amd64.tar.gz

    # Make the binary executable
    chmod +x insider-monitor

    # Optionally, move to a directory in your PATH
    sudo mv insider-monitor /usr/local/bin/
    ```

## Verifying Installation

After installation, verify it works by running:

```bash
insider-monitor -version
```

Or if running from source:

```bash
go run cmd/monitor/main.go -version
```

You should see output showing the current version of Solana Insider Monitor.

## Initial Configuration

After installation, you'll need to create a configuration file:

1. Copy the example configuration:
   ```bash
   cp config.example.json config.json
   ```

2. Edit the file to include your Solana wallet addresses and network details.

See the [Configuration Guide](../configuration/index.md) for detailed information on configuration options.

## Next Steps

- [Quick Start Guide](quick-start.md) - Get up and running quickly
- [Configuration Guide](../configuration/index.md) - Learn how to configure the monitor
- [Running Modes](running-modes.md) - Understand different ways to run the monitor
