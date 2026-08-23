# AGENTS.md

## Project Overview

`ye-userinfo-bot` is a lightweight, ultra-fast, and privacy-first Telegram bot written in Go. It discovers Telegram IDs for users, groups, supergroups, channels, forum topics, and forwarded message authors. Deployed serverless on **Vercel** with zero logging and reproducible builds.

- **Telegram Bot Username:** [@ye_info_bot](https://t.me/ye_info_bot)
- **Repository:** `https://github.com/justprox/ye-userinfo-bot`
- **Go Version:** Go 1.23+

---

## Architecture & Codebase Layout

```text
ye-userinfo-bot/
├── .github/
│   └── workflows/
│       └── release.yml     # CI test suite & reproducible Linux x64 binary builds
├── pkg/
│   └── bot/
│       ├── bot.go          # Bot constructor, client initialization, menu setup
│       ├── formatter.go    # HTML formatting helpers, i18n localization (EN/RU)
│       ├── handlers.go     # Safe update dispatching, private/group/forward handlers
│       └── version.go      # VCS transparency info, commit hash, build checksums
├── .env.example            # Template for environment variables
├── .gitignore              # Ignores .env*, IDE files, binaries, .vercel
├── AGENTS.md               # Agent and contributor technical documentation
├── go.mod                  # Go module definition (ye-userinfo-bot)
├── go.sum                  # Checksums of dependencies
├── main.go                 # Vercel Webhook HTTP server entrypoint
└── README.md               # User-facing guide, setup, and ID reference
```

---

## Key Design Principles & Security Guardrails

1. **Zero-Logging Enforced**:
   - Never log incoming updates, message contents, chat IDs, or personal user data to stdout, stderr, disk, or remote logging services.
   - Keep bot errors silent or sanitized.

2. **Synchronous Webhook Execution for Serverless**:
   - In serverless environments (e.g. Vercel), execution freezes immediately after HTTP `200 OK` is returned.
   - All handlers must execute synchronously via `bot.HandleUpdate(context.Background(), b, &update)` so outbound Telegram API calls finish before sending the HTTP response.

3. **Scanner Protection**:
   - Support `WEBHOOK_PATH` (e.g. `/webhook_<random-hex>`) to return `404 Not Found` for any unauthorized root `/` or scanner requests.
   - Support `WEBHOOK_SECRET` to verify the `X-Telegram-Bot-Api-Secret-Token` header.

4. **Safe HTML Formatting**:
   - All dynamic strings from users (names, titles, signatures) must pass through `escape()` (`html.EscapeString`) before embedding into HTML messages.

5. **Panic Recovery**:
   - Top-level recovery in `HandleUpdate` ensures no panics crash the HTTP server.

---

## Local Development & Testing

### Prerequisites
- Go 1.23 or newer.

### Commands
```bash
# Verify dependencies
go mod verify

# Run tests
go test -v ./...

# Build binary locally
go build -trimpath -ldflags="-s -w -buildid=" -o ye-userinfo-bot .

# Run local HTTP server
BOT_TOKEN="your_token" go run .
```

---

## CI/CD & Reproducible Builds

The GitHub Actions workflow (`.github/workflows/release.yml`) builds bit-for-bit reproducible binaries for `linux/amd64` using:
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -trimpath \
  -ldflags="-s -w -buildid= -X ye-userinfo-bot/pkg/bot.GitCommit=${COMMIT} -X ye-userinfo-bot/pkg/bot.BuildTime=${BUILD_TIME} -X ye-userinfo-bot/pkg/bot.RepoURL=${REPO}" \
  -o dist/ye-userinfo-bot-linux-amd64 .
```
SHA-256 hashes are automatically generated in `checksums.txt` and verifiable against `/about` command output.
