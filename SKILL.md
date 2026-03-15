---
name: bit
description: Manage Bit URL shortener links — create, list, get, update, delete short links and view click analytics.
metadata: {"openclaw": {"requires": {"bins": ["bit"], "env": ["BIT_API_KEY"]}, "primaryEnv": "BIT_API_KEY"}}
---

# Bit URL Shortener

Use this skill to interact with the Bit URL shortener API via the `bit` CLI.

## Available Commands

- `bit ping` — Check API health
- `bit list [--limit N] [--cursor X]` — List all short links
- `bit create <url>` — Create a new short link
- `bit get <id>` — Get link details and recent click analytics
- `bit update <id> <new-url>` — Update a link's destination
- `bit delete <id>` — Delete a link
- `bit clicks <id> [--limit N] [--cursor X]` — List click records

## Environment Variables

- `BIT_API_KEY` (required) — API authentication key
- `BIT_API_URL` (optional) — API base URL, defaults to `http://localhost:4000`

## Original GitHub Sources

- This CLI repository: `https://github.com/ParinLL/bit-cli`
- Bit server (upstream project): `https://github.com/sjdonado/bit`
- Bit API docs: `https://sjdonado.github.io/bit/`

## Installation

### Build Locally

```bash
git clone git@github.com:ParinLL/bit-cli.git
cd bit-cli
go build -o bit .
```

### Install to PATH

```bash
cd bit-cli
go build -o bit .
sudo mv bit /usr/local/bin/
```

### Docker Image

```bash
docker build -t bit-cli .
docker buildx build --platform linux/amd64,linux/arm64 -t bit-cli .
```

### Docker Compose

```bash
export BIT_ADMIN_API_KEY=$(openssl rand -base64 32)
docker compose up -d bit
docker compose run --rm bit-cli ping
```

## Usage Examples

When the user asks to shorten a URL, run: `bit create <url>`
When the user asks to see all links, run: `bit list`
When the user asks about clicks on a link, run: `bit clicks <id>`
