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

## Usage Examples

When the user asks to shorten a URL, run: `bit create <url>`
When the user asks to see all links, run: `bit list`
When the user asks about clicks on a link, run: `bit clicks <id>`
