package bot

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/go-telegram/bot/models"
)

// Lang represents supported languages
type Lang string

const (
	LangEN Lang = "en"
	LangRU Lang = "ru"
)

// detectLang determines the user language based on Telegram language_code.
// Defaults to English (en).
func detectLang(u *models.User) Lang {
	if u == nil || u.LanguageCode == "" {
		return LangEN
	}
	code := strings.ToLower(strings.TrimSpace(u.LanguageCode))
	if strings.HasPrefix(code, "ru") {
		return LangRU
	}
	return LangEN
}

// escape safely escapes dynamic user input to prevent HTML formatting issues.
func escape(s string) string {
	return html.EscapeString(s)
}

// formatUser returns formatted HTML info about a Telegram user.
func formatUser(u *models.User, header string, lang Lang) string {
	if u == nil {
		return ""
	}

	var sb strings.Builder
	if header != "" {
		sb.WriteString(fmt.Sprintf("<b>%s</b>\n", escape(header)))
	}

	sb.WriteString(fmt.Sprintf("🆔 <b>ID:</b> <code>%d</code>\n", u.ID))

	name := u.FirstName
	if u.LastName != "" {
		name += " " + u.LastName
	}

	nameLabel := "Name:"
	langLabel := "Language:"
	premiumLabel := "⭐ <b>Telegram Premium:</b> Yes\n"
	botLabel := "🤖 <b>Bot:</b> Yes\n"

	if lang == LangRU {
		nameLabel = "Имя:"
		langLabel = "Язык:"
		premiumLabel = "⭐ <b>Telegram Premium:</b> Да\n"
		botLabel = "🤖 <b>Бот:</b> Да\n"
	}

	sb.WriteString(fmt.Sprintf("👤 <b>%s</b> %s\n", nameLabel, escape(name)))

	if u.Username != "" {
		sb.WriteString(fmt.Sprintf("🔗 <b>Username:</b> @%s\n", escape(u.Username)))
	}

	if u.LanguageCode != "" {
		sb.WriteString(fmt.Sprintf("🌐 <b>%s</b> <code>%s</code>\n", langLabel, escape(u.LanguageCode)))
	}

	if u.IsPremium {
		sb.WriteString(premiumLabel)
	}

	if u.IsBot {
		sb.WriteString(botLabel)
	}

	return sb.String()
}

// formatChat returns formatted HTML info about a Telegram chat (group, supergroup, channel, private).
func formatChat(c *models.Chat, header string, lang Lang) string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	if header != "" {
		sb.WriteString(fmt.Sprintf("<b>%s</b>\n", escape(header)))
	}

	chatIDLabel := "Chat ID:"
	typeLabel := "Type:"
	titleLabel := "Title:"
	forumLabel := "📑 <b>Forum (Topics):</b> Enabled\n"

	if lang == LangRU {
		chatIDLabel = "ID чата:"
		typeLabel = "Тип:"
		titleLabel = "Название:"
		forumLabel = "📑 <b>Форум (Темы/Топики):</b> Включены\n"
	}

	sb.WriteString(fmt.Sprintf("🆔 <b>%s</b> <code>%d</code>\n", chatIDLabel, c.ID))

	chatTypeDesc := ""
	switch c.Type {
	case models.ChatTypePrivate:
		if lang == LangRU {
			chatTypeDesc = "Личный диалог (Private)"
		} else {
			chatTypeDesc = "Private Chat"
		}
	case models.ChatTypeGroup:
		if lang == LangRU {
			chatTypeDesc = "Обычная группа (Group, ID начинается с -)"
		} else {
			chatTypeDesc = "Basic Group (ID starts with -)"
		}
	case models.ChatTypeSupergroup:
		if lang == LangRU {
			chatTypeDesc = "Супергруппа (Supergroup, ID начинается с -100)"
		} else {
			chatTypeDesc = "Supergroup (ID starts with -100)"
		}
	case models.ChatTypeChannel:
		if lang == LangRU {
			chatTypeDesc = "Канал (Channel, ID начинается с -100)"
		} else {
			chatTypeDesc = "Channel (ID starts with -100)"
		}
	default:
		chatTypeDesc = string(c.Type)
	}

	sb.WriteString(fmt.Sprintf("📁 <b>%s</b> %s\n", typeLabel, escape(chatTypeDesc)))

	if c.Title != "" {
		sb.WriteString(fmt.Sprintf("🏷 <b>%s</b> %s\n", titleLabel, escape(c.Title)))
	}

	if c.Username != "" {
		sb.WriteString(fmt.Sprintf("🔗 <b>Username:</b> @%s\n", escape(c.Username)))
	}

	if c.IsForum {
		sb.WriteString(forumLabel)
	}

	return sb.String()
}

