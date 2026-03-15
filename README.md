# bit-cli

CLI tool for [Bit](https://github.com/sjdonado/bit) URL Shortener, written in Go.

## OpenClaw / ClawHub

This skill can be used in OpenClaw.

- ClawHub: https://clawhub.ai/ParinLL/bit

## Features

- `ping` — Health check
- `list` — List all short links
- `create <url>` — Create a short link
- `get <id>` — Get link details with recent click records
- `update <id> <url>` — Update a link's URL
- `delete <id>` — Delete a link
- `clicks <id>` — List click records for a link

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BIT_API_URL` | Bit API base URL | `http://localhost:4000` |
| `BIT_API_KEY` | API authentication key | (empty) |

## Build Locally

```bash
cd bit-cli
go build -o bit .
```

## Install to PATH

```bash
cd bit-cli
go build -o bit .
sudo mv bit /usr/local/bin/
```

## Usage

```bash
# Set environment variables
export BIT_API_URL="http://localhost:4000"
export BIT_API_KEY="your-api-key"

# Health check
bit ping

# Create a short link
bit create https://example.com

# List all links
bit list
bit list --limit 10 --cursor 5

# Get link details
bit get 1

# Update a link
bit update 1 https://new-url.com

# Delete a link
bit delete 1

# View click records
bit clicks 1
bit clicks 1 --limit 50 --cursor 10
```

## Docker

### Build Image (Multi-arch)

```bash
# Single platform
docker build -t bit-cli .

# Multi-platform (requires buildx)
docker buildx build --platform linux/amd64,linux/arm64 -t bit-cli .
```

### Run Directly

```bash
docker run --rm \
  -e BIT_API_URL="http://host.docker.internal:4000" \
  -e BIT_API_KEY="your-api-key" \
  bit-cli create https://example.com
```

## Docker Compose

The docker-compose.yaml includes both the Bit server and CLI services.

### Start Bit Server

```bash
# Generate an API key
export BIT_ADMIN_API_KEY=$(openssl rand -base64 32)

cd bit-cli
docker compose up -d bit
```

### Create a User to Get an API Key

```bash
docker compose exec bit cli --create-user=Admin
```

### Use the CLI

```bash
export BIT_API_KEY="your-api-key"

# Run CLI commands via docker compose
docker compose run --rm bit-cli ping
docker compose run --rm bit-cli create https://example.com
docker compose run --rm bit-cli list
```

## References

- Bit API docs: https://sjdonado.github.io/bit/
- Bit source code: https://github.com/sjdonado/bit
