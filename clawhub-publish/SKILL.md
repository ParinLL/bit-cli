---
name: bit
description: Explain bit-cli skill purpose, installation, required setup, and troubleshooting.
metadata: {"openclaw": {"requires": {"bins": ["bit", "git", "go", "sudo"], "env": ["BIT_API_KEY"]}, "primaryEnv": "BIT_API_KEY", "install": [{"id": "go-install", "kind": "go", "label": "Install bit via Go", "bins": ["bit"], "module": "github.com/ParinLL/bit-cli"}]}}
---

# Bit CLI Skill (Documentation-Only)

## Skill Purpose And Trigger Scenarios

- Purpose: Provide a usage entry point for the Bit URL Shortener CLI (`bit`) to create, query, update, delete short links, and view click data.
- Trigger scenarios:
- The user mentions needs like "short URL", "bit-cli", or commands such as `bit create/list/get/update/delete/clicks`.
- The user wants to run Bit API operations with OpenClaw.
- The user needs to verify Bit API availability (for example, a health check).

## Installation Commands (Or GitHub Link To Installation Section)

- Installation mechanism: This skill provides an OpenClaw `install spec` (`kind: go`) and also keeps manual build instructions as a fallback.
- Risk note: Building arbitrary source code carries supply-chain risk. Review repository contents and version sources before running install steps.
- GitHub (this project): `https://github.com/ParinLL/bit-cli`
- README install section: `https://github.com/ParinLL/bit-cli#build-locally`
- Safety reminder: Review source repository contents and trust level before building from source.
- Installation commands:

```bash
git clone https://github.com/ParinLL/bit-cli.git
cd bit-cli
go build -o bit .
sudo mv bit /usr/local/bin/
```

- Alternative without `sudo` (install to user directory):

```bash
git clone https://github.com/ParinLL/bit-cli.git
cd bit-cli
go build -o bit .
mkdir -p "$HOME/.local/bin"
mv bit "$HOME/.local/bin/"
export PATH="$HOME/.local/bin:$PATH"
```

## Required Environment Variables / Permissions

- Required environment variables:
- `BIT_API_KEY` (required): Bit API authentication key.
- `BIT_API_URL` (optional): Bit API base URL, default `http://localhost:4000`.
- Permission requirements:
- The `bit` executable must be callable from PATH.
- Installing to `/usr/local/bin` with `sudo mv` requires administrator privileges.
- If the target API is remote, network connectivity to that API is required.

## Common Troubleshooting

- `bit: command not found`
- Cause: The CLI is not installed or not in PATH.
- Fix: Rebuild with `go build` and verify `which bit` returns a valid path.
- `401 Unauthorized` / `403 Forbidden`
- Cause: `BIT_API_KEY` is missing or invalid.
- Fix: Reset `BIT_API_KEY` and confirm the key is still valid on the server.
- `connection refused` / timeout
- Cause: `BIT_API_URL` is incorrect, the Bit service is not running, or the network is unreachable.
- Fix: Run `bit ping` first, then verify API service status and URL.
- Command succeeds but data is unexpected
- Cause: The target ID does not exist, data was deleted, or the update payload format is incorrect.
- Fix: Validate current state with `bit list` or `bit get <id>`, then retry.