// formatForwardOrigin formats information about where the message was forwarded from.
func formatForwardOrigin(fo *models.MessageOrigin, lang Lang) string {
	if fo == nil {
		return ""
	}

	var sb strings.Builder
	if lang == LangRU {
		sb.WriteString("<b>📨 Информация о пересланном сообщении:</b>\n")
	} else {
		sb.WriteString("<b>📨 Forwarded Message Information:</b>\n")
	}

	dateLabel := "Original Date:"
	sigLabel := "Signature:"
	if lang == LangRU {
		dateLabel = "Дата оригинала:"
		sigLabel = "Подпись:"
	}

	switch fo.Type {
	case models.MessageOriginTypeUser:
		if fo.MessageOriginUser != nil {
			if lang == LangRU {
				sb.WriteString("📌 <b>Источник:</b> Пользователь / Бот\n")
			} else {
				sb.WriteString("📌 <b>Source:</b> User / Bot\n")
			}
			sb.WriteString(formatUser(&fo.MessageOriginUser.SenderUser, "", lang))
			if fo.MessageOriginUser.Date > 0 {
				t := time.Unix(int64(fo.MessageOriginUser.Date), 0).UTC()
				sb.WriteString(fmt.Sprintf("🕒 <b>%s</b> %s UTC\n", dateLabel, t.Format("2006-01-02 15:04:05")))
			}
		}

	case models.MessageOriginTypeHiddenUser:
		if fo.MessageOriginHiddenUser != nil {
			if lang == LangRU {
				sb.WriteString("📌 <b>Источник:</b> Скрытый автор (приватный профиль)\n")
				sb.WriteString(fmt.Sprintf("👤 <b>Имя автора:</b> %s\n", escape(fo.MessageOriginHiddenUser.SenderUserName)))
				sb.WriteString("🔒 <i>ID скрыт настройками конфиденциальности пересылки автора.</i>\n")
			} else {
				sb.WriteString("📌 <b>Source:</b> Hidden User (Private Profile)\n")
				sb.WriteString(fmt.Sprintf("👤 <b>Author Name:</b> %s\n", escape(fo.MessageOriginHiddenUser.SenderUserName)))
				sb.WriteString("🔒 <i>ID is hidden by the author's forward privacy settings.</i>\n")
			}
			if fo.MessageOriginHiddenUser.Date > 0 {
				t := time.Unix(int64(fo.MessageOriginHiddenUser.Date), 0).UTC()
				sb.WriteString(fmt.Sprintf("🕒 <b>%s</b> %s UTC\n", dateLabel, t.Format("2006-01-02 15:04:05")))
			}
		}

	case models.MessageOriginTypeChat:
		if fo.MessageOriginChat != nil {
			if lang == LangRU {
				sb.WriteString("📌 <b>Источник:</b> Группа / Чат\n")
			} else {
				sb.WriteString("📌 <b>Source:</b> Group / Chat\n")
			}
			sb.WriteString(formatChat(&fo.MessageOriginChat.SenderChat, "", lang))
			if fo.MessageOriginChat.AuthorSignature != nil && *fo.MessageOriginChat.AuthorSignature != "" {
				sb.WriteString(fmt.Sprintf("✍️ <b>%s</b> %s\n", sigLabel, escape(*fo.MessageOriginChat.AuthorSignature)))
			}
			if fo.MessageOriginChat.Date > 0 {
				t := time.Unix(int64(fo.MessageOriginChat.Date), 0).UTC()
				sb.WriteString(fmt.Sprintf("🕒 <b>%s</b> %s UTC\n", dateLabel, t.Format("2006-01-02 15:04:05")))
			}
		}

	case models.MessageOriginTypeChannel:
		if fo.MessageOriginChannel != nil {
			if lang == LangRU {
				sb.WriteString("📌 <b>Источник:</b> Канал\n")
			} else {
				sb.WriteString("📌 <b>Source:</b> Channel\n")
			}
			sb.WriteString(formatChat(&fo.MessageOriginChannel.Chat, "", lang))
			if fo.MessageOriginChannel.MessageID > 0 {
				postIDLabel := "Channel Post ID:"
				if lang == LangRU {
					postIDLabel = "ID сообщения в канале:"
				}
				sb.WriteString(fmt.Sprintf("🔢 <b>%s</b> <code>%d</code>\n", postIDLabel, fo.MessageOriginChannel.MessageID))
			}
			if fo.MessageOriginChannel.AuthorSignature != nil && *fo.MessageOriginChannel.AuthorSignature != "" {
				sb.WriteString(fmt.Sprintf("✍️ <b>%s</b> %s\n", sigLabel, escape(*fo.MessageOriginChannel.AuthorSignature)))
			}
			if fo.MessageOriginChannel.Date > 0 {
				t := time.Unix(int64(fo.MessageOriginChannel.Date), 0).UTC()
				sb.WriteString(fmt.Sprintf("🕒 <b>%s</b> %s UTC\n", dateLabel, t.Format("2006-01-02 15:04:05")))
			}
		}
	}

	return sb.String()
}

