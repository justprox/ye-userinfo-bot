# Telegram User & Chat ID Bot

[![Telegram Bot](https://img.shields.io/badge/Telegram-@ye__info__bot-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/ye_info_bot)

> 🤖 **Official Telegram Bot:** [@ye_info_bot](https://t.me/ye_info_bot)

A lightweight, ultra-fast, and privacy-first Telegram bot written in Go for discovering Telegram IDs of users, groups, supergroups, channels, forum topics, and forwarded message authors. Deployed serverless on **Vercel** with zero logging, sub-50ms response times, and scanner protection.

---

## ✨ Features / Возможности

- **👤 Personal Profile ID**:
  - Send any message or `/start` to get your `User ID`, Full Name, `@username`, `Language Code`, and `Telegram Premium` badge.
- **📨 Forwarded Message Author ID (`forward_origin`)**:
  - **From Users & Bots**: Extracts original author's `User ID`, Name, and Username.
  - **From Hidden Profiles**: Detects when an author has enabled Telegram forward privacy settings and explains how to obtain their ID.
  - **From Channels**: Returns `Channel ID`, Channel Title, Username, and original Post ID.
  - **From Groups**: Returns `Chat ID` and Title.
- **📇 Contact Card Sharing**:
  - Share any contact from your phonebook (📎 Attachment → Contact) to resolve their Telegram `User ID`.
- **👥 Group & Forum Topic Support**:
  - Send `/id` in any group to get `Group Chat ID`, Chat Type, and Forum `MessageThreadID`.
  - Reply with `/id` to any member's message in a group to fetch that specific user's info.
  - Send `/quit` (or `/leave`) to make the bot leave the group gracefully.
- **🌐 Multi-Language Support (i18n)**:
  - Automatic language detection based on user's Telegram client locale.
  - Fully localized in **English** (default) and **Russian** (`ru`).
- **🔒 Zero-Logging & Security-by-Default**:
  - Zero logging of incoming updates, user IDs, message texts, or personal data.
  - All outputs dynamically sanitized with `html.EscapeString` against HTML injection.
  - Completely stateless: 0 databases, 0 persistent storage.
- **🛡️ Scanner & Probe Protection**:
  - Configurable secret webhook path (`WEBHOOK_PATH`) to return `404 Not Found` to any unauthorized crawlers or scanners.
  - Optional `WEBHOOK_SECRET` token validation via `X-Telegram-Bot-Api-Secret-Token`.
- **🔨 Verifiable Transparency**:
  - `/about` command shows the exact Git revision, commit timestamp, Go compiler version, binary SHA-256 checksum, and repository link.
  - Bit-for-bit reproducible builds verifiable against GitHub Actions artifacts and releases.

---

## 💡 How IDs Work in Telegram / Как устроены ID

| Entity Type / Тип | ID Format / Формат | Example / Пример | Description / Описание |
| :--- | :--- | :--- | :--- |
| **User / Bot** | Positive Integer | `123456789` | Unique immutable account identifier / Уникальный ID аккаунта. |
| **Basic Group** | Negative Integer | `-987654321` | Legacy small group chats / Обычные небольшие группы. |
| **Supergroup** | Negative 13-digit (`-100...`) | `-1001234567890` | Modern group chats and forum supergroups / Супергруппы. |
| **Channel** | Negative 13-digit (`-100...`) | `-1001987654321` | Public and private channels / Каналы. |
| **Forum Topic** | Positive Integer (`MessageThreadID`) | `42` | Initial message ID creating the forum thread / ID ветки форума. |

---

## 🤖 Commands / Команды

| Command | Scope | Description (EN) | Описание (RU) |
| :--- | :--- | :--- | :--- |
| `/start` | Private / Groups | Start the bot, display help and personal profile ID | Запустить бота и узнать свой ID |
| `/id` | Private / Groups | Get current chat ID, thread ID, or replied user's info | Узнать ID чата, топика или пользователя |
| `/about` | Private / Groups | Display build transparency, Git revision, and SHA-256 | О боте, Git ревизия и верификация |
| `/quit` | Groups | Make the bot leave the group chat | Покинуть группу |
| `/help` | Private / Groups | Help guide and Telegram ID architecture overview | Справка об устройстве Telegram ID |

---

## 🚀 Serverless Deployment on Vercel

Deploy your own instance of the bot for free in less than 1 minute:

1. **Import Repository**:
   Import this repository in your [Vercel Dashboard](https://vercel.com/dashboard).

2. **Add Environment Variables**:
   In Project Settings → **Environment Variables** (`Production`):
   - `BOT_TOKEN`: `your-telegram-bot-token-from-botfather`
   - `WEBHOOK_PATH` *(recommended)*: `/webhook_<random-hex>`
   - `WEBHOOK_SECRET` *(optional)*: `<random-secret-token>`

3. **Deploy & Set Webhook**:
   After deployment, register your Vercel webhook with Telegram:
   ```bash
   curl "https://api.telegram.org/bot<BOT_TOKEN>/setWebhook?url=https://<your-vercel-domain>.vercel.app<WEBHOOK_PATH>&secret_token=<WEBHOOK_SECRET>"
   ```

---

## 🔨 Reproducible Builds & Verification

Every commit generates bit-for-bit reproducible static binaries for **Linux x86_64 (`linux/amd64`)** via GitHub Actions.

### Verify Build Independently

To compile the exact same binary locally and verify the SHA-256 hash:

```bash
git checkout <commit-hash>
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" .
sha256sum ye-userinfo-bot
```

Compare the output hash with the **Binary SHA-256** reported by `/about` in the bot or in `checksums.txt` on the GitHub Releases page.

---

## 🛡️ Security & Privacy

- **Zero Logging**: Incoming messages, user identities, IPs, and queries are never logged to console, disk, or external services.
- **No Data Persistence**: The bot operates completely stateless and stores zero database entries.
- **Safe HTML Escaping**: Dynamic user names and titles are safely escaped before rendering.
- **Fail-Safe Panic Recovery**: Centralized recovery prevents crashes while keeping errors silent and safe.

---

## 📄 License

MIT License. Open source and free to use.
