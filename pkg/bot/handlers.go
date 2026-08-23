package bot

import (
	"context"
	"fmt"
	"strings"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// setupBotCommands registers bot commands for both default (English) and Russian languages.
func setupBotCommands(ctx context.Context, b *tgbot.Bot) {
	// 1. Commands (Default / English)
	_, _ = b.SetMyCommands(ctx, &tgbot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "start", Description: "Start the bot and get your ID"},
			{Command: "id", Description: "Get ID and chat details"},
			{Command: "about", Description: "Build info, Git commit & transparency"},
			{Command: "quit", Description: "Leave group chat"},
			{Command: "help", Description: "Help & Telegram ID guide"},
		},
	})

	// 2. Commands (Russian)
	_, _ = b.SetMyCommands(ctx, &tgbot.SetMyCommandsParams{
		LanguageCode: "ru",
		Commands: []models.BotCommand{
			{Command: "start", Description: "Запустить бота и узнать свой ID"},
			{Command: "id", Description: "Узнать ID пользователя или чата"},
			{Command: "about", Description: "О боте, Git ревизия и верификация"},
			{Command: "quit", Description: "Покинуть группу"},
			{Command: "help", Description: "Справка и как устроены Telegram ID"},
		},
	})
}

// HandleUpdate is the centralized dispatcher for all Telegram updates.
// Designed with zero logging and robust panic recovery.
func HandleUpdate(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	defer func() {
		// Suppress panics safely to prevent crashing
		_ = recover()
	}()

	if update == nil {
		return
	}

	if update.Message != nil {
		handleMessage(ctx, b, update.Message)
		return
	}

	if update.ChannelPost != nil {
		handleChannelPost(ctx, b, update.ChannelPost)
		return
	}

	if update.MyChatMember != nil {
		handleMyChatMember(ctx, b, update.MyChatMember)
		return
	}
}

// handleMessage processes incoming messages in private chats and groups.
func handleMessage(ctx context.Context, b *tgbot.Bot, msg *models.Message) {
	if msg == nil {
		return
	}

	lang := detectLang(msg.From)
	text := strings.TrimSpace(msg.Text)
	cmd := ""
	if strings.HasPrefix(text, "/") {
		parts := strings.Fields(text)
		if len(parts) > 0 {
			cmd = strings.ToLower(parts[0])
			// Strip @botusername from command if present (e.g. /start@MyBot -> /start)
			if idx := strings.Index(cmd, "@"); idx != -1 {
				cmd = cmd[:idx]
			}
		}
	}

	// /start command: show localized welcome message + sender's user info
	if cmd == "/start" {
		replyMessage(ctx, b, msg, getWelcomeText(msg.From, lang))
		return
	}

	// /help command: show help text
	if cmd == "/help" {
		replyMessage(ctx, b, msg, getHelpText(lang))
		return
	}

	// /about or /version command: show build transparency & git revision
	if cmd == "/about" || cmd == "/version" {
		replyMessage(ctx, b, msg, formatVersion(lang))
		return
	}

	switch msg.Chat.Type {
	case models.ChatTypePrivate:
		handlePrivateMessage(ctx, b, msg, lang)
	case models.ChatTypeGroup, models.ChatTypeSupergroup:
		handleGroupMessage(ctx, b, msg, cmd, lang)
	}
}

// handlePrivateMessage handles messages received in a direct message dialog.
func handlePrivateMessage(ctx context.Context, b *tgbot.Bot, msg *models.Message, lang Lang) {
	var sb strings.Builder

	// Check if this message is a shared contact card
	if msg.Contact != nil {
		sb.WriteString(formatContact(msg.Contact, lang))
	} else if msg.ForwardOrigin != nil {
		// Check if this message was forwarded: only send forwarded author / chat info
		sb.WriteString(formatForwardOrigin(msg.ForwardOrigin, lang))
	} else {
		// Normal message: send user profile info
		header := "👤 Your Profile Information:"
		if lang == LangRU {
			header = "👤 Информация о вашем профиле:"
		}
		sb.WriteString(formatUser(msg.From, header, lang))
	}

	replyMessage(ctx, b, msg, sb.String())
}