// formatContact formats shared Telegram contact card information.
func formatContact(c *models.Contact, lang Lang) string {
	if c == nil {
		return ""
	}

	var sb strings.Builder
	if lang == LangRU {
		sb.WriteString("<b>📇 Информация об отправленном контакте:</b>\n")
		if c.UserID != 0 {
			sb.WriteString(fmt.Sprintf("🆔 <b>ID пользователя:</b> <code>%d</code>\n", c.UserID))
		} else {
			sb.WriteString("⚠️ <i>Этот контакт не привязан к Telegram-аккаунту или ID недоступен.</i>\n")
		}
		name := c.FirstName
		if c.LastName != "" {
			name += " " + c.LastName
		}
		sb.WriteString(fmt.Sprintf("👤 <b>Имя:</b> %s\n", escape(name)))
		if c.PhoneNumber != "" {
			sb.WriteString(fmt.Sprintf("📞 <b>Телефон:</b> <code>%s</code>\n", escape(c.PhoneNumber)))
		}
	} else {
		sb.WriteString("<b>📇 Shared Contact Information:</b>\n")
		if c.UserID != 0 {
			sb.WriteString(fmt.Sprintf("🆔 <b>User ID:</b> <code>%d</code>\n", c.UserID))
		} else {
			sb.WriteString("⚠️ <i>This contact is not linked to a Telegram account or ID is unavailable.</i>\n")
		}
		name := c.FirstName
		if c.LastName != "" {
			name += " " + c.LastName
		}
		sb.WriteString(fmt.Sprintf("👤 <b>Name:</b> %s\n", escape(name)))
		if c.PhoneNumber != "" {
			sb.WriteString(fmt.Sprintf("📞 <b>Phone:</b> <code>%s</code>\n", escape(c.PhoneNumber)))
		}
	}

	return sb.String()
}

// getWelcomeText returns a friendly localized welcome message with user profile information.
func getWelcomeText(u *models.User, lang Lang) string {
	var sb strings.Builder
	if lang == LangRU {
		sb.WriteString("👋 <b>Добро пожаловать в Telegram ID Helper Bot!</b>\n\n")
		sb.WriteString("Этот бот мгновенно помогает узнать ваш Telegram ID, ID чатов, каналов и авторов пересланных сообщений без сбора и логирования данных (Zero-Logging).\n\n")
		sb.WriteString("<b>Быстрые действия:</b>\n")
		sb.WriteString("• Перешлите любое сообщение — узнать ID автора.\n")
		sb.WriteString("• Отправьте контакт (📎 Скрепка → Контакт) — узнать ID контакта.\n")
		sb.WriteString("• Добавьте бота в группу и напишите <code>/id</code> — узнать ID группы.\n\n")
		sb.WriteString("━━━━━━━━━━━━━━━━━━━━\n")
		sb.WriteString(formatUser(u, "👤 Ваш профиль:", lang))
		sb.WriteString("\n💡 Отправьте <code>/help</code> для подробной справки об устройстве Telegram ID.")
	} else {
		sb.WriteString("👋 <b>Welcome to Telegram ID Helper Bot!</b>\n\n")
		sb.WriteString("This bot helps you quickly find your Telegram ID, Chat/Group IDs, Channel IDs, and original author IDs of forwarded messages with complete privacy and zero data logging.\n\n")
		sb.WriteString("<b>Quick actions:</b>\n")
		sb.WriteString("• Forward any message here — to get the author's ID.\n")
		sb.WriteString("• Share a contact card (📎 Attachment → Contact) — to get contact ID.\n")
		sb.WriteString("• Add the bot to a group and send <code>/id</code> — to get group ID.\n\n")
		sb.WriteString("━━━━━━━━━━━━━━━━━━━━\n")
		sb.WriteString(formatUser(u, "👤 Your Profile:", lang))
		sb.WriteString("\n💡 Send <code>/help</code> for detailed information about how Telegram IDs work.")
	}
	return sb.String()
}

// getHelpText returns explanatory text about Telegram IDs and bot usage.
func getHelpText(lang Lang) string {
	if lang == LangRU {
		return `🤖 <b>Telegram ID Helper Bot</b>

Этот бот помогает узнать Telegram ID пользователей, групп, каналов и авторов пересланных сообщений.

<b>Как пользоваться:</b>
• <b>Личный ID:</b> Отправьте любое сообщение боту в ЛС.
• <b>ID автора сообщения:</b> Перешлите любое сообщение в ЛС боту.
• <b>Если ID автора скрыт приватностью:</b>
  1) Отправьте карточку контакта (📎 Скрепка → Контакт), если человек есть в телефонной книге.
  2) Или добавьте человека в группу и ответьте на его сообщение командой <code>/id</code>.
• <b>Команды в группах:</b>
  • <code>/id</code> — выводит ID группы (или данные пользователя при ответе на сообщение).
  • <code>/quit</code> — бот выходит из группы.
• <b>ID канала:</b> Перешлите сообщение из канала боту в ЛС.
• <b>О боте / Версия:</b> Отправьте команду <code>/about</code>.

━━━━━━━━━━━━━━━━━━━━
💡 <b>Как устроены ID в Telegram:</b>
• <b>Пользователи:</b> Положительное число (например: <code>123456789</code>).
• <b>Обычные группы:</b> Отрицательное число с минусом (например: <code>-123456789</code>).
• <b>Супергруппы и Каналы:</b> Отрицательное 13-значное число с префиксом <code>-100</code> (например: <code>-1001234567890</code>).
• <b>Темы/Топики форума:</b> Каждая ветка (Topic) имеет свой <code>MessageThreadID</code> — ID первого сообщения темы.`
	}

	return `🤖 <b>Telegram ID Helper Bot</b>

This bot helps you discover Telegram IDs for users, groups, channels, and authors of forwarded messages.

<b>How to use:</b>
• <b>Your ID:</b> Send any message to the bot in direct messages.
• <b>Forwarded Author's ID:</b> Forward any message to the bot.
• <b>If Forwarded Author ID is Hidden by Privacy:</b>
  1) Send a contact card (📎 Attachment → Contact) if saved in your phonebook.
  2) Or add them to a group and Reply to their message with <code>/id</code>.
• <b>Commands in Groups:</b>
  • <code>/id</code> — returns group info (or replied user info when replying to a message).
  • <code>/quit</code> — forces the bot to leave the group.
• <b>Channel ID:</b> Forward any post from the channel to the bot.
• <b>About & Version:</b> Send <code>/about</code>.

━━━━━━━━━━━━━━━━━━━━
💡 <b>How IDs work in Telegram:</b>
• <b>Users:</b> Positive integer (e.g., <code>123456789</code>).
• <b>Basic Groups:</b> Negative integer (e.g., <code>-123456789</code>).
• <b>Supergroups & Channels:</b> Negative 13-digit number with <code>-100</code> prefix (e.g., <code>-1001234567890</code>).
• <b>Forum Topics / Threads:</b> Unique <code>MessageThreadID</code> for each topic (ID of the initial topic message).`
}