// handleGroupMessage handles messages sent in basic groups and supergroups.
func handleGroupMessage(ctx context.Context, b *tgbot.Bot, msg *models.Message, cmd string, lang Lang) {
	// Ignore service messages (member added/removed, chat created, etc.)
	if len(msg.NewChatMembers) > 0 || msg.LeftChatMember != nil || msg.GroupChatCreated || msg.SupergroupChatCreated {
		return
	}

	// Handle /quit or /leave commands to exit group explicitly
	if cmd == "/quit" || cmd == "/leave" || cmd == "/exit" {
		byeText := "👋 Leaving group..."
		if lang == LangRU {
			byeText = "👋 Выхожу из группы..."
		}
		replyMessage(ctx, b, msg, byeText)
		_, _ = b.LeaveChat(ctx, &tgbot.LeaveChatParams{
			ChatID: msg.Chat.ID,
		})
		return
	}

	// In groups, ONLY respond to explicit commands (/id, /info, /chatid)
	isExplicitCommand := cmd == "/id" || cmd == "/info" || cmd == "/chatid"
	if !isExplicitCommand {
		return
	}

	var sb strings.Builder
	groupHeader := "👥 Group Information:"
	threadLabel := "Topic / Thread ID:"
	replyHeader := "🎯 Replied User Information:"
	senderHeader := "👤 Command Sender:"

	if lang == LangRU {
		groupHeader = "👥 Информация о текущей группе:"
		threadLabel = "ID темы (Topic/Thread ID):"
		replyHeader = "🎯 Информация о пользователе из ответа (Reply):"
		senderHeader = "👤 Отправитель команды:"
	}

	sb.WriteString(formatChat(&msg.Chat, groupHeader, lang))

	// If forum topic/thread is active
	if msg.MessageThreadID != 0 {
		sb.WriteString(fmt.Sprintf("🧵 <b>%s</b> <code>%d</code>\n", threadLabel, msg.MessageThreadID))
	}

	sb.WriteString("\n━━━━━━━━━━━━━━━━━━━━\n")

	if msg.ReplyToMessage != nil {
		if msg.ReplyToMessage.SenderChat != nil {
			chatHeader := "📢 Replied Channel / Anonymous Admin:"
			if lang == LangRU {
				chatHeader = "📢 Канал / Анонимный администратор из ответа:"
			}
			sb.WriteString(formatChat(msg.ReplyToMessage.SenderChat, chatHeader, lang))
			sb.WriteString("\n━━━━━━━━━━━━━━━━━━━━\n")
		} else if msg.ReplyToMessage.From != nil {
			sb.WriteString(formatUser(msg.ReplyToMessage.From, replyHeader, lang))
			if msg.ReplyToMessage.ForwardOrigin != nil {
				sb.WriteString("\n")
				sb.WriteString(formatForwardOrigin(msg.ReplyToMessage.ForwardOrigin, lang))
			}
			sb.WriteString("\n━━━━━━━━━━━━━━━━━━━━\n")
		}
	}

	// Sender user info
	sb.WriteString(formatUser(msg.From, senderHeader, lang))

	replyMessage(ctx, b, msg, sb.String())
}

// handleChannelPost handles posts in channels where the bot is added as an administrator.
func handleChannelPost(ctx context.Context, b *tgbot.Bot, msg *models.Message) {
	lang := LangEN // Default to English for channel posts

	var sb strings.Builder
	sb.WriteString(formatChat(&msg.Chat, "📢 Channel Information:", lang))
	if msg.ID > 0 {
		sb.WriteString(fmt.Sprintf("🔢 <b>Channel Message ID:</b> <code>%d</code>\n", msg.ID))
	}

	replyMessage(ctx, b, msg, sb.String())
}

// handleMyChatMember sends a welcome message when the bot is added to a group or channel.
func handleMyChatMember(ctx context.Context, b *tgbot.Bot, update *models.ChatMemberUpdated) {
	status := update.NewChatMember.Type
	if status != models.ChatMemberTypeMember && status != models.ChatMemberTypeAdministrator {
		return
	}

	lang := detectLang(&update.From)

	var sb strings.Builder
	if lang == LangRU {
		sb.WriteString("👋 <b>Привет! Бот добавлен в чат.</b>\n\n")
		sb.WriteString(formatChat(&update.Chat, "👥 Информация о чате:", lang))
		sb.WriteString("\n💡 <b>Команды:</b>\n• <code>/id</code> — получить информацию о группе или о пользователе (при ответе на сообщение)\n• <code>/quit</code> — попросить бота выйти из группы.")
	} else {
		sb.WriteString("👋 <b>Hello! Bot has been added to this chat.</b>\n\n")
		sb.WriteString(formatChat(&update.Chat, "👥 Chat Information:", lang))
		sb.WriteString("\n💡 <b>Commands:</b>\n• <code>/id</code> — get chat/user info (or reply to a message)\n• <code>/quit</code> — make bot leave the group.")
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:    update.Chat.ID,
		Text:      sb.String(),
		ParseMode: models.ParseModeHTML,
	})
}

// replyMessage sends an HTML-formatted reply to the given message.
func replyMessage(ctx context.Context, b *tgbot.Bot, msg *models.Message, text string) {
	if msg == nil || text == "" {
		return
	}

	params := &tgbot.SendMessageParams{
		ChatID: msg.Chat.ID,
		Text:   text,
		ReplyParameters: &models.ReplyParameters{
			MessageID: msg.ID,
		},
		ParseMode: models.ParseModeHTML,
	}

	if msg.MessageThreadID != 0 {
		params.MessageThreadID = msg.MessageThreadID
	}

	_, err := b.SendMessage(ctx, params)
	if err != nil {
		// Fallback: try sending without reply parameter in case original message was deleted
		params.ReplyParameters = nil
		_, _ = b.SendMessage(ctx, params)
	}
}
